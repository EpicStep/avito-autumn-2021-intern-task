package v1

import "time"

type GetHistoryResponse struct {
	Count   int            `json:"count"`
	History []*Transaction `json:"history"`
}

type Transaction struct {
	IDFrom    int       `json:"id_from"`
	IDTo      int       `json:"id_to"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
