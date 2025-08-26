package queue

import (
	"log"
	"sync"
)

// ConcurrentQueue provides a concurrent execution queue with limits
type ConcurrentQueue struct {
	capacity int            // Maximum concurrency count
	sem      chan struct{}  // Semaphore channel for controlling concurrency
	wg       sync.WaitGroup // Wait for all tasks to complete
}

// NewConcurrentQueue creates a new concurrent queue
// capacity: Maximum concurrency count, must be greater than 0
func NewConcurrentQueue(capacity int) *ConcurrentQueue {
	if capacity <= 0 {
		panic("queue capacity must be greater than 0")
	}
	return &ConcurrentQueue{
		capacity: capacity,
		sem:      make(chan struct{}, capacity),
	}
}

// Go submits a task to the queue for asynchronous execution
// If the queue is full, it will block until there is an available slot
func (q *ConcurrentQueue) Go(task func()) {
	q.wg.Add(1)
	go func() {
		q.sem <- struct{}{} // Acquire semaphore
		defer func() {
			<-q.sem // Release semaphore
			q.wg.Done()
		}()

		// Execute task and handle potential panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("task panic recovered: %v", r)
			}
		}()

		task()
	}()
}

// Wait waits for all submitted tasks to complete
func (q *ConcurrentQueue) Wait() {
	q.wg.Wait()
}

// TryGo attempts to submit a task, returns false immediately if the queue is full
func (q *ConcurrentQueue) TryGo(task func()) bool {
	select {
	case q.sem <- struct{}{}: // Try to acquire semaphore
		q.wg.Add(1)
		go func() {
			defer func() {
				<-q.sem
				q.wg.Done()
				if r := recover(); r != nil {
					log.Printf("task panic recovered: %v", r)
				}
			}()
			task()
		}()
		return true
	default:
		return false
	}
}

// CurrentCount returns the number of tasks currently executing
func (q *ConcurrentQueue) CurrentCount() int {
	return len(q.sem)
}
