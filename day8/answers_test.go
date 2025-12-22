package main

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Day 8: Concurrency Testing - Answer Key
// Reference implementations for all exercises

// ============================================
// Test 1: SafeCounter - Basic Operations
// ============================================

func TestSafeCounter_Basic_Answer(t *testing.T) {
	counter := NewSafeCounter()

	// Test initial value
	if counter.Value() != 0 {
		t.Errorf("Initial value should be 0, got %d", counter.Value())
	}

	// Test Increment
	counter.Increment()
	if counter.Value() != 1 {
		t.Errorf("After Increment, expected 1, got %d", counter.Value())
	}

	// Test Decrement
	counter.Decrement()
	if counter.Value() != 0 {
		t.Errorf("After Decrement, expected 0, got %d", counter.Value())
	}

	// Test Add
	counter.Add(10)
	if counter.Value() != 10 {
		t.Errorf("After Add(10), expected 10, got %d", counter.Value())
	}

	counter.Add(-5)
	if counter.Value() != 5 {
		t.Errorf("After Add(-5), expected 5, got %d", counter.Value())
	}
}

// ============================================
// Test 2: SafeCounter - Concurrent Access
// ============================================

func TestSafeCounter_Concurrent_Answer(t *testing.T) {
	counter := NewSafeCounter()
	var wg sync.WaitGroup

	numGoroutines := 100
	incrementsPerGoroutine := 1000

	// Launch goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment()
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	expected := int64(numGoroutines * incrementsPerGoroutine)
	if counter.Value() != expected {
		t.Errorf("Expected %d, got %d", expected, counter.Value())
	}
}

// ============================================
// Test 3: ConcurrentCache - Basic Operations
// ============================================

func TestConcurrentCache_Basic_Answer(t *testing.T) {
	cache := NewConcurrentCache(5 * time.Second)

	// Test Set and Get
	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")
	if !ok {
		t.Error("Expected key1 to exist")
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	// Test non-existent key
	_, ok = cache.Get("nonexistent")
	if ok {
		t.Error("Expected nonexistent key to return false")
	}

	// Test Size
	cache.Set("key2", "value2")
	if cache.Size() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Size())
	}

	// Test Delete
	cache.Delete("key1")
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Expected key1 to be deleted")
	}
	if cache.Size() != 1 {
		t.Errorf("Expected size 1 after delete, got %d", cache.Size())
	}

	// Test Clear
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", cache.Size())
	}
}

// ============================================
// Test 4: ConcurrentCache - Concurrent Access
// ============================================

func TestConcurrentCache_Concurrent_Answer(t *testing.T) {
	cache := NewConcurrentCache(5 * time.Second)
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + (i % 26)))
			cache.Set(key, i)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + (i % 26)))
			cache.Get(key)
		}(i)
	}

	wg.Wait()

	// If we get here without race detector complaints, the test passes
	t.Log("Concurrent access test passed - no race conditions detected")
}

// ============================================
// Test 5: ConcurrentCache - TTL Expiration
// ============================================

func TestConcurrentCache_TTL_Answer(t *testing.T) {
	ttl := 100 * time.Millisecond
	cache := NewConcurrentCache(ttl)

	cache.Set("key", "value")

	// Should exist immediately
	val, ok := cache.Get("key")
	if !ok || val != "value" {
		t.Error("Value should exist immediately after setting")
	}

	// Wait for TTL to expire
	time.Sleep(ttl + 50*time.Millisecond)

	// Should be expired now
	_, ok = cache.Get("key")
	if ok {
		t.Error("Value should have expired")
	}
}

// ============================================
// Test 6: WorkerPool - Basic Processing
// ============================================

func TestWorkerPool_Basic_Answer(t *testing.T) {
	processor := func(job Job) Result {
		return Result{
			JobID:  job.ID,
			Output: job.ID * 2,
		}
	}

	pool := NewWorkerPool(2, processor)
	pool.Start()

	// Submit jobs
	numJobs := 5
	for i := 0; i < numJobs; i++ {
		pool.Submit(Job{ID: i, Payload: nil})
	}

	// Collect results
	results := make(map[int]int)
	timeout := time.After(2 * time.Second)

	for len(results) < numJobs {
		select {
		case result := <-pool.Results():
			results[result.JobID] = result.Output.(int)
		case <-timeout:
			t.Fatal("Timeout waiting for results")
		}
	}

	pool.Stop()

	// Verify results
	for i := 0; i < numJobs; i++ {
		expected := i * 2
		if results[i] != expected {
			t.Errorf("Job %d: expected %d, got %d", i, expected, results[i])
		}
	}
}

// ============================================
// Test 7: WorkerPool - Concurrent Processing
// ============================================

func TestWorkerPool_Concurrent_Answer(t *testing.T) {
	processingTime := 50 * time.Millisecond
	processor := func(job Job) Result {
		time.Sleep(processingTime)
		return Result{JobID: job.ID, Output: "done"}
	}

	numWorkers := 4
	pool := NewWorkerPool(numWorkers, processor)
	pool.Start()

	numJobs := 4
	start := time.Now()

	// Submit jobs
	for i := 0; i < numJobs; i++ {
		pool.Submit(Job{ID: i})
	}

	// Collect results
	timeout := time.After(2 * time.Second)
	collected := 0
	for collected < numJobs {
		select {
		case <-pool.Results():
			collected++
		case <-timeout:
			t.Fatal("Timeout waiting for results")
		}
	}

	elapsed := time.Since(start)
	pool.Stop()

	// With 4 workers processing 4 jobs simultaneously,
	// it should take ~50-100ms, not 200ms (4 * 50ms)
	maxExpected := processingTime * 2 // Allow some overhead
	if elapsed > maxExpected {
		t.Errorf("Expected parallel processing in ~%v, took %v", processingTime, elapsed)
	}
}

// ============================================
// Test 8: Pipeline - Stage Processing
// ============================================

func TestPipeline_Basic_Answer(t *testing.T) {
	ctx := context.Background()

	pipeline := NewPipeline()
	pipeline.AddStage(Multiplier(2)) // Multiply by 2

	// Input: 0, 1, 2, 3, 4
	input := Generator(ctx, 5)

	// Run pipeline
	output := pipeline.Run(ctx, input)

	// Collect results
	var results []int
	for val := range output {
		results = append(results, val)
	}

	// Verify: should be 0, 2, 4, 6, 8
	expected := []int{0, 2, 4, 6, 8}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

// ============================================
// Test 9: Pipeline - Context Cancellation
// ============================================

func TestPipeline_Cancellation_Answer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	pipeline := NewPipeline()
	pipeline.AddStage(func(ctx context.Context, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for val := range in {
				select {
				case out <- val:
					time.Sleep(50 * time.Millisecond) // Slow processing
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	})

	// Generate many values
	input := Generator(ctx, 100)
	output := pipeline.Run(ctx, input)

	// Collect some results, then cancel
	collected := 0
	go func() {
		for range output {
			collected++
		}
	}()

	// Cancel after brief time
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Wait for goroutines to clean up
	time.Sleep(100 * time.Millisecond)

	// Should have collected only some values, not all 100
	if collected >= 100 {
		t.Error("Expected pipeline to stop early due to cancellation")
	}
	t.Logf("Collected %d values before cancellation", collected)
}

// ============================================
// Test 10: Semaphore - Limit Concurrency
// ============================================

func TestSemaphore_Basic_Answer(t *testing.T) {
	limit := 2
	sem := NewSemaphore(limit)

	var (
		active    int32
		maxActive int32
		wg        sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			// Track active count
			current := atomic.AddInt32(&active, 1)
			for {
				old := atomic.LoadInt32(&maxActive)
				if current <= old || atomic.CompareAndSwapInt32(&maxActive, old, current) {
					break
				}
			}

			// Do some work
			time.Sleep(50 * time.Millisecond)

			atomic.AddInt32(&active, -1)
		}()
	}

	wg.Wait()

	if maxActive > int32(limit) {
		t.Errorf("Max concurrent exceeded limit: got %d, limit was %d", maxActive, limit)
	}
	t.Logf("Max concurrent operations: %d (limit: %d)", maxActive, limit)
}

// ============================================
// Test 11: Semaphore - TryAcquire
// ============================================

func TestSemaphore_TryAcquire_Answer(t *testing.T) {
	sem := NewSemaphore(1)

	// First acquire should succeed
	if !sem.TryAcquire() {
		t.Error("First TryAcquire should succeed")
	}

	// Second should fail (non-blocking)
	if sem.TryAcquire() {
		t.Error("Second TryAcquire should fail")
	}

	// Release and try again
	sem.Release()

	if !sem.TryAcquire() {
		t.Error("TryAcquire after release should succeed")
	}

	sem.Release()
}

// ============================================
// Test 12: FanOut/FanIn - Parallel Processing
// ============================================

func TestFanOutFanIn_Answer(t *testing.T) {
	ctx := context.Background()

	// Create input channel
	input := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		input <- i
	}
	close(input)

	// Fan out to 3 workers that square the input
	worker := func(n int) int {
		return n * n
	}
	outputs := FanOut(ctx, input, 3, worker)

	// Fan in results
	results := FanIn(ctx, outputs...)

	// Collect results
	var collected []int
	for val := range results {
		collected = append(collected, val)
	}

	// Verify we got correct squares (order may vary)
	expected := map[int]bool{1: true, 4: true, 9: true, 16: true, 25: true}
	if len(collected) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(collected))
	}

	for _, v := range collected {
		if !expected[v] {
			t.Errorf("Unexpected result: %d", v)
		}
	}
}

// ============================================
// Test 14: Goroutine Leak Detection
// ============================================

func TestNoGoroutineLeak_Answer(t *testing.T) {
	// Record initial goroutine count
	initialCount := runtime.NumGoroutine()

	// Create and use worker pool
	processor := func(job Job) Result {
		return Result{JobID: job.ID, Output: "done"}
	}
	pool := NewWorkerPool(4, processor)
	pool.Start()

	// Submit some jobs
	for i := 0; i < 10; i++ {
		pool.Submit(Job{ID: i})
	}

	// Collect results
	timeout := time.After(2 * time.Second)
	collected := 0
	for collected < 10 {
		select {
		case <-pool.Results():
			collected++
		case <-timeout:
			t.Fatal("Timeout")
		}
	}

	// Stop pool
	pool.Stop()

	// Wait for cleanup
	time.Sleep(100 * time.Millisecond)

	// Force garbage collection
	runtime.GC()

	// Check goroutine count
	finalCount := runtime.NumGoroutine()

	// Allow small variance (1-2 goroutines for GC, etc)
	if finalCount > initialCount+2 {
		t.Errorf("Possible goroutine leak: initial=%d, final=%d", initialCount, finalCount)
	}
	t.Logf("Goroutines: initial=%d, final=%d", initialCount, finalCount)
}

// ============================================
// Test: Mutex vs RWMutex Performance
// ============================================

func TestMutexVsRWMutex_Performance_Answer(t *testing.T) {
	// This test demonstrates when RWMutex is beneficial
	// (when reads greatly outnumber writes)

	type MutexMap struct {
		data map[string]int
		mu   sync.Mutex
	}

	type RWMutexMap struct {
		data map[string]int
		mu   sync.RWMutex
	}

	// Setup
	mutexMap := &MutexMap{data: make(map[string]int)}
	rwMutexMap := &RWMutexMap{data: make(map[string]int)}

	// Populate initial data
	for i := 0; i < 100; i++ {
		key := string(rune('a' + i%26))
		mutexMap.data[key] = i
		rwMutexMap.data[key] = i
	}

	// Test with mostly reads (90% read, 10% write)
	numOps := 10000
	var wg sync.WaitGroup

	// Mutex version
	start := time.Now()
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		if i%10 == 0 {
			// Write
			go func(i int) {
				defer wg.Done()
				mutexMap.mu.Lock()
				mutexMap.data["test"] = i
				mutexMap.mu.Unlock()
			}(i)
		} else {
			// Read
			go func() {
				defer wg.Done()
				mutexMap.mu.Lock()
				_ = mutexMap.data["test"]
				mutexMap.mu.Unlock()
			}()
		}
	}
	wg.Wait()
	mutexTime := time.Since(start)

	// RWMutex version
	start = time.Now()
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		if i%10 == 0 {
			// Write
			go func(i int) {
				defer wg.Done()
				rwMutexMap.mu.Lock()
				rwMutexMap.data["test"] = i
				rwMutexMap.mu.Unlock()
			}(i)
		} else {
			// Read
			go func() {
				defer wg.Done()
				rwMutexMap.mu.RLock()
				_ = rwMutexMap.data["test"]
				rwMutexMap.mu.RUnlock()
			}()
		}
	}
	wg.Wait()
	rwMutexTime := time.Since(start)

	t.Logf("Mutex time: %v", mutexTime)
	t.Logf("RWMutex time: %v", rwMutexTime)
	t.Logf("RWMutex is %.2fx faster for read-heavy workload",
		float64(mutexTime)/float64(rwMutexTime))
}

// ============================================
// Test: Context Timeout Pattern
// ============================================

func TestContextTimeout_Answer(t *testing.T) {
	// Create a slow operation
	slowOperation := func(ctx context.Context) (string, error) {
		select {
		case <-time.After(1 * time.Second):
			return "completed", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	// Test with timeout shorter than operation
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := slowOperation(ctx)
	if err == nil {
		t.Error("Expected timeout error")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded, got %v", err)
	}

	// Test with timeout longer than operation
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	result, err := slowOperation(ctx2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "completed" {
		t.Errorf("Expected 'completed', got %s", result)
	}
}

// ============================================
// Test: Channel Select Patterns
// ============================================

func TestChannelSelectPatterns_Answer(t *testing.T) {
	t.Run("Non-blocking send", func(t *testing.T) {
		ch := make(chan int, 1)
		ch <- 1 // Fill the buffer

		// Non-blocking send
		select {
		case ch <- 2:
			t.Error("Should not have sent")
		default:
			// Expected - channel is full
		}
	})

	t.Run("Non-blocking receive", func(t *testing.T) {
		ch := make(chan int, 1)

		// Non-blocking receive from empty channel
		select {
		case <-ch:
			t.Error("Should not have received")
		default:
			// Expected - channel is empty
		}
	})

	t.Run("Multi-channel select", func(t *testing.T) {
		ch1 := make(chan int, 1)
		ch2 := make(chan int, 1)

		ch2 <- 42

		var result int
		select {
		case result = <-ch1:
		case result = <-ch2:
		}

		if result != 42 {
			t.Errorf("Expected 42, got %d", result)
		}
	})

	t.Run("Timeout with select", func(t *testing.T) {
		ch := make(chan int)

		select {
		case <-ch:
			t.Error("Should have timed out")
		case <-time.After(50 * time.Millisecond):
			// Expected - timeout
		}
	})
}
