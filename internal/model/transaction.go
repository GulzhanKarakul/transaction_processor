package model

import "time"

// Status представляет статус обработки транзакции
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// Transaction представляет финансовую транзакцию в приложении
type Transaction struct {
	ID          string     // уникальный ID транзакции
	UserID      string     // ID пользователя
	Amount      float64    // сумма транзакции
	Status      Status     // статус обработки транзакции
	BonusAmount float64    // начисленные бонусы
	CreatedAt   time.Time  // время создания транзакции
	ProcessedAt *time.Time // время обработки, nil если не обработана
	Error       string     // ошибка при failed
}

// Result представляет результат обработки транзакций в WorkerPool
type Result struct {
	Transaction Transaction
	Err         error
}
