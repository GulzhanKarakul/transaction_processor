package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GulzhanKarakul/transaction_processor/internal/handler"
	"github.com/GulzhanKarakul/transaction_processor/internal/middleware"
	"github.com/GulzhanKarakul/transaction_processor/internal/pipeline"
	"github.com/GulzhanKarakul/transaction_processor/internal/storage"
	"github.com/GulzhanKarakul/transaction_processor/internal/worker"
)

func main() {
	// зависимости создаются снизу вверх — каждый слой зависит от предыдущего
	store := storage.NewStorage()
	pipe := pipeline.NewPipeline(store)
	pool := worker.NewPool(pipe)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	h := handler.NewHandler(pool, store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /transactions", h.CreateTransaction)
	mux.HandleFunc("GET /transactions", h.GetTransactions)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.Recovery(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// запускаем сервер в горутине чтобы не блокировать main
	go func() {
		log.Println("server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// блокируемся до сигнала от OS (Ctrl+C или kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")

	// даём активным запросам 5 секунд на завершение
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	pool.Stop()
	log.Println("server stopped")
}
