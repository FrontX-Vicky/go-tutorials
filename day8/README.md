# Day 8 — Concurrency Testing in Go

## Concepts to Learn

### 1. Testing Concurrent Code
- Race conditions and data races
- Using `go test -race` to detect race conditions
- Testing goroutine behavior
- Synchronization primitives (Mutex, RWMutex, WaitGroup)

### 2. The `-race` Flag
```bash
go test -race ./...
```
- Detects data races at runtime
- Reports exact locations of race conditions
- Essential for concurrent code
- Slight performance overhead (10-20x slower)

**⚠️ Important Note for Windows Users**:
- The `-race` flag requires CGO (C compiler) to be installed
- If you don't have `gcc` installed, you'll get: `cgo: C compiler "gcc" not found`
- **Solution**: Run tests without `-race` flag - they work perfectly fine!
- Race conditions still produce incorrect results even without the detector
- The race detector just makes them easier to find

**To install CGO on Windows** (optional):
1. Install TDM-GCC: https://jmeubank.github.io/tdm-gcc/
2. Or install MinGW: http://www.mingw.org/
3. Add to PATH and restart terminal

**Alternative**: Use WSL (Windows Subsystem for Linux) for race detection

### 3. Testing with Goroutines
- Waiting for goroutines to complete
- Using channels for synchronization
- Testing concurrent access to shared state
- Verifying goroutine completion

Example:
```go
func TestConcurrentAccess(t *testing.T) {
    counter := NewSafeCounter()
    var wg sync.WaitGroup
    
    // Launch 100 goroutines
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    
    if counter.Value() != 100 {
        t.Errorf("expected 100, got %d", counter.Value())
    }
}
```

### 4. Channel Testing
- Testing channel sends and receives
- Testing channel timeouts
- Testing channel closing behavior
- Buffered vs unbuffered channels

Example:
```go
func TestChannelTimeout(t *testing.T) {
    ch := make(chan int)
    
    select {
    case <-ch:
        t.Error("should not receive")
    case <-time.After(100 * time.Millisecond):
        // Expected timeout
    }
}
```

### 5. Testing Worker Pools
- Multiple workers processing from queue
- Load balancing verification
- Graceful shutdown testing
- Error handling in workers

### 6. Mutex and Synchronization Testing
- Testing lock contention
- Verifying mutex correctness
- Deadlock detection
- Testing RWMutex read/write patterns

### 7. Context Cancellation Testing
- Testing context timeout behavior
- Testing context cancellation propagation
- Testing graceful shutdown
- Testing cleanup on cancellation

Example:
```go
func TestContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    done := make(chan bool)
    go func() {
        DoWork(ctx)
        done <- true
    }()
    
    cancel() // Cancel the context
    
    select {
    case <-done:
        // Work stopped as expected
    case <-time.After(1 * time.Second):
        t.Error("work did not stop after cancel")
    }
}
```

### 8. Testing Concurrent Data Structures
- Thread-safe maps
- Concurrent queues
- Lock-free data structures
- Atomic operations

---

## Today's Exercise

### Goal
Build and test concurrent data structures and patterns:
1. Thread-safe counter with atomic operations
2. Concurrent-safe cache with RWMutex
3. Worker pool with job queue
4. Pipeline pattern with channels
5. Fan-out/Fan-in patterns

### Test Requirements

1. **SafeCounter Tests**
   - Test concurrent increments
   - Test concurrent reads
   - Test with race detector
   - Verify final count accuracy

2. **ConcurrentCache Tests**
   - Test concurrent reads and writes
   - Test cache eviction
   - Test TTL expiration
   - Test with race detector

3. **WorkerPool Tests**
   - Test job processing
   - Test graceful shutdown
   - Test error handling
   - Test worker count

4. **Pipeline Tests**
   - Test data flow through stages
   - Test cancellation propagation
   - Test error handling
   - Test backpressure

5. **Fan-Out/Fan-In Tests**
   - Test parallel processing
   - Test result aggregation
   - Test timeout handling

---

## Key Testing Commands

```bash
# Run tests with race detector
go test -race -v

# Run specific test with race detector
go test -race -run TestConcurrentAccess -v

# Run benchmarks
go test -bench=. -benchmem

# Run with timeout
go test -timeout 30s -race -v

# Run with coverage
go test -race -cover -v
```

---

## Common Race Condition Patterns

### 1. Read-Modify-Write Race
```go
// UNSAFE - Race condition!
func (c *Counter) Increment() {
    c.value++  // Read, increment, write - not atomic
}

// SAFE - Using mutex
func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// SAFE - Using atomic
func (c *Counter) Increment() {
    atomic.AddInt64(&c.value, 1)
}
```

### 2. Check-Then-Act Race
```go
// UNSAFE - Race condition!
func (m *Map) SetIfAbsent(key, value string) bool {
    if _, exists := m.data[key]; !exists {
        m.data[key] = value  // Another goroutine might have set it
        return true
    }
    return false
}

// SAFE - Using mutex
func (m *Map) SetIfAbsent(key, value string) bool {
    m.mu.Lock()
    defer m.mu.Unlock()
    if _, exists := m.data[key]; !exists {
        m.data[key] = value
        return true
    }
    return false
}
```

### 3. Closure Variable Capture
```go
// UNSAFE - All goroutines see same i
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)  // Race condition!
    }()
}

// SAFE - Pass i as parameter
for i := 0; i < 10; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}
```

---

## Files to Create

1. `main.go` - Concurrent data structures implementation
2. `main_test.go` - Test suite (fill in TODOs)
3. `main_test_answers.go` - Reference implementations

---

## Tips

1. **Always run with `-race`** during development
2. **Use `sync.WaitGroup`** to wait for goroutines
3. **Prefer channels** for communication between goroutines
4. **Use `context`** for cancellation and timeouts
5. **Avoid sharing memory** - share by communicating instead
6. **Test with high concurrency** (100+ goroutines)
7. **Test edge cases**: empty, single item, many items
8. **Use `time.After`** for test timeouts

---

## Expected Output

```
$ go test -race -v
=== RUN   TestSafeCounter_Concurrent
--- PASS: TestSafeCounter_Concurrent (0.01s)
=== RUN   TestConcurrentCache_ReadWrite
--- PASS: TestConcurrentCache_ReadWrite (0.02s)
=== RUN   TestWorkerPool_ProcessJobs
--- PASS: TestWorkerPool_ProcessJobs (0.05s)
=== RUN   TestPipeline_DataFlow
--- PASS: TestPipeline_DataFlow (0.01s)
=== RUN   TestFanOutFanIn_Parallel
--- PASS: TestFanOutFanIn_Parallel (0.03s)
PASS
ok      go-tutorials/day8    0.12s
```

Good luck with Day 8! 🚀
