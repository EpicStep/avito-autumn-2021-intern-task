package balance

import (
	"errors"
	balanceDB "github.com/EpicStep/avito-autumn-2021-intern-task/internal/balance/database"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/convertor"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/jsonutil"
	v1 "github.com/EpicStep/avito-autumn-2021-intern-task/pkg/api/v1"
	"github.com/EpicStep/avito-autumn-2021-intern-task/pkg/database"
	"github.com/jackc/pgx/v4"
	"net/http"
	"strconv"
)

// Service balance.
type Service struct {
	db         *balanceDB.BalanceDB
	cConvertor *convertor.CurrencyConvertor
}

// New returns new balance service.
func New(db *database.DB, cc *convertor.CurrencyConvertor) *Service {
	return &Service{
		db:         balanceDB.NewBalanceDB(db),
		cConvertor: cc,
	}
}

// GetBalance GET /api/balance
func (s *Service) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, "Validation error"))
		return
	}

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "RUB"
	}

	balanceAccount, err := s.db.GetBalanceAccountByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonutil.MarshalResponse(w, http.StatusNotFound, jsonutil.NewError(3, "Account not found"))
			return
		} else {
			jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(3, "Error while get balance data"))
			return
		}
	}

	if currency == "RUB" {
		jsonutil.MarshalResponse(w, http.StatusOK, v1.GetBalanceResponse{
			Balance:  balanceAccount.Balance,
			Currency: currency,
		})
	} else {
		c, err := s.cConvertor.Convert(balanceAccount.Balance, currency)
		if err != nil {
			jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(4, err.Error()))
			return
		}

		jsonutil.MarshalResponse(w, http.StatusOK, v1.GetBalanceResponse{
			Balance:  c,
			Currency: currency,
		})
	}
}

type controlBalanceRequest struct {
	Amount  float64 `json:"amount"`
	Comment string  `json:"comment"`
}

func (r *controlBalanceRequest) validate() error {
	if r.Amount == 0 {
		return errors.New("amount must to be not 0")
	}

	return nil
}

// ControlBalance POST /api/balance
func (s *Service) ControlBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, "Validation error"))
		return
	}

	var req controlBalanceRequest

	unmarshallStatusCode, err := jsonutil.Unmarshal(w, r, &req)
	if err != nil {
		jsonutil.MarshalResponse(w, unmarshallStatusCode, jsonutil.NewError(3, err.Error()))
		return
	}

	if err := req.validate(); err != nil {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, err.Error()))
		return
	}

	err = s.db.UpdateBalance(ctx, id, req.Amount, req.Comment)
	if err != nil {
		if errors.Is(err, balanceDB.ErrBalanceMustBePositive) {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(5, "Balance can't be negative"))
			return
		} else {
			jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(3, "Error while update account"))
			return
		}
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(1))
}
