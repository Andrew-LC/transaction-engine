package domain

type TransactionRequest struct {
	CardNumber  int64    `json:"card_number"`
	Pin         string    `json:"pin"`
	Type        string   `json:"type`
	Amount      int64     `json:"amount"`
}

type BalanceRequest struct {
	CardNumber  int64    `json:"card_number"`
	Pin         string    `json:"pin"`
}
