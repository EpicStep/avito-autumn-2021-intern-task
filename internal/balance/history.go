package balance

import (
	"errors"
	balanceDB "github.com/EpicStep/avito-autumn-2021-intern-task/internal/balance/database"
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/jsonutil"
	v1 "github.com/EpicStep/avito-autumn-2021-intern-task/pkg/api/v1"
	"net/http"
	"strconv"
	"strings"
)

// TransactionsHistory GET /api/balance/history
func (s *Service) TransactionsHistory(w http.ResponseWriter, r *http.Request) {
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

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	if limit > 100 {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, "Limit must be <= 100"))
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	if sortBy != "created_at" && sortBy != "amount" {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, "SortBy param must be created_at or amount"))
		return
	}

	sortOrder := r.URL.Query().Get("sort_order")
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	if strings.ToTitle(sortOrder) != "DESC" && strings.ToTitle(sortOrder) != "ASC" {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, "SortOrder param must be DESC or ASC"))
		return
	}

	history, count, err := s.db.GetHistory(ctx, id, limit, offset, sortBy, strings.ToTitle(sortOrder))
	if err != nil {
		if errors.Is(err, balanceDB.ErrAccountNotFound) {
			jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(2, "Account not found"))
		} else {
			jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(3, "Cannot get history data"))
		}

		return
	}

	response := v1.GetHistoryResponse{
		Count: count,
	}

	for _, v := range history {
		t := v1.Transaction{
			IDFrom:    v.IDFrom,
			IDTo:      v.IDTo,
			Currency:  currency,
			CreatedAt: v.CreatedAt,
			Comment:   v.Comment,
		}

		if currency == "RUB" {
			t.Amount = v.Amount
		} else {
			c, err := s.cConvertor.Convert(v.Amount, currency)
			if err != nil {
				jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(4, err.Error()))
				return
			}

			t.Amount = c
		}

		response.History = append(response.History, &t)
	}

	jsonutil.MarshalResponse(w, http.StatusOK, response)
}
