package models

import "time"

type CardStatus int
type TransactionStatus int
type TransactionType int

const (
	CardStatusActive CardStatus = iota
	CardStatusBlocked
)

const (
	TransactionStatusSuccess TransactionStatus = iota
	TransactionStatusFailed
)

const (
	TransactionTypeWithdraw TransactionType = iota
	TransactionTypeTopUp
)

type Card struct {
	CardNumber int64     `json:"card_number"`
	CardHolder string    `json:"card_holder"`
	PinHash    string    `json:"-"`
	Balance    int64     `json:"balance"`
	Status     CardStatus    `json:"status"`  
}

type Transaction struct {
	TransactionID string             `json:"transaction_id"`
	CardNumber    int64              `json:"card_number"`
	Type          TransactionType    `json:"type"`
	Amount        int64              `json:"amount"`
	Status        TransactionStatus  `json:"transaction_status"`
	TimeStamp     time.Time          `json:"created_at"`
}
