package model

import "time"

type Account struct {
	ID                   int           `json:"id"`
	Balance              float64       `json:"balance"`
	LastTransferTime     time.Time     `json:"last_transfer_time"`
	OutgoingTransactions []Transaction `json:"outgoing_transactions,omitempty"`
	IncomingTransactions []Transaction `json:"incoming_transactions,omitempty"`
}
