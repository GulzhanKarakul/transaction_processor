package pipeline

import (
	"errors"
	"fmt"
	"time"

	"github.com/GulzhanKarakul/transaction_processor/internal/model"
	"github.com/GulzhanKarakul/transaction_processor/internal/storage"
)

// Pipeline для обраюотки транзакций последовательно: валидация -> бонус -> сохранение
type Pipeline struct {
	storage *storage.Storage
}

// NewPipeline для создания экземпляра Pipeline с переданным хранилищем
func NewPipeline(s *storage.Storage) *Pipeline {
	return &Pipeline{storage: s}
}

// Process пропускает каждую транзакцию через Pipeline
// при ошибке возвращает Result с ошибкой Err
func (p *Pipeline) Process(t model.Transaction) model.Result {
	if err := validate(t); err != nil {
		t.Status = model.StatusFailed
		t.Error = err.Error()
		return model.Result{Transaction: t, Err: err}
	}

	t = addBonus(t)

	if err := p.save(t); err != nil {
		t.Status = model.StatusFailed
		t.Error = err.Error()
		return model.Result{Transaction: t, Err: err}
	}

	return model.Result{Transaction: t}
}

func validate(t model.Transaction) error {
	var errs []error
	if t.Amount <= 0 {
		errs = append(errs, fmt.Errorf("pipeline: validation error transaction %s: amount %.2f is invalid", t.ID, t.Amount))
	}
	if t.UserID == "" {
		errs = append(errs, fmt.Errorf("pipeline: validation error transaction %s: user id is empty", t.ID))
	}
	return errors.Join(errs...)
}

func addBonus(t model.Transaction) model.Transaction {
	t.BonusAmount = t.Amount * 0.05
	return t
}

func (p *Pipeline) save(t model.Transaction) error {
	now := time.Now()
	t.Status = model.StatusCompleted
	t.ProcessedAt = &now
	p.storage.Save(t)
	return nil
}
