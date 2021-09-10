package v1

// GetBalanceResponse struct.
type GetBalanceResponse struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}
