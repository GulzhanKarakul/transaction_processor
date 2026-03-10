package worker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/GulzhanKarakul/transaction_processor/internal/model"
	"github.com/GulzhanKarakul/transaction_processor/internal/pipeline"
	"github.com/GulzhanKarakul/transaction_processor/internal/storage"
)

func TestPool_ProcessesTransactions(t *testing.T) {
	store := storage.NewStorage()
	pipe := pipeline.NewPipeline(store)
	pool := NewPool(pipe)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// отправляем 10 транзакций
	for i := 0; i < 10; i++ {
		pool.Submit(model.Transaction{
			ID:     fmt.Sprintf("tx-%d", i),
			UserID: "user-1",
			Amount: 1000,
			Status: model.StatusPending,
		})
	}

	pool.Stop()

	// собираем результаты
	results := []model.Result{}
	for r := range pool.Results() {
		results = append(results, r)
	}

	if len(results) != 10 {
		t.Errorf("expected 10 results, got %d", len(results))
	}

	for _, r := range results {
		if r.Err != nil {
			t.Errorf("expected no error, got %v", r.Err)
		}
		if r.Transaction.Status != model.StatusCompleted {
			t.Errorf("expected status completed, got %s", r.Transaction.Status)
		}
	}
}

func TestPool_ContextCancel(t *testing.T) {
	store := storage.NewStorage()
	pipe := pipeline.NewPipeline(store)
	pool := NewPool(pipe)

	ctx, cancel := context.WithCancel(context.Background())
	pool.Start(ctx)

	// отменяем контекст — воркеры должны остановиться
	cancel()

	// даём горутинам время завершиться
	time.Sleep(50 * time.Millisecond)

	// Stop не должен завис
	done := make(chan struct{})
	go func() {
		pool.Stop()
		close(done)
	}()

	select {
	case <-done:
		// всё хорошо
	case <-time.After(2 * time.Second):
		t.Error("Stop() hung — possible deadlock")
	}
}
