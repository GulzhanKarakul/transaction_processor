package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GulzhanKarakul/transaction_processor/internal/model"
	"github.com/GulzhanKarakul/transaction_processor/internal/storage"
	"github.com/GulzhanKarakul/transaction_processor/internal/worker"
)

// Handler обрабатывает HTTP запросы и связывает transport слой с бизнес логикой
type Handler struct {
	pool    *worker.Pool
	storage *storage.Storage
}

// NewHandler создаёт Handler с переданными зависимостями
func NewHandler(p *worker.Pool, s *storage.Storage) *Handler {
	return &Handler{pool: p, storage: s}
}

// CreateTransaction принимает POST /transactions.
// Декодирует тело запроса, инициализирует транзакцию и отправляет в Worker Pool.
// Отвечает 202 Accepted — это значит "принято в обработку", не "выполнено".
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tx model.Transaction

	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json",
		})
		return
	}

	// инициализируем поля которые клиент не должен задавать сам
	tx.ID = fmt.Sprintf("tx-%d", time.Now().UnixNano())
	tx.Status = model.StatusPending
	tx.CreatedAt = time.Now()

	// Submit не блокирует — кладёт транзакцию в буферизированный канал и возвращается
	// воркеры подхватят её асинхронно
	h.pool.Submit(tx)

	// 202 а не 200 — потому что транзакция принята но ещё не обработана
	writeJSON(w, http.StatusAccepted, tx)
}

// GetTransactions принимает GET /transactions.
// Возвращает все транзакции из хранилища на момент запроса.
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	txs := h.storage.GetAll()
	writeJSON(w, http.StatusOK, txs)
}

// writeJSON — вспомогательная функция для JSON ответов.
// Header всегда устанавливается ДО WriteHeader — иначе Go их игнорирует.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
