package worker

import (
	"context"
	"go-fiber-gorm/core/logger"
	"sync"
	"time"
)

// Task represents a job to be executed by workers
type Task func() error

// Pool represents a pool of workers
type Pool struct {
	tasks       chan Task
	concurrency int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewPool creates a new worker pool
func NewPool(concurrency int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		tasks:       make(chan Task, concurrency*10), // Buffer size is 10x concurrency
		concurrency: concurrency,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start starts the worker pool
func (p *Pool) Start() {
	logger.Info("Starting worker pool with", p.concurrency, "workers")

	// Start workers
	for i := 0; i < p.concurrency; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker is the goroutine processing tasks
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	logger.Info("Worker ", id, " started.")

	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				// Channel closed, worker should exit
				logger.Info("Worker ", id, " shutting down")
				return
			}

			// Execute the task
			startTime := time.Now()
			if err := task(); err != nil {
				logger.Error("Worker", id, "task error:", err)
			}

			logger.Info("Worker", id, "finished task in", time.Since(startTime))

		case <-p.ctx.Done():
			// Context canceled, worker should exit
			logger.Info("Worker", id, "received cancellation signal")
			return
		}
	}
}

// Submit adds a task to the worker pool
func (p *Pool) Submit(task Task) {
	select {
	case p.tasks <- task:
		// Task submitted successfully
	case <-p.ctx.Done():
		logger.Warn("Cannot submit task: worker pool is shutting down")
	}
}

// Stop gracefully stops the worker pool
func (p *Pool) Stop() {
	logger.Info("Stopping worker pool")

	// Signal for workers to exit
	p.cancel()

	// Close the task channel
	close(p.tasks)

	// Wait for all workers to finish
	p.wg.Wait()

	logger.Info("Worker pool stopped")
}
