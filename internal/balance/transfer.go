package balance

import (
	"errors"
	"fmt"
	balanceDB "github.com/EpicStep/avito-autumn-2021-intern-task/internal/balance/database"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/balance/model"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/jsonutil"
	"net/http"
)

type transferRequest struct {
	IDFrom  int     `json:"id_from"`
	IDTo    int     `json:"id_to"`
	Amount  float64 `json:"amount"`
	Comment string  `json:"comment"`
}

func (r *transferRequest) validate() error {
	if r.IDFrom == 0 || r.IDTo == 0 {
		return errors.New("you can't transfer money to/from system")
	}

	if r.Amount <= 0 {
		return errors.New("amount must be > 0")
	}

	return nil
}

// Transfer POST /api/balance/transfer
func (s *Service) Transfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req transferRequest

	unmarshallStatusCode, err := jsonutil.Unmarshal(w, r, &req)
	if err != nil {
		jsonutil.MarshalResponse(w, unmarshallStatusCode, jsonutil.NewError(3, err.Error()))
		return
	}

	if err := req.validate(); err != nil {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, err.Error()))
		return
	}

	th := model.TransactionHistory{
		IDFrom:  req.IDFrom,
		IDTo:    req.IDTo,
		Amount:  req.Amount,
		Comment: req.Comment,
	}

	th.Prepare()

	err = s.db.Transfer(ctx, &th)
	if err != nil {
		if errors.Is(err, balanceDB.ErrBalanceMustBePositive) {
			jsonutil.MarshalResponse(w, http.StatusConflict, jsonutil.NewError(5, "After transfer your balance will be < 0"))
			return
		}

		if errors.Is(err, balanceDB.ErrSenderNotExist) {
			jsonutil.MarshalResponse(w, http.StatusConflict, jsonutil.NewError(6, fmt.Sprintf("Transfer sender account #%d dosent exist", th.IDFrom)))
			return
		}

		if errors.Is(err, balanceDB.ErrReceiverNotExist) {
			jsonutil.MarshalResponse(w, http.StatusConflict, jsonutil.NewError(7, fmt.Sprintf("Transfer reciever account #%d dosent exist", th.IDTo)))
			return
		}

		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Failed to create transfer"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(1))
}
