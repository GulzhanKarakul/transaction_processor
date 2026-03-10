package model

import "time"

// Status представляет статус транзакции
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// Transaction представляет финансовую транзакцию
type Transaction struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Amount      float64    `json:"amount"`
	Status      Status     `json:"status"`
	BonusAmount float64    `json:"bonus_amount"`
	CreatedAt   time.Time  `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	Error       string     `json:"error,omitempty"`
}

// Result представляет результат обработки транзакции в Worker Pool
type Result struct {
	Transaction Transaction
	Err         error
}