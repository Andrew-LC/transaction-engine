package service

import (
	"errors"

	"transaction-engine/domain"
	"transaction-engine/models"
	"transaction-engine/store"
	"transaction-engine/utils"
)

var (
	ErrInvalidPin             = errors.New("invalid pin")
	ErrCardExists             = errors.New("card already exists")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

type Service struct {
	store store.Store
}

func NewService(store store.Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateCard(cardRequest domain.NewCardRequest) (models.Card, error) {
	newCard := models.NewCard(
		cardRequest.CardNumber,
		cardRequest.CardHolder,
		cardRequest.Pin,
		cardRequest.Amount,
	)

	err := s.store.CreateCard(&newCard)
	if err != nil {
		return models.Card{}, err
	}

	return newCard, nil
}

func (s *Service) GetBalance(cardNumber int64) (int64, error) {
	card, err := s.store.GetCard(cardNumber)
	if err != nil {
		return 0, err
	}
	return card.Balance, nil
}

func (s *Service) GetTransactions(cardNumber int64) ([]models.Transaction, error) {
	_, err := s.store.GetCard(cardNumber)
	if err != nil {
		return []models.Transaction{}, err
	}
	return s.store.GetTransactions(cardNumber)
}

func (s *Service) ProcessTransaction(req domain.TransactionRequest) (int64, error) {
	card, err := s.store.GetCard(req.CardNumber)
	if err != nil {
		return 0, err
	}

	hash := utils.Hash(req.Pin)
	if card.PinHash != hash {
		return 0, ErrInvalidPin
	}

	switch req.Type {
	case "withdraw":
		if req.Amount > card.Balance {
			return 0, ErrInsufficientFunds
		}
		newBalance := card.Balance - req.Amount
		if err := s.store.UpdateBalance(card.CardNumber, newBalance); err != nil {
			return 0, err
		}
		if err := s.store.AddTransaction(models.NewTransaction(
			card.CardNumber,
			domain.TransactionTypeWithdraw,
			req.Amount,
			domain.TransactionStatusSuccess,
		)); err != nil {
			return 0, err
		}
		return newBalance, nil

	case "topup":
		newBalance := card.Balance + req.Amount
		if err := s.store.UpdateBalance(card.CardNumber, newBalance); err != nil {
			return 0, err
		}
		if err := s.store.AddTransaction(models.NewTransaction(
			card.CardNumber,
			domain.TransactionTypeTopUp,
			req.Amount,
			domain.TransactionStatusSuccess,
		)); err != nil {
			return 0, err
		}
		return newBalance, nil

	default:
		return 0, ErrInvalidTransactionType
	}
}
