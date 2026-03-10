package pipeline

import (
	"testing"

	"github.com/GulzhanKarakul/transaction_processor/internal/model"
	"github.com/GulzhanKarakul/transaction_processor/internal/storage"
)

// вспомогательная функция — создаёт валидную транзакцию для тестов
func newTestTransaction(amount float64, userID string) model.Transaction {
	return model.Transaction{
		ID:     "tx-test-1",
		UserID: userID,
		Amount: amount,
		Status: model.StatusPending,
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		tx      model.Transaction
		wantErr bool
	}{
		{
			name:    "valid transaction",
			tx:      newTestTransaction(1000, "user-1"),
			wantErr: false,
		},
		{
			name:    "negative amount",
			tx:      newTestTransaction(-100, "user-1"),
			wantErr: true,
		},
		{
			name:    "zero amount",
			tx:      newTestTransaction(0, "user-1"),
			wantErr: true,
		},
		{
			name:    "empty userID",
			tx:      newTestTransaction(1000, ""),
			wantErr: true,
		},
		{
			name:    "both invalid",
			tx:      newTestTransaction(-100, ""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddBonus(t *testing.T) {
	tx := newTestTransaction(1000, "user-1")
	result := addBonus(tx)

	if result.BonusAmount != 50.0 {
		t.Errorf("expected bonus 50.0, got %.2f", result.BonusAmount)
	}

	// оригинал не должен измениться
	if tx.BonusAmount != 0 {
		t.Error("original transaction should not be modified")
	}
}

func TestProcess_Success(t *testing.T) {
	store := storage.NewStorage()
	pipe := NewPipeline(store)

	tx := newTestTransaction(1000, "user-1")
	result := pipe.Process(tx)

	if result.Err != nil {
		t.Fatalf("expected no error, got %v", result.Err)
	}
	if result.Transaction.Status != model.StatusCompleted {
		t.Errorf("expected status completed, got %s", result.Transaction.Status)
	}
	if result.Transaction.BonusAmount != 50.0 {
		t.Errorf("expected bonus 50.0, got %.2f", result.Transaction.BonusAmount)
	}
	if result.Transaction.ProcessedAt == nil {
		t.Error("ProcessedAt should not be nil after processing")
	}
}

func TestProcess_ValidationFail(t *testing.T) {
	store := storage.NewStorage()
	pipe := NewPipeline(store)

	tx := newTestTransaction(-500, "")
	result := pipe.Process(tx)

	if result.Err == nil {
		t.Fatal("expected error, got nil")
	}
	if result.Transaction.Status != model.StatusFailed {
		t.Errorf("expected status failed, got %s", result.Transaction.Status)
	}
}
