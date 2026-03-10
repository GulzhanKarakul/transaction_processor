package storage

import (
	"sync"

	"github.com/GulzhanKarakul/transaction_processor/internal/model"
)

// Структура Storage  in memory хранилище
type Storage struct {
	transactions map[string]model.Transaction
	mu           sync.RWMutex
}

// NewStorage для создания экземпляра хранилища
func NewStorage() *Storage {
	return &Storage{transactions: make(map[string]model.Transaction)}
}

// Save— сохранить транзакцию
func (s *Storage) Save(t model.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[t.ID] = t
}

// GetByID — получить транзакцию по ID
func (s *Storage) GetByID(id string) (model.Transaction, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.transactions[id]
	return t, ok
}

// GetAll() — все транзакции
func (s *Storage) GetAll() []model.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Transaction, 0, len(s.transactions))
	for _, t := range s.transactions {
		result = append(result, t)
	}
	return result
}
