package storage

import (
	"fmt"
	"sync"
	"testing"

	"github.com/GulzhanKarakul/transaction_processor/internal/model"
)

func newTestTransaction(id string, amount float64) model.Transaction {
	return model.Transaction{
		ID:     id,
		UserID: "user-1",
		Amount: amount,
		Status: model.StatusPending,
	}
}

func TestSave_And_GetByID(t *testing.T) {
	store := NewStorage()
	tx := newTestTransaction("tx-1", 1000)

	store.Save(tx)

	got, ok := store.GetByID("tx-1")
	if !ok {
		t.Fatal("expected transaction, got nothing")
	}
	if got.ID != tx.ID {
		t.Errorf("expected ID %s, got %s", tx.ID, got.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	store := NewStorage()

	_, ok := store.GetByID("non-existent")
	if ok {
		t.Error("expected not found, got true")
	}
}

func TestGetAll(t *testing.T) {
	store := NewStorage()
	store.Save(newTestTransaction("tx-1", 1000))
	store.Save(newTestTransaction("tx-2", 2000))
	store.Save(newTestTransaction("tx-3", 3000))

	all := store.GetAll()
	if len(all) != 3 {
		t.Errorf("expected 3 transactions, got %d", len(all))
	}
}

// TestConcurrent_Save — 100 горутин пишут одновременно
// race detector не должен ничего найти: go test -race ./...
func TestConcurrent_Save(t *testing.T) {
	store := NewStorage()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			store.Save(newTestTransaction(fmt.Sprintf("tx-%d", i), float64(i*100)))
		}(i)
	}

	wg.Wait()

	all := store.GetAll()
	if len(all) != 100 {
		t.Errorf("expected 100 transactions, got %d", len(all))
	}
}
