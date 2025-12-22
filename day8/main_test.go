package main


import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	// "golang.org/x/text/currency"
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
	counter := NewSafeCounter()
	// 2. Test Increment() increases value by 1
	counter.Increment()
	if counter.Value() != 1 {
		t.Errorf("After increment, expected 1, got %d", counter.Value())
	}
	// 3. Test Decrement() decreases value by 1
	counter.Decrement()
	if counter.Value() != 0 {
		t.Errorf("After decrement, expected 0, got %d", counter.Value())
	}
	// 4. Test Add() adds the correct amount
	counter.Add(5)
	if counter.Value() != 5 {
		t.Errorf("After add(5), expected 5, got %d", counter.Value())
	}

	counter.Add(-3)
	if counter.Value() != 2 {
		t.Errorf("After add(-3), expected 2, got %d", counter.Value())
	}
	// 5. Verify Value() returns correct result



	t.Skip("TODO: Implement TestSafeCounter_Basic") // why this line has used here?
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

	counter := NewSafeCounter()
	var wg sync.WaitGroup

	numGoroutines := 100
	increamentsPerGoroutine := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increamentsPerGoroutine; j++ {
				counter.Increment()
				fmt.Printf("%d\n", counter.Value())
			}
		}()
	}

	wg.Wait()

	expected := int64(numGoroutines * increamentsPerGoroutine)

	if counter.Value() != expected {
		t.Errorf("Expected %d, got %d", expected, counter.Value())
	}

	// t.Skip("TODO: Implement TestSafeCounter_Concurrent")
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

	cache := NewConcurrentCache(5 * time.Second)

	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")
	if !ok {
		t.Errorf("Expected key1 to exist")
	}

	if val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	_, ok = cache.Get("nonexistent")
	if ok {
		t.Errorf("expected non existent key to return false")
	}

	cache.Set("key2", "value2")
	if cache.Size() != 2 {
		t.Errorf("Expected 2, got %d", cache.Size())
	}

	cache.Delete("key1")
	_, ok = cache.Get("key1")
	if ok {
		t.Errorf("Expected key1 to be deleted")
	}
	if cache.Size() != 1 {
		t.Errorf("Expected size 1 after delete, got %d", cache.Size())
	}

	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", cache.Size())
	}
	//  t.Skip("TODO: Implement TestConcurrentCache_Basic")
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

	cache := NewConcurrentCache(5 * time.Second)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + (i % 26)))
			cache.Set(key, i)
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + (i % 26)))
			_, _ = cache.Get(key)
		}(i)
	}

	wg.Wait()

	// t.Skip("TODO: Implement TestConcurrentCache_Concurrent")
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

	ttl := 100 * time.Millisecond

	cache := NewConcurrentCache(ttl)
	cache.Set("key", "value")

	val, ok := cache.Get("key")
	if !ok || val != "value" {
		t.Errorf("value should exist immediately after setting")
	}

	time.Sleep(ttl + 50*time.Millisecond)

	_, ok = cache.Get("key")
	if ok {
		t.Errorf("value should have expired after TTL")
	}

	// t.Skip("TODO: Implement TestConcurrentCache_TTL")
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

	processor := func(job Job) Result { // why processor funtion declared like this?
		return Result{
			JobID: job.ID,
			Output: job.ID * 2,
		}
	}

	pool := NewWorkerPool(2, processor)
	pool.Start()

	numJobs := 5
	for i := 0; i < numJobs; i++ {
		pool.Submit(Job{ID: i, Payload: nil})
	}

	results := make(map[int]int)
	timeout := time.After(2 * time.Second)

	for len(results) < numJobs {
		select {
		case result := <-pool.Results():
			results[result.JobID] = result.Output.(int)
		case <-timeout:
			t.Fatal("Timeout ")
		}
	}

	pool.Stop()

	for i := 0; i < numJobs; i++ {
		expected := i * 2
		if results[i] != expected {
			t.Errorf("Job %d: expected %d, got %d", i, expected, results[i])
		}
	}

	// t.Skip("TODO: Implement TestWorkerPool_Basic")
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

	processingTime := 50 * time.Millisecond
	processor := func(job Job) Result {
		time.Sleep(processingTime)
		return Result{JobID: job.ID, Output: "done"}
	}

	numWorkers := 4
	pool := NewWorkerPool(numWorkers, processor)
	pool.Start() // how this will work and why Start is written here?

	numJobs := 4
	start := time.Now()

	for i := 0; i < numJobs; i++ {
		pool.Submit(Job{ID: i}) // what submit will do?
	}

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

	maxExpected := processingTime * 2
	if elapsed > maxExpected {
		t.Errorf("Expected parallel processing in ~%v, took %v", processingTime, elapsed)
	}



	// t.Skip("TODO: Implement TestWorkerPool_Concurrent")
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

	ctx := context.Background()

	pipeline := NewPipeline()
	pipeline.AddStage(Multiplier(2))

	input := Generator(ctx, 5)

	output := pipeline.Run(ctx, input)

	var results []int
	for val := range output {
		results = append(results, val)
	}

	expected := []int{0, 2, 4, 6, 8}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Index %d: expected %d, got %d", i, expected[i], v)
		}
	}



	// t.Skip("TODO: Implement TestPipeline_Basic")
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

	ctx, cancel := context.WithCancel(context.Background())

	pipeline := NewPipeline() // what this function actually does
	pipeline.AddStage(func(ctx context.Context, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for val := range in {
				select {
				case out <- val:
					time.Sleep(50 * time.Millisecond)
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	})
	
	input := Generator(ctx, 100)
	output := pipeline.Run(ctx, input)

	collected := 0
	go func() {
		for range output {
			collected++
		}
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	time.Sleep(100 * time.Millisecond)

	if collected >= 100 {
		t.Errorf("Expected pipeline to stop early due to cancellation")
	}
	t.Logf("Collected %d values before cancellation", collected)
	

	// t.Skip("TODO: Implement TestPipeline_Cancellation")
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

	limit := 2
	sem := NewSemaphore(limit)

	var (
		active int32
		maxActive int32
		wg sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			current := atomic.AddInt32(&active, 1)
			for {
				old := atomic.LoadInt32(&maxActive)
				if current <= old || atomic.CompareAndSwapInt32(&maxActive, old, current) {
					break
				}
			}

			time.Sleep(50 * time.Millisecond)

			atomic.AddInt32(&active, -1)
		}()
	}

	wg.Wait()

	if maxActive > int32(limit) {
		t.Errorf("Max concurrent exceeded limit: got %d, limit was %d", maxActive, limit)
	}
	t.Logf("Max concurrent operations: %d (limit: %d)", maxActive, limit)

	// t.Skip("TODO: Implement TestSemaphore_Basic")
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
