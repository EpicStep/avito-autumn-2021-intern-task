package convertor

import (
	"fmt"
	"math"
)

// CurrencyConvertor is a struct with map and method to convert currency.
type CurrencyConvertor struct {
	currency map[string]float64
}

// NewCurrencyConvertor returns new CurrencyConvertor.
func NewCurrencyConvertor(list map[string]float64) *CurrencyConvertor {
	return &CurrencyConvertor{currency: list}
}

// Convert rub amount to another currency.
func (cc *CurrencyConvertor) Convert(rubAmount float64, toCurrency string) (float64, error) {
	// API на бесплатной версии предлагает только EUR как base валюту, приходится изворачиваться
	eurRubCurrency, ok := cc.currency["RUB"]
	if !ok {
		return 0, fmt.Errorf("currency %s dosent exist in DB", "RUB")
	}

	eurAmount := rubAmount / eurRubCurrency

	cToEur, ok := cc.currency[toCurrency]
	if !ok {
		return 0, fmt.Errorf("currency %s dosent exist in DB", toCurrency)
	}

	result := eurAmount * cToEur

	return math.Round(result*100) / 100, nil
}
