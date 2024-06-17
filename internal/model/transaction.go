package model

import "time"

type Transaction struct {
	ID            int       `json:"id"`
	Amount        float64   `json:"amount"`
	Timestamp     time.Time `json:"timestamp"`
	FromAccountID int       `json:"from_account_id"`
	ToAccountID   int       `json:"to_account_id"`
}
