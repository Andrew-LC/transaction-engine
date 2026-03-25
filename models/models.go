package models

import (
	"time"
	"github.com/google/uuid"
	"transaction-engine/domain"
	"transaction-engine/utils"
)

type Card struct {
	CardNumber int64        `json:"card_number"`
	CardHolder string       `json:"card_holder"`
	PinHash    string       `json:"-"`
	Balance    int64        `json:"balance"`
	Status     domain.CardStatus   `json:"status"`  
}

type Transaction struct {
	TransactionID string                    `json:"transaction_id"`
	CardNumber    int64                     `json:"card_number"`
	Type          domain.TransactionType    `json:"type"`
	Amount        int64                     `json:"amount"`
	Status        domain.TransactionStatus  `json:"transaction_status"`
	TimeStamp     time.Time                 `json:"created_at"`
}


func NewCard(
	cardNumber int64,
	cardHolder string,
	pin        string,
	amount int64,
) Card {
	return Card{
		CardNumber: cardNumber,
		CardHolder: cardHolder,
		PinHash:    utils.Hash(pin),
		Status:     domain.CardStatusActive,
		Balance:    amount,      
	}
}

func NewTransaction(
	cardNumber int64,
	txType domain.TransactionType,
	amount int64,
	status domain.TransactionStatus,
) Transaction {
	return Transaction{
		TransactionID: uuid.NewString(),
		CardNumber:    cardNumber,
		Type:          txType,
		Amount:        amount,
		Status:        status, 
		TimeStamp:     time.Now(),
	}
}
