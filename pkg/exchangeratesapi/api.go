package exchangeratesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// GetCurrencyListResponse struct.
type GetCurrencyListResponse struct {
	Success   bool               `json:"success"`
	Timestamp int                `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}

// ExchangerAPI is a client to ExchangeratesAPI.
type ExchangerAPI struct {
	client *http.Client
	token  string
}

// New returns new ExchangerAPI.
func New(token string) *ExchangerAPI {
	return &ExchangerAPI{
		client: &http.Client{
			Timeout: time.Second * 3,
		},
		token: token,
	}
}

// GetCurrencyList returns GetCurrencyListResponse.
func (c *ExchangerAPI) GetCurrencyList() (*GetCurrencyListResponse, error) {
	var r GetCurrencyListResponse

	resp, err := http.Get(fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=%s&format=1", c.token))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get currency list")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
