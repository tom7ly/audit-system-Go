package utils

import (
	"context"
	"sync"
	"time"
)

// Operation represents a function that performs a database operation.
type Operation func(ctx context.Context) error

// Queue manages a queue of operations to be processed.
type Queue struct {
	ops       chan Operation
	wg        sync.WaitGroup
	close     chan struct{}
	batchSize int
}

// NewQueue creates a new Queue with the specified buffer and batch size.
func NewQueue(bufferSize, batchSize int) *Queue {
	q := &Queue{
		ops:       make(chan Operation, bufferSize),
		close:     make(chan struct{}),
		batchSize: batchSize,
	}
	go q.run()
	return q
}

// run processes operations in batches.
func (q *Queue) run() {
	for {
		select {
		case <-q.close:
			return
		default:
			q.processBatch()
		}
	}
}

// processBatch processes a batch of operations.
func (q *Queue) processBatch() {
	batch := make([]Operation, 0, q.batchSize)
	timeout := time.After(100 * time.Millisecond)

	for {
		select {
		case op := <-q.ops:
			batch = append(batch, op)
			if len(batch) >= q.batchSize {
				q.executeBatch(batch)
				return
			}
		case <-timeout:
			if len(batch) > 0 {
				q.executeBatch(batch)
			}
			return
		}
	}
}

// executeBatch executes a batch of operations.
func (q *Queue) executeBatch(batch []Operation) {
	q.wg.Add(len(batch))
	for _, op := range batch {
		go func(op Operation) {
			defer q.wg.Done()
			op(context.Background())
		}(op)
	}
	q.wg.Wait()
}

// Add adds a new operation to the queue.
func (q *Queue) Add(op Operation) {
	q.ops <- op
}

// Shutdown shuts down the queue and waits for all operations to complete.
func (q *Queue) Shutdown() {
	close(q.ops)
	q.wg.Wait()
	close(q.close)
}
