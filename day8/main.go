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
//
// WHY RWMutex instead of regular Mutex?
// -------------------------------------
// - Regular Mutex: Only ONE goroutine
//  can access at a time (read OR write)
// - RWMutex: MULTIPLE goroutines can READ simultaneously, but only ONE can WRITE
//
// This is perfect for caches because:
// - Reads are typically much more frequent than writes (90% reads, 10% writes)
// - Multiple readers don't conflict with each other
// - RWMutex allows concurrent reads = better performance
//
// LOCK TYPES:
// - Lock()/Unlock()   = EXCLUSIVE lock (for writes) - blocks ALL other access
// - RLock()/RUnlock() = SHARED lock (for reads) - allows other readers

// CacheItem represents an item in the cache with TTL (Time-To-Live)
//
// WHY these data types?
//   - Value: interface{} = Can store ANY type (string, int, struct, etc.)
//     This makes the cache generic - you can cache users, products, anything!
//   - ExpiresAt: time.Time = Exact timestamp when this item becomes invalid
//     Using time.Time (not duration) makes expiration checks simple: time.Now().After(ExpiresAt)
type CacheItem struct {
	Value     interface{} // interface{} = Go's "any type" - flexible but requires type assertion when reading
	ExpiresAt time.Time   // time.Time = precise moment of expiration (not relative duration)
}

// ConcurrentCache is a thread-safe cache with TTL support
//
// WHY these data types?
//   - data: map[string]CacheItem = Fast O(1) lookup by key
//     string keys are common (userID, productID, etc.)
//   - mu: sync.RWMutex = Protects 'data' from concurrent access
//     Without this, concurrent reads/writes would corrupt the map!
//   - ttl: time.Duration = How long items stay valid (e.g., 5*time.Minute)
//     Stored once, applied to all items for consistency
type ConcurrentCache struct {
	data map[string]CacheItem // The actual storage - map provides O(1) access
	mu   sync.RWMutex         // Protects 'data' - RWMutex allows concurrent reads
	ttl  time.Duration        // Default lifetime for cache items (e.g., 5 minutes)
}

// NewConcurrentCache creates a new cache with default TTL
//
// WHY return *ConcurrentCache (pointer)?
// - Mutex/RWMutex should NEVER be copied (causes bugs)
// - Returning pointer ensures everyone uses the SAME mutex
// - Also more efficient (no copying large struct)
//
// WHY make(map[string]CacheItem)?
// - Maps must be initialized before use (nil map panics on write)
// - make() allocates and initializes the map
func NewConcurrentCache(ttl time.Duration) *ConcurrentCache {
	return &ConcurrentCache{
		data: make(map[string]CacheItem), // Initialize empty map (required!)
		ttl:  ttl,                        // Store the TTL for all future items
	}
}

// Set adds or updates a value in the cache
//
// WHY Lock() (not RLock())?
// - We're WRITING to the map = need EXCLUSIVE access
// - Lock() blocks ALL other goroutines (readers AND writers)
// - This prevents data corruption from concurrent writes
//
// WHY defer Unlock()?
// - defer guarantees Unlock() runs even if panic occurs
// - Prevents deadlock (forgetting to unlock = program hangs)
// - Best practice: Lock + defer Unlock on same line
func (c *ConcurrentCache) Set(key string, value interface{}) {
	c.mu.Lock()         // Acquire EXCLUSIVE lock - blocks everyone
	defer c.mu.Unlock() // Release lock when function returns (even on panic)

	c.data[key] = CacheItem{
		Value:     value,                 // Store the actual value
		ExpiresAt: time.Now().Add(c.ttl), // Calculate expiration: now + TTL duration
	}
}

// Get retrieves a value from the cache
//
// WHY RLock() (not Lock())?
// - We're only READING = use SHARED lock
// - RLock() allows OTHER readers to access simultaneously
// - Much better performance when many goroutines read at once
//
// WHY return (interface{}, bool)?
// - interface{} = the cached value (caller must type-assert)
// - bool = "ok" pattern - tells if key exists AND is not expired
// - Idiomatic Go: val, ok := cache.Get("key")
func (c *ConcurrentCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()         // Acquire SHARED lock - allows other readers
	defer c.mu.RUnlock() // Release read lock when done

	// Try to find the item in the map
	item, exists := c.data[key] // map access returns (value, ok)
	if !exists {
		return nil, false // Key not found
	}

	// Check if expired (even if key exists, it might be stale)
	// time.Now().After(item.ExpiresAt) = "is current time past expiration?"
	if time.Now().After(item.ExpiresAt) {
		return nil, false // Item expired - treat as not found
	}

	return item.Value, true // Success! Return value and true
}

// Delete removes a key from the cache
//
// WHY Lock() (not RLock())?
// - delete() MODIFIES the map = need exclusive access
func (c *ConcurrentCache) Delete(key string) {
	c.mu.Lock() // Exclusive lock for write operation
	defer c.mu.Unlock()
	delete(c.data, key) // Built-in delete() removes key from map
}

// Size returns the number of items in the cache
//
// NOTE: This returns total items INCLUDING expired ones
// (expired items are only removed on Get() or CleanupExpired())
//
// WHY RLock()?
// - len() only READS the map = shared lock is sufficient
func (c *ConcurrentCache) Size() int {
	c.mu.RLock() // Shared lock - just reading
	defer c.mu.RUnlock()
	return len(c.data) // len() returns number of keys in map
}

// Clear removes all items from the cache
//
// WHY make(map[string]CacheItem) instead of delete loop?
// - Creating new map is O(1) - instant
// - Deleting each key would be O(n) - slow for large cache
// - Old map gets garbage collected automatically
func (c *ConcurrentCache) Clear() {
	c.mu.Lock() // Exclusive lock - we're replacing the map
	defer c.mu.Unlock()
	c.data = make(map[string]CacheItem) // Replace with fresh empty map
}

// CleanupExpired removes expired items from the cache
//
// WHY return int?
// - Tells caller how many items were removed (useful for logging/monitoring)
//
// WHY store time.Now() in variable?
// - Calling time.Now() repeatedly in loop would give different times
// - Using single 'now' ensures consistent expiration check
//
// WHY is it safe to delete during iteration?
// - Go allows deleting current key while iterating over map
// - This is explicitly supported in the Go spec
func (c *ConcurrentCache) CleanupExpired() int {
	c.mu.Lock() // Exclusive lock - modifying map
	defer c.mu.Unlock()

	count := 0                      // Track how many items we remove
	now := time.Now()               // Get current time ONCE for consistency
	for key, item := range c.data { // Iterate all items
		if now.After(item.ExpiresAt) { // Is this item expired?
			delete(c.data, key) // Remove expired item from map
			count++             // Increment removed counter
		}
	}
	return count // Return total removed (for logging/monitoring)
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
