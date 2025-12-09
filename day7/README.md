# Day 7 — Testing HTTP Handlers & Integration Testing

## Concepts to Learn

### 1. The `httptest` Package
- `httptest.NewRequest()` — create HTTP requests for testing
- `httptest.NewRecorder()` — capture HTTP responses without a real server
- `httptest.NewServer()` — start a real test server if needed
- No network overhead, fast and isolated tests

Example:
```go
func TestGetHandler(t *testing.T) {
    // Create a test request
    req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Call the handler
    handler := http.HandlerFunc(GetUserHandler)
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", rr.Code)
    }
}
```

### 2. Testing Different HTTP Methods
- **GET**: Retrieve data (no body)
- **POST**: Create data (with request body)
- **PUT/PATCH**: Update data
- **DELETE**: Remove data
- Test both success and error cases for each method

### 3. Request Body & JSON Parsing
- Use `strings.NewReader()` to create request body
- Parse JSON in tests with `json.Unmarshal()`
- Test malformed JSON rejection
- Test missing required fields

Example:
```go
func TestPostHandler(t *testing.T) {
    body := `{"name": "Alice", "age": 30}`
    req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()
    PostUserHandler(rr, req)
    
    if rr.Code != http.StatusCreated {
        t.Errorf("expected 201, got %d", rr.Code)
    }
}
```

### 4. Response Headers & Body Validation
- Check `Content-Type` header
- Parse response body (JSON, text, etc.)
- Verify header values
- Test redirect status codes (301, 302, 307)

### 5. Testing Middleware
- Chain middleware in tests same way as production
- Test middleware order (does it execute first?)
- Test error responses from middleware
- Test middleware side effects (logging, rate limiting, etc.)

### 6. Table-Driven HTTP Tests
- Combine table-driven pattern with HTTP testing
- Test multiple endpoints in one test
- Test multiple methods on same endpoint
- Easy to add new test cases

Example:
```go
func TestHTTPEndpoints(t *testing.T) {
    tests := []struct {
        name           string
        method         string
        path           string
        body           string
        expectedStatus int
    }{
        {"GET /users", http.MethodGet, "/users", "", http.StatusOK},
        {"POST /users valid", http.MethodPost, "/users", `{"name":"Bob"}`, http.StatusCreated},
        {"POST /users invalid", http.MethodPost, "/users", `{invalid json}`, http.StatusBadRequest},
        {"DELETE /users/notfound", http.MethodDelete, "/users/notfound", "", http.StatusNotFound},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
            rr := httptest.NewRecorder()
            
            // Call your router/handler here
            YourHTTPHandler(rr, req)
            
            if rr.Code != tt.expectedStatus {
                t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
            }
        })
    }
}
```

## Tasks

### 1. Test GET /users Endpoint
- Empty store (should return 200 with empty list)
- Populated store (should return 200 with user list)
- Test response JSON structure
- Test Content-Type header is `application/json`

### 2. Test POST /users Endpoint
- Valid user creation (should return 201)
- Invalid JSON (should return 400)
- Duplicate ID (should return 409)
- Missing required fields (should return 400)
- Test response contains created user data

### 3. Test GET /users/:id Endpoint
- Existing user (should return 200 with user data)
- Non-existent user (should return 404)
- Invalid ID format (should return 400)
- Test response includes correct user details

### 4. Test DELETE /users/:id Endpoint
- Existing user (should return 200)
- Non-existent user (should return 404)
- Invalid ID format (should return 400)
- Verify user is actually deleted

### 5. Test Middleware Integration
- Test that logging middleware executes
- Test that timeout middleware cancels long requests
- Test rate limiter blocks after threshold
- Verify middleware order is correct

### 6. Test Error Responses
- Malformed JSON returns 400
- Missing Content-Type header handling
- Server error responses (500)
- Test error message format

## Extra Challenge
- Test with HTTP client in real server (use `httptest.NewServer()`)
- Test concurrent requests to verify thread safety
- Mock external dependencies (database, API calls)
- Test WebSocket connections
- Test streaming responses
- Create a test helper that sets up common scenarios
- Test authentication and authorization
- Benchmark handler performance with `b.ReportAllocs()`

