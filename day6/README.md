# Day 6 — Testing, Benchmarking, and Test Coverage

## Concepts to Learn

### 1. Testing Fundamentals
- Go's built-in `testing` package—no external frameworks needed
- Test file naming: `filename_test.go`
- Test function signature: `func TestXxx(t *testing.T)`
- Running tests: `go test`, `go test -v`, `go test ./...`

### 2. Table-Driven Tests
- Most idiomatic Go testing pattern
- Define test cases as a slice of structs
- Loop through cases, run same logic
- Easy to add new cases without duplicating code

Example:
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -1, -2},
        {"zero", 0, 5, 5},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```


### 3. Test Coverage
- Measure how much of your code is exercised by tests
- Commands:
  - `go test -cover` — show coverage percentage
  - `go test -coverprofile=coverage.out` — generate coverage report
  - `go tool cover -html=coverage.out` — view in browser

### 4. Benchmarking
- Measure performance of your code
- Benchmark function signature: `func BenchmarkXxx(b *testing.B)`
- Run with: `go test -bench=.`
- The `b.N` loop: framework determines how many iterations to run

Example:
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(10, 20)
    }
}
```

### 5. Testable Code Patterns
- **Dependency injection**: Pass dependencies as parameters or interfaces
- **Interface-based design**: Mock implementations for testing
- **Pure functions**: No side effects, easier to test
- **Small, focused functions**: One responsibility per function

### 6. Testing HTTP Handlers
- Use `httptest.NewRecorder()` to capture responses
- Use `httptest.NewRequest()` to create test requests
- No need to start real server

Example:
```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/users", nil)
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(YourHandler)
    handler.ServeHTTP(rr, req)
    
    if rr.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rr.Code)
    }
}
```

## Tasks

### 1. Write Unit Tests for Helper Functions
Create functions and write table-driven tests:
- `Add(a, b int) int` — addition
- `IsPalindrome(s string) bool` — check if string is palindrome
- `Reverse(s string) string` — reverse a string
- `ValidateEmail(email string) error` — email validation

### 2. Test UserStore from Day 5
Write table-driven tests for:
- `Create` (success, duplicate error, context cancellation)
- `Get` (success, not found, context cancellation)
- `List` (empty store, multiple users, context cancellation)
- `Delete` (success, not found, context cancellation)

### 3. Test HTTP Handlers
Write tests for handlers from Day 5:
- POST `/users` — valid input, invalid JSON, missing fields, duplicate ID
- GET `/users` — empty list, multiple users
- GET `/users/{id}` — existing user, not found
- DELETE `/users/{id}` — success, not found

### 4. Write Benchmarks
Benchmark performance:
- `BenchmarkUserStoreCreate`
- `BenchmarkUserStoreGet`
- `BenchmarkIsPalindrome`

### 5. Achieve High Test Coverage
- Run `go test -cover` and aim for >80% coverage
- Identify untested code paths
- Add tests for edge cases

## Extra Challenge
- Mock external dependencies (e.g., database, HTTP client)
- Use `t.Parallel()` to run tests concurrently
- Write a custom test helper function
- Test error conditions and panics using `defer` + `recover`
- Generate coverage report and view in browser

## Key Takeaways

### ✅ Day 6 Completed Successfully

**Test Results**:
- **All Tests Passing**: ✅ 100% success rate
- **Test Coverage**: 94.1% of statements
- **Test Functions**: 8 test functions with subtests
- **Benchmarks**: 4 benchmark functions

**Important Concepts Mastered**:

1. **Table-Driven Tests**: The idiomatic Go pattern for testing
   - Define test cases as a slice of structs
   - Use `t.Run()` for subtests and clean output
   - Easy to scale and maintain

2. **Context Handling in Tests**:
   - Always test context cancellation paths
   - Use `context.WithCancel()` to immediately cancel
   - Test timeouts with `context.WithTimeout()`
   - Remember: context checks happen BEFORE slow operations

3. **Test Coverage**:
   - Aim for >80% coverage (you achieved 94.1%!)
   - Test success cases, error cases, and edge cases
   - Coverage ≠ Quality, but it helps ensure completeness

4. **Benchmarking**:
   - Use `b.N` loop—framework determines iterations
   - Run with `go test -bench=.`
   - Measure performance and detect regressions

5. **Helper Functions**:
   - `createTestUser()` — reduce duplication
   - `createTimeoutContext()` — reusable test setup
   - Keep test code DRY (Don't Repeat Yourself)

**Best Practices Applied**:
- ✅ All test function names follow `TestXxx` convention
- ✅ All benchmarks follow `BenchmarkXxx` convention
- ✅ Subtests organized with `t.Run()`
- ✅ Clear error messages with `t.Errorf()`
- ✅ Edge cases covered (empty strings, zero values, not found, duplicates)
- ✅ Context cancellation tested for concurrent operations

**Next Steps**:
- Move to Day 7 to learn about package management and modules
- Or explore advanced testing patterns (mocking, interfaces)
- Practice writing tests BEFORE implementation (TDD)
