package model

import "time"

// TransactionHistory struct.
type TransactionHistory struct {
	IDFrom    int
	IDTo      int
	Amount    float64
	Comment   string
	CreatedAt time.Time
}

// Prepare model to insert to DB.
func (m *TransactionHistory) Prepare() {
	m.CreatedAt = time.Now()
}
