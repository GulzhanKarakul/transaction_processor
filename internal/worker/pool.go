package worker

import (
	"context"
	"sync"

	"github.com/GulzhanKarakul/transaction_processor/internal/model"
	"github.com/GulzhanKarakul/transaction_processor/internal/pipeline"
)

const NumWorkers = 5

// Worker Pool - параллельно обрабоатывает транзакции
type Pool struct {
	pipeline *pipeline.Pipeline
	jobs     chan model.Transaction
	results  chan model.Result
	wg       sync.WaitGroup
}

// NewPool создает экземпляр Worker Pool c буф каналами
func NewPool(p *pipeline.Pipeline) *Pool {
	return &Pool{pipeline: p, jobs: make(chan model.Transaction, 100), results: make(chan model.Result, 100)}
}

// Start - запускает параллельно число(NumWorkers) горутин
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < NumWorkers; i++ {
		p.wg.Add(1)
		go p.runWorker(ctx)
	}
}

// Один воркер - грутина, которая запустится параллельно в Start с другими
// закроется когда закончатся jobs или контекст отменен
func (p *Pool) runWorker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case t, ok := <-p.jobs:
			if !ok {
				return
			}
			p.results <- p.pipeline.Process(t)
		case <-ctx.Done():
			return
		}
	}
}

// Отправляет транзакцию в очередь
func (p *Pool) Submit(t model.Transaction) {
	p.jobs <- t
}

// возращает канал результатов только для чтения
func (p *Pool) Results() <-chan model.Result {
	return p.results
}

// graceful shutdown
func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}
