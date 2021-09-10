package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/balance/model"
	"github.com/EpicStep/avito-autumn-2021-intern-task/pkg/database"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// BalanceDB struct.
type BalanceDB struct {
	db *database.DB
}

// NewBalanceDB returns new BalanceDB.
func NewBalanceDB(db *database.DB) *BalanceDB {
	return &BalanceDB{db: db}
}

// ErrAccountNotFound error.
var ErrAccountNotFound = errors.New("account not found")

// ErrBalanceMustBePositive error.
var ErrBalanceMustBePositive = errors.New("balance must be positive number")

// ErrSenderNotExist error.
var ErrSenderNotExist = errors.New("sender not exist")

// ErrReceiverNotExist error.
var ErrReceiverNotExist = errors.New("receiver not exist")

// GetBalanceAccountByID from database.
func (db *BalanceDB) GetBalanceAccountByID(ctx context.Context, id int) (*model.Account, error) {
	var a model.Account

	err := db.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `
			SELECT
				id, balance
			FROM
				accounts
			WHERE
				id = $1
		`, id).Scan(&a.ID, &a.Balance)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetHistory from database.
func (db *BalanceDB) GetHistory(ctx context.Context, id int, limit int, offset int, sortBy string, sortOrder string) ([]*model.TransactionHistory, int, error) {
	var ths []*model.TransactionHistory
	var count int

	err := db.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf(`
			SELECT 
				id_from, id_to, amount, comment, created_at, count(*) OVER() AS count
			FROM
				transaction_history
			WHERE
				id_from = $1 OR id_to = $1
			ORDER BY %s %s
			LIMIT $2
			OFFSET $3
		`, sortBy, sortOrder), id, limit, offset)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var th model.TransactionHistory

			err := rows.Scan(&th.IDFrom, &th.IDTo, &th.Amount, &th.Comment, &th.CreatedAt, &count)
			if err != nil {
				return err
			}

			ths = append(ths, &th)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, 0, err
	}

	if len(ths) <= 0 {
		return nil, 0, ErrAccountNotFound
	}

	return ths, count, nil
}

// UpdateBalance in database.
func (db *BalanceDB) UpdateBalance(ctx context.Context, id int, amount float64, comment string) error {
	err := db.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		r, err := tx.Exec(ctx, `
			UPDATE
				accounts
			SET 
				balance = balance + $1
			WHERE
				id = $2
		`, amount, id)
		if err != nil {
			if pgerr, ok := err.(*pgconn.PgError); ok {
				if pgerr.Code == "23514" {
					return ErrBalanceMustBePositive
				}
			}

			return err
		}

		if r.RowsAffected() == 0 {
			if amount >= 0 {
				id, err = db.CreateUserInTx(ctx, tx, amount)
				if err != nil {
					return err
				}
			} else {
				return ErrBalanceMustBePositive
			}
		}

		th := model.TransactionHistory{
			Amount:  amount,
			Comment: comment,
		}

		if amount < 0 {
			th.IDFrom = id
			th.IDTo = 0
		} else {
			th.IDTo = id
		}

		th.Prepare()

		if err := db.CreateHistoryLog(ctx, tx, &th); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (db *BalanceDB) Transfer(ctx context.Context, h *model.TransactionHistory) error {
	err := db.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		rs, err := tx.Exec(ctx, `
			UPDATE
				accounts
			SET
				balance = balance - $1
			WHERE
				id = $2
		`, h.Amount, h.IDFrom)
		if err != nil {
			if pgerr, ok := err.(*pgconn.PgError); ok {
				if pgerr.Code == "23514" {
					return ErrBalanceMustBePositive
				}
			}

			return err
		}

		if rs.RowsAffected() == 0 {
			return ErrSenderNotExist
		}

		rr, err := tx.Exec(ctx, `
			UPDATE
				accounts
			SET
				balance = balance + $1
			WHERE
				id = $2
		`, h.Amount, h.IDTo)

		if err != nil {
			return err
		}

		if rr.RowsAffected() == 0 {
			return ErrReceiverNotExist
		}

		err = db.CreateHistoryLog(ctx, tx, h)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// CreateUserInTx is a function to create new user in DB.
func (db *BalanceDB) CreateUserInTx(ctx context.Context, tx pgx.Tx, amount float64) (int, error) {
	var id int

	err := tx.QueryRow(ctx, `
		INSERT INTO
			accounts(balance)
		VALUES
			($1)
		RETURNING id
	`, amount).Scan(&id)

	return id, err
}

// CreateHistoryLog is a function to create new history log in DB.
func (db *BalanceDB) CreateHistoryLog(ctx context.Context, tx pgx.Tx, h *model.TransactionHistory) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO 
			transaction_history
			(id_from, id_to, amount, comment, created_at)
		VALUES
			($1, $2, $3, $4, $5)
	`, h.IDFrom, h.IDTo, h.Amount, h.Comment, h.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
