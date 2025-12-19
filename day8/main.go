package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Day 8: Concurrency Testing in Go

// ============================================
// 1. SafeCounter - Thread-safe counter using atomic operations
// ============================================

// SafeCounter is a thread-safe counter using atomic operations
type SafeCounter struct {
	value int64
}

// NewSafeCounter creates a new SafeCounter
func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}

// Increment adds 1 to the counter atomically
func (c *SafeCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

// Decrement subtracts 1 from the counter atomically
func (c *SafeCounter) Decrement() {
	atomic.AddInt64(&c.value, -1)
}

// Value returns the current counter value
func (c *SafeCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Add adds n to the counter atomically
func (c *SafeCounter) Add(n int64) {
	atomic.AddInt64(&c.value, n)
}

// ============================================
// 2. ConcurrentCache - Thread-safe cache with RWMutex
// ============================================

// CacheItem represents an item in the cache with TTL
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// ConcurrentCache is a thread-safe cache with TTL support
type ConcurrentCache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
	ttl  time.Duration
}

// NewConcurrentCache creates a new cache with default TTL
func NewConcurrentCache(ttl time.Duration) *ConcurrentCache {
	return &ConcurrentCache{
		data: make(map[string]CacheItem),
		ttl:  ttl,
	}
}

// Set adds or updates a value in the cache
func (c *ConcurrentCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves a value from the cache
func (c *ConcurrentCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Value, true
}

// Delete removes a key from the cache
func (c *ConcurrentCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Size returns the number of items in the cache
func (c *ConcurrentCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// Clear removes all items from the cache
func (c *ConcurrentCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]CacheItem)
}

// CleanupExpired removes expired items from the cache
func (c *ConcurrentCache) CleanupExpired() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	now := time.Now()
	for key, item := range c.data {
		if now.After(item.ExpiresAt) {
			delete(c.data, key)
			count++
		}
	}
	return count
}

// ============================================
// 3. WorkerPool - Process jobs with multiple workers
// ============================================

// Job represents a unit of work
type Job struct {
	ID      int
	Payload interface{}
}

// Result represents the result of processing a job
type Result struct {
	JobID  int
	Output interface{}
	Error  error
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	done       chan struct{}
	wg         sync.WaitGroup
	processor  func(Job) Result
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, processor func(Job) Result) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, 100),
		results:    make(chan Result, 100),
		done:       make(chan struct{}),
		processor:  processor,
	}
}

// Start begins the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker is a goroutine that processes jobs
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			result := wp.processor(job)
			wp.results <- result
		case <-wp.done:
			return
		}
	}
}

// Submit adds a job to the queue
func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}

// Results returns the results channel
func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.done)
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

// ============================================
// 4. Pipeline - Multi-stage data processing
// ============================================

// Pipeline represents a multi-stage data processing pipeline
type Pipeline struct {
	stages []func(ctx context.Context, in <-chan int) <-chan int
}

// NewPipeline creates a new pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStage adds a processing stage to the pipeline
func (p *Pipeline) AddStage(stage func(ctx context.Context, in <-chan int) <-chan int) {
	p.stages = append(p.stages, stage)
}

// Run executes the pipeline
func (p *Pipeline) Run(ctx context.Context, input <-chan int) <-chan int {
	current := input
	for _, stage := range p.stages {
		current = stage(ctx, current)
	}
	return current
}

// Generator creates a channel that emits numbers from 0 to n-1
func Generator(ctx context.Context, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Multiplier returns a stage that multiplies each input by n
func Multiplier(n int) func(ctx context.Context, in <-chan int) <-chan int {
	return func(ctx context.Context, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for val := range in {
				select {
				case out <- val * n:
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	}
}

// Filter returns a stage that only passes values matching the predicate
func Filter(predicate func(int) bool) func(ctx context.Context, in <-chan int) <-chan int {
	return func(ctx context.Context, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for val := range in {
				if predicate(val) {
					select {
					case out <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
		return out
	}
}

// ============================================
// 5. Fan-Out/Fan-In Pattern
// ============================================

// FanOut distributes work from one channel to multiple workers
func FanOut(ctx context.Context, input <-chan int, numWorkers int, worker func(int) int) []<-chan int {
	outputs := make([]<-chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		outputs[i] = func() <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for val := range input {
					select {
					case out <- worker(val):
					case <-ctx.Done():
						return
					}
				}
			}()
			return out
		}()
	}

	return outputs
}

// FanIn merges multiple channels into one
func FanIn(ctx context.Context, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Output function for each channel
	output := func(ch <-chan int) {
		defer wg.Done()
		for val := range ch {
			select {
			case out <- val:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	// Close out after all outputs are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ============================================
// 6. Semaphore - Limit concurrent operations
// ============================================

// Semaphore limits the number of concurrent operations
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a new semaphore with the given limit
func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, limit),
	}
}

// Acquire blocks until a slot is available
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// TryAcquire attempts to acquire without blocking
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release frees a slot
func (s *Semaphore) Release() {
	<-s.ch
}

// ============================================
// Main - Demo
// ============================================

func main() {
	fmt.Println("Day 8: Concurrency Testing")
	fmt.Println("Run 'go test -race -v' to execute tests with race detector")
	fmt.Println("")

	// Demo: SafeCounter
	fmt.Println("=== SafeCounter Demo ===")
	counter := NewSafeCounter()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("Counter value after 100 concurrent increments: %d\n", counter.Value())

	// Demo: ConcurrentCache
	fmt.Println("\n=== ConcurrentCache Demo ===")
	cache := NewConcurrentCache(5 * time.Second)
	cache.Set("key1", "value1")
	if val, ok := cache.Get("key1"); ok {
		fmt.Printf("Cache get: %v\n", val)
	}

	// Demo: WorkerPool
	fmt.Println("\n=== WorkerPool Demo ===")
	processor := func(job Job) Result {
		return Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("Processed job %d", job.ID),
		}
	}

	pool := NewWorkerPool(3, processor)
	pool.Start()

	// Submit jobs
	for i := 0; i < 5; i++ {
		pool.Submit(Job{ID: i, Payload: fmt.Sprintf("data-%d", i)})
	}

	// Collect results (with timeout)
	timeout := time.After(2 * time.Second)
	collected := 0
	for collected < 5 {
		select {
		case result := <-pool.Results():
			fmt.Printf("Result: Job %d -> %v\n", result.JobID, result.Output)
			collected++
		case <-timeout:
			fmt.Println("Timeout waiting for results")
			break
		}
	}

	pool.Stop()
	fmt.Println("\nDone! Run 'go test -race -v' to run tests.")
}
