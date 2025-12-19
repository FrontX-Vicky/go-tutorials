package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

// Day 8: Concurrency Testing - Exercises
// Run with: go test -race -v
// The -race flag enables the race detector to catch data races

// ============================================
// Test 1: SafeCounter - Basic Operations
// ============================================
// TODO: Test that SafeCounter correctly increments and decrements
// Test concurrent access to ensure no race conditions

func TestSafeCounter_Basic(t *testing.T) {
	// TODO: Implement this test
	// 1. Create a new SafeCounter
	// 2. Test Increment() increases value by 1
	// 3. Test Decrement() decreases value by 1
	// 4. Test Add() adds the correct amount
	// 5. Verify Value() returns correct result

	t.Skip("TODO: Implement TestSafeCounter_Basic")
}

// ============================================
// Test 2: SafeCounter - Concurrent Access
// ============================================
// TODO: Test concurrent increments from multiple goroutines
// This test will detect race conditions when run with -race flag

func TestSafeCounter_Concurrent(t *testing.T) {
	// TODO: Implement this test
	// 1. Create a SafeCounter
	// 2. Launch 100 goroutines that each increment 1000 times
	// 3. Use sync.WaitGroup to wait for all goroutines
	// 4. Verify final count is 100,000

	t.Skip("TODO: Implement TestSafeCounter_Concurrent")
}

// ============================================
// Test 3: ConcurrentCache - Basic Operations
// ============================================
// TODO: Test basic cache operations (Set, Get, Delete)

func TestConcurrentCache_Basic(t *testing.T) {
	// TODO: Implement this test
	// 1. Create cache with 5 second TTL
	// 2. Test Set() and Get() work correctly
	// 3. Test Get() returns false for non-existent keys
	// 4. Test Delete() removes keys
	// 5. Test Size() returns correct count

	t.Skip("TODO: Implement TestConcurrentCache_Basic")
}

// ============================================
// Test 4: ConcurrentCache - Concurrent Access
// ============================================
// TODO: Test concurrent read/write access to cache

func TestConcurrentCache_Concurrent(t *testing.T) {
	// TODO: Implement this test
	// 1. Create a cache
	// 2. Launch multiple goroutines that write different keys
	// 3. Launch multiple goroutines that read keys
	// 4. Use WaitGroup to coordinate
	// 5. Verify no race conditions (run with -race flag)

	t.Skip("TODO: Implement TestConcurrentCache_Concurrent")
}

// ============================================
// Test 5: ConcurrentCache - TTL Expiration
// ============================================
// TODO: Test that cache items expire after TTL

func TestConcurrentCache_TTL(t *testing.T) {
	// TODO: Implement this test
	// 1. Create cache with very short TTL (100ms)
	// 2. Set a value
	// 3. Verify Get() returns the value immediately
	// 4. Wait for TTL to expire
	// 5. Verify Get() returns false after expiration

	t.Skip("TODO: Implement TestConcurrentCache_TTL")
}

// ============================================
// Test 6: WorkerPool - Basic Processing
// ============================================
// TODO: Test that WorkerPool processes jobs correctly

func TestWorkerPool_Basic(t *testing.T) {
	// TODO: Implement this test
	// 1. Create a processor function that doubles the job ID
	// 2. Create a WorkerPool with 2 workers
	// 3. Start the pool
	// 4. Submit 5 jobs
	// 5. Collect results and verify correctness
	// 6. Stop the pool

	t.Skip("TODO: Implement TestWorkerPool_Basic")
}

// ============================================
// Test 7: WorkerPool - Concurrent Processing
// ============================================
// TODO: Test that multiple workers process jobs concurrently

func TestWorkerPool_Concurrent(t *testing.T) {
	// TODO: Implement this test
	// 1. Create a processor that takes 50ms per job
	// 2. Create WorkerPool with 4 workers
	// 3. Submit 4 jobs
	// 4. Measure time to complete all jobs
	// 5. Verify parallel processing (should be ~50-100ms, not 200ms)

	t.Skip("TODO: Implement TestWorkerPool_Concurrent")
}

// ============================================
// Test 8: Pipeline - Stage Processing
// ============================================
// TODO: Test pipeline processes data through stages correctly

func TestPipeline_Basic(t *testing.T) {
	// TODO: Implement this test
	// 1. Create a pipeline
	// 2. Add a Multiplier(2) stage
	// 3. Add a Filter(even numbers) stage
	// 4. Generate input [0,1,2,3,4]
	// 5. Run pipeline and collect output
	// 6. Verify output contains correct values

	t.Skip("TODO: Implement TestPipeline_Basic")
}

// ============================================
// Test 9: Pipeline - Context Cancellation
// ============================================
// TODO: Test that pipeline respects context cancellation

func TestPipeline_Cancellation(t *testing.T) {
	// TODO: Implement this test
	// 1. Create a context with cancel
	// 2. Create a pipeline with slow processing
	// 3. Start pipeline in goroutine
	// 4. Cancel context after short time
	// 5. Verify pipeline stops gracefully

	t.Skip("TODO: Implement TestPipeline_Cancellation")
}

// ============================================
// Test 10: Semaphore - Limit Concurrency
// ============================================
// TODO: Test that semaphore limits concurrent operations

func TestSemaphore_Basic(t *testing.T) {
	// TODO: Implement this test
	// 1. Create semaphore with limit 2
	// 2. Track active goroutines
	// 3. Launch 10 goroutines that:
	//    - Acquire semaphore
	//    - Increment active count
	//    - Verify active <= 2
	//    - Sleep briefly
	//    - Decrement and release
	// 4. Verify max concurrent was never exceeded

	t.Skip("TODO: Implement TestSemaphore_Basic")
}

// ============================================
// Test 11: Semaphore - TryAcquire
// ============================================
// TODO: Test TryAcquire non-blocking behavior

func TestSemaphore_TryAcquire(t *testing.T) {
	// TODO: Implement this test
	// 1. Create semaphore with limit 1
	// 2. Acquire the slot
	// 3. TryAcquire should return false
	// 4. Release the slot
	// 5. TryAcquire should return true

	t.Skip("TODO: Implement TestSemaphore_TryAcquire")
}

// ============================================
// Test 12: FanOut/FanIn - Parallel Processing
// ============================================
// TODO: Test fan-out/fan-in pattern

func TestFanOutFanIn(t *testing.T) {
	// TODO: Implement this test
	// 1. Create input channel with values
	// 2. FanOut to 3 workers that square the input
	// 3. FanIn results back to single channel
	// 4. Collect and verify results

	t.Skip("TODO: Implement TestFanOutFanIn")
}

// ============================================
// Test 13: Race Condition Detection (INTENTIONALLY BUGGY)
// ============================================
// This demonstrates what a race condition looks like
// Run with -race flag to see it detected

// UnsafeCounter has an intentional race condition for demo
type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Increment() {
	c.value++ // Race condition!
}

func (c *UnsafeCounter) Value() int {
	return c.value
}

func TestUnsafeCounter_Race(t *testing.T) {
	// This test demonstrates race detection
	// When run with -race, Go will detect the data race

	t.Skip("Skipped: Run manually with 'go test -race -run TestUnsafeCounter_Race' to see race detection")

	// Uncomment below to see race detection in action:
	/*
		counter := &UnsafeCounter{}
		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					counter.Increment()
				}
			}()
		}

		wg.Wait()
		t.Logf("Final value: %d (expected 100000, but may differ due to race)", counter.Value())
	*/
}

// ============================================
// Test 14: Goroutine Leak Detection
// ============================================
// TODO: Test that resources are properly cleaned up

func TestNoGoroutineLeak(t *testing.T) {
	// TODO: Implement this test
	// 1. Record initial goroutine count (runtime.NumGoroutine())
	// 2. Create and use WorkerPool
	// 3. Properly stop the pool
	// 4. Wait briefly for cleanup
	// 5. Verify goroutine count returned to initial

	t.Skip("TODO: Implement TestNoGoroutineLeak")
}

// ============================================
// Benchmark: SafeCounter vs UnsafeCounter
// ============================================

func BenchmarkSafeCounter_Increment(b *testing.B) {
	counter := NewSafeCounter()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkConcurrentCache_Set(b *testing.B) {
	cache := NewConcurrentCache(5 * time.Minute)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set("key", i)
			i++
		}
	})
}

func BenchmarkConcurrentCache_Get(b *testing.B) {
	cache := NewConcurrentCache(5 * time.Minute)
	cache.Set("key", "value")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get("key")
		}
	})
}

// ============================================
// Helper: Wait for condition with timeout
// ============================================

func waitFor(t *testing.T, timeout time.Duration, condition func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("Timeout waiting for: %s", msg)
}

// Dummy usage to avoid unused import errors
var _ = sync.WaitGroup{}
var _ = context.Background()
