# Learning Log

## Day 1 — Go Basics
**Date**: November 17, 2025

**Concepts Covered**:
- Basic syntax and structure
- Variables and types
- Control flow (if, for, switch)
- Functions

**Files Created**:
- `day1/main.go`

**Key Takeaways**:
- Go's simple syntax and explicit error handling
- No semicolons needed
- Package system and imports

---

## Day 2 — Structs, Methods, Interfaces
**Date**: November 19, 2025

**Concepts Covered**:
- Structs as composite types
- Methods with value and pointer receivers
- Interfaces for polymorphism
- Struct composition (embedding)

**Files Created/Updated**:
- `day2/README.md`
- `day2/main.go`

**Key Takeaways**:
- Use pointer receivers when you need to modify the receiver or avoid copying
- Interfaces are implemented implicitly (no `implements` keyword)
- Composition over inheritance through embedding
- Interface satisfaction is checked at compile time

**Tasks Completed**:
- Created `User` struct with `Name`, `Age`, and `Skills` fields
- Implemented methods: `Greet()`, `IsAdult()`, `AddSkill()`
- Defined and implemented `Profile` interface
- Created `Employee` struct with embedded `User`
- Implemented `PrintProfiles()` helper for polymorphic behavior

---

## Day 3 — Error Handling, Pointers, Slices & Maps
**Date**: November 19-20, 2025

**Concepts Covered**:
- Error handling patterns (no exceptions)
- Pointers and nil checks
- Slice operations and behavior
- Map operations and iteration

**Files Created**:
- `day3/README.md`
- `day3/main.go`

**Key Takeaways**:
- Go uses error values instead of exceptions—always check `if err != nil`
- Must initialize maps with `make()` before use (nil maps panic on write)
- Map-based deduplication (O(n)) is more idiomatic than nested loops
- `strings.Fields()` is the standard way to split text into words
- `defer` + `recover()` can catch panics and convert them to errors

**Tasks Completed**:
- Built `UserRegistry` with CRUD operations (AddUser, GetUser, UpdateUser, DeleteUser, ListUsers)
- Implemented `Divide` and `Sqrt` functions with proper error handling
- Created `RemoveDuplicates` using map for O(n) performance
- Implemented `CountWords` using `strings.Fields()` and map
- Implemented `Safe` wrapper to handle panics gracefully

---

## Day 4 — Concurrency: Goroutines, Channels & Sync
**Date**: November 20, 2025

**Concepts Covered**:
- Goroutines (lightweight threads)
- Channels (unbuffered and buffered)
- Select statement for multiplexing
- WaitGroups for synchronization
- Concurrency patterns (worker pool, pipeline, fan-out/fan-in)

**Files Created**:
- `day4/README.md`
- `day4/main.go`

**Key Takeaways**:
- Prefer WaitGroups over sleeps to coordinate goroutines
- Directional channel types (`chan<-`, `<-chan`) clarify intent
- Buffered vs unbuffered channels change blocking behavior
- Use select with `time.After`/timers for timeouts
- Worker pool with `close` + `range` is an idiomatic pattern

**Tasks Completed**:
- Launch goroutines and coordinate with channels
- Understand buffered vs unbuffered channels
- Implement worker pool pattern
- Use select with timeout
- (Optional) Build a pipeline with multiple stages

---

## Day 5 — Context, HTTP Server, Graceful Shutdown, Middleware
**Date**: November 25, 2025

**Concepts Covered**:
- Contexts for cancellation and timeouts
- HTTP handlers and JSON encoding
- Middleware (logging, request timeout, rate limiting)
- Graceful shutdown with `Server.Shutdown`

**Files Created**:
- `day5/README.md`
- `day5/main.go`

**Key Takeaways**:
- Always propagate and honor `context.Context` in handlers and stores
- Protect shared state with `sync.RWMutex`
- Compose middlewares to add cross-cutting concerns cleanly
- Implement graceful shutdown to drain in-flight requests safely
- Token bucket algorithm for rate limiting
- `struct{}` for memory-efficient signaling channels

**Tasks Completed**:
- Built context-aware UserStore with RWMutex protection
- Implemented RESTful HTTP handlers (POST, GET, DELETE)
- Created middleware: Logging, RequestTimeout, RateLimiter
- Implemented graceful shutdown with signal handling
- Health check endpoint with context timeout

---

## Day 6 — Testing, Benchmarking, and Test Coverage
**Date**: November 28 - December 8, 2025

**Concepts Covered**:
- Go's built-in testing package
- Table-driven tests (idiomatic Go pattern)
- Test coverage measurement
- Benchmarking performance
- Testing HTTP handlers with httptest
- Testable code patterns

**Files Created**:
- `day6/README.md`
- `day6/main.go`
- `day6/main_test.go`
- `day6/main_test_answers.go` (answer sheet)

**Key Takeaways**:

1. **Table-Driven Tests**: The idiomatic Go pattern
   - Define test cases as a slice of structs
   - Use `t.Run()` for subtests and clean hierarchical output
   - Easy to scale—add new cases without code duplication

2. **Test Function Purpose**: `func TestXxx(t *testing.T)`
   - `t *testing.T` parameter gives access to testing methods
   - `t.Errorf()` - report failure
   - `t.Run()` - create subtests
   - `t.Fail()`, `t.Fatal()`, `t.Skip()` - test control

3. **Context Testing**:
   - Always test context cancellation paths
   - Use `context.WithCancel()` to cancel immediately
   - Context checks happen BEFORE slow operations
   - Test context.WithTimeout() for deadline scenarios

4. **Test Coverage**:
   - Achieved **94.1% coverage** (exceeded 80% goal)
   - Test success cases, error cases, and edge cases
   - Use `go test -cover` to measure
   - Coverage ≠ Quality, but helps identify untested paths

5. **Benchmarking**:
   - Use `b.N` loop—Go framework determines iterations
   - Run with `go test -bench=.`
   - Measures performance and detects regressions
   - `b.ReportAllocs()` shows memory allocations

6. **Helper Functions**:
   - `createTestUser()` - reduce test setup duplication
   - `createTimeoutContext()` - reusable context creation
   - Keep test code DRY

**Tasks Completed**:
- ✅ Wrote table-driven tests for `Add`, `IsPalindrome`, `Reverse`, `ValidateEmail`
- ✅ Tested UserStore CRUD (Create, Get, List, Delete) with edge cases
- ✅ Tested context cancellation for concurrent operations
- ✅ Wrote benchmarks for performance measurement
- ✅ Achieved 94.1% test coverage
- ✅ All 8 test functions + subtests passing (100% success)
- ✅ Created answer sheet with complete solution

**Test Results**:
- **Status**: ✅ PASS (all tests passing)
- **Coverage**: 94.1% of statements
- **Test Functions**: 8 (TestAdd, TestIsPalindrome, TestReverse, TestValidateEmail, TestUserStore_Create, TestUserStore_Get, TestUserStore_List, TestUserStore_Delete)
- **Subtests**: 16+ (organized with t.Run())
- **Benchmarks**: 4 (BenchmarkAdd, BenchmarkIsPalindrome, BenchmarkUserStore_Create, BenchmarkUserStore_Get)

**Code Quality**:
- ✅ All test functions follow `TestXxx` convention
- ✅ All benchmarks follow `BenchmarkXxx` convention
- ✅ Subtests organized with `t.Run(name, func(t *testing.T) {})`
- ✅ Clear error messages with `t.Errorf()`
- ✅ Edge cases covered: empty strings, zero values, not found errors, duplicates, context cancellation
- ✅ Idiomatic Go patterns throughout

**Grade**: 95/100

---

## Day 7 — HTTP Handler Testing
**Date**: December 8-9, 2025

**Concepts Covered**:
- HTTP handler testing with `httptest` package
- `httptest.NewRecorder()` for capturing responses
- `httptest.NewRequest()` for creating test requests
- Table-driven HTTP tests
- Middleware testing (Logging, Timeout, RateLimiter)
- Error response format testing
- Body size limiting with `http.MaxBytesReader`

**Files Created**:
- `day7/README.md`
- `day7/main.go`
- `day7/main_test.go`
- `day7/answers_test.go` (answer sheet)
- `day7/CODE_REVIEW.md`
- `day7/DATA_FLOW_EXPLANATION.md`

**Key Takeaways**:

1. **httptest Package**: Test HTTP handlers without a network
   - `httptest.NewRecorder()` captures the response
   - `httptest.NewRequest()` creates a request with method, path, body
   - Direct handler invocation for fast, isolated tests

2. **Table-Driven HTTP Tests**: Comprehensive endpoint testing
   - Test multiple HTTP methods (GET, POST, DELETE)
   - Test valid/invalid inputs
   - Test different status codes and response bodies

3. **Middleware Testing**: Test cross-cutting concerns
   - Logging middleware: verify log output
   - Timeout middleware: use `context.WithTimeout()`
   - Rate limiter: test request limiting behavior

4. **Body Size Limiting**: Protect against large payloads
   - `http.MaxBytesReader()` limits request body size
   - Returns 413 Payload Too Large when exceeded
   - Critical for production security

5. **Error Response Testing**: Consistent error format
   - JSON error responses with structured format
   - Test Content-Type headers
   - Verify error message content

**Tasks Completed**:
- ✅ Implemented REST API with User CRUD handlers
- ✅ Tested handleListUsers (empty list, with users)
- ✅ Tested handleCreateUser (valid, invalid JSON, missing fields, large payload)
- ✅ Tested handleGetUser (valid ID, not found, invalid ID)
- ✅ Tested handleDeleteUser (existing, not found, invalid ID)
- ✅ Table-driven tests for all HTTP endpoints
- ✅ Middleware testing (Logging, Timeout, RateLimiter)
- ✅ Error response format testing
- ✅ Body size limit implementation (5KB max)

**Test Results**:
- **Status**: ✅ PASS (24/24 tests passing)
- **Coverage**: Comprehensive HTTP handler coverage
- **Test Functions**: 8 test functions with multiple subtests
- **Benchmarks**: 2 (handler performance)

**Grade**: 100/100

---

## Day 8 — Concurrency Testing
**Date**: December 19, 2025

**Concepts to Learn**:
- Race condition detection with `-race` flag
- Testing concurrent access to shared data
- Atomic operations with `sync/atomic` package
- RWMutex for read-heavy concurrent access
- Worker pool pattern testing
- Channel and pipeline testing
- Semaphore pattern for limiting concurrency
- Fan-out/Fan-in patterns
- Context cancellation in concurrent code
- Goroutine leak detection

**Files Created**:
- `day8/README.md` - Learning guide
- `day8/main.go` - Concurrency implementations
- `day8/main_test.go` - Exercise tests (TODO)
- `day8/answers_test.go` - Answer key (reference)

**Your Tasks**:
1. Study `main.go` - understand SafeCounter, ConcurrentCache, WorkerPool, Pipeline, Semaphore
2. Complete the 14 TODO tests in `main_test.go`
3. Run tests with `go test -v` to verify your solutions
4. Compare with `answers_test.go` after completing

**Commands**:
```bash
# Run your exercise tests
go test -v

# Run answer tests to see expected behavior
go test -v -run "_Answer$"

# Run benchmarks
go test -bench=.
```

**Status**: ✅ COMPLETE (with one deadlock learning experience!)

**Tasks Completed**:
- ✅ TestSafeCounter_Basic
- ✅ TestSafeCounter_Concurrent
- ✅ TestConcurrentCache_Basic
- ✅ TestConcurrentCache_Concurrent
- ✅ TestConcurrentCache_TTL
- ✅ TestWorkerPool_Basic (found deadlock issue - learning opportunity!)
- ✅ TestWorkerPool_Concurrent
- ✅ TestPipeline_Basic
- ✅ TestPipeline_Cancellation
- ✅ TestSemaphore_Basic
- ✅ TestSemaphore_TryAcquire
- ✅ TestFanOutFanIn
- ✅ TestNoGoroutineLeak
- ✅ TestUnsafeCounter_Race (demonstrated race condition)

**Key Takeaways**:

1. **Race Detector**: Requires CGO on Windows - tests work without it
2. **Atomic Operations**: `atomic.AddInt64()` for lock-free counting
3. **RWMutex Optimization**: Multiple readers OR one writer pattern
4. **WorkerPool Pattern**: Use goroutine to submit jobs to avoid deadlock
5. **Channel Deadlock**: Never block on send while needing to receive from same system
6. **Context Cancellation**: Stops pipelines gracefully via `ctx.Done()`
7. **Semaphore Pattern**: Buffered channel limits concurrent operations
8. **Goroutine Leak Detection**: `runtime.NumGoroutine()` before/after

**Important Lesson - Deadlock in TestWorkerPool_Basic**:
- **Problem**: Submitting 300 jobs with 100-item buffer causes deadlock
- **Why**: Main thread blocks submitting job 101, never reaches result collection
- **Solution**: Submit jobs in goroutine, collect results in main thread
- **Key Insight**: Concurrent submission + collection prevents channel deadlock

**Grade**: 98/100 (Excellent work! Minor deduction for initial deadlock, but great learning!)

