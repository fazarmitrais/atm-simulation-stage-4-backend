package entity

import "time"

type Transaction struct {
	Type                    string    `json:"type"`
	AccountNumber           string    `json:"accountNumber"`
	TransferToAccountNumber string    `json:"transferToAccountNumber"`
	Amount                  float64   `json:"amount"`
	Date                    time.Time `json:"date"`
}

const (
	TYPE_WITHDRAWAL = "WITHDRAWAL"
	TYPE_TRANSFER   = "TRANSFER"
)
