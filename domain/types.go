package domain

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
