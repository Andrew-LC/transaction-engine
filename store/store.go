package store

import (
    "errors"
    "sync"
    "time"

    "github.com/google/uuid"
    "transaction-engine/models"
    "transaction-engine/domain"
)

var (
    ErrCardNotFound      = errors.New("card not found")
    ErrCardBlocked       = errors.New("card is blocked")
    ErrInvalidPIN        = errors.New("invalid PIN")
    ErrInsufficientFunds = errors.New("insufficient balance")
    ErrInvalidAmount     = errors.New("amount must be greater than 0")
)

type Store interface {
    GetCard(cardNumber int64) (*models.Card, error)
    UpdateBalance(cardNumber int64, newBalance int64) error
    AddTransaction(transaction models.Transaction) error
    GetTransactions(cardNumber int64) ([]models.Transaction, error)
}

type InMemoryStore struct {
    mu           sync.RWMutex
    cards        map[int64]*models.Card
    transactions map[int64][]models.Transaction 
}

func NewStore() *InMemoryStore {
    s := &InMemoryStore{
        cards:        make(map[int64]*models.Card),
        transactions: make(map[int64][]models.Transaction),
    }
    s.seed()
    return s
}

func (s *InMemoryStore) seed() {
    s.cards[4123456789012345] = &models.Card{
        CardNumber: 4123456789012345,
        CardHolder: "John Doe",
        PinHash:    "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4", 
        Balance:    1000,
        Status:     domain.CardStatusActive,
    }
}

func (s *InMemoryStore) CreateCard(card *models.Card) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.cards[card.CardNumber] = card
    return nil
}

func (s *InMemoryStore) GetCard(cardNumber int64) (*models.Card, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    card, ok := s.cards[cardNumber]
    if !ok {
        return nil, ErrCardNotFound
    }
    copy := *card
    return &copy, nil
}

func (s *InMemoryStore) UpdateBalance(cardNumber int64, newBalance int64) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    card, ok := s.cards[cardNumber]
    if !ok {
        return ErrCardNotFound
    }
    card.Balance = newBalance
    return nil
}

func (s *InMemoryStore) AddTransaction(tx models.Transaction) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    tx.TransactionID = uuid.New().String() 
    tx.TimeStamp = time.Now()
    s.transactions[tx.CardNumber] = append(s.transactions[tx.CardNumber], tx) 
    return nil
}

func (s *InMemoryStore) GetTransactions(cardNumber int64) ([]models.Transaction, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    txs, ok := s.transactions[cardNumber]
    if !ok {
        return []models.Transaction{}, nil 
    }
    return txs, nil
}
