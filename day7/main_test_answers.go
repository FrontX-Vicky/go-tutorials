package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Answer Sheet for Day 7 - HTTP Handler Testing

// TestHandleListUsers_EmptyStore
func TestHandleListUsers_EmptyStore_Answer(t *testing.T) {
	store := NewSimpleUserStore()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler := handleListUsers(store)
	handler.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Check Content-Type header
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	// Parse response
	var users []User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	// Should be empty array or nil
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

// TestHandleListUsers_WithUsers
func TestHandleListUsers_WithUsers_Answer(t *testing.T) {
	store := NewSimpleUserStore()
	store.Create(User{ID: "1", Name: "Alice", Age: 30})
	store.Create(User{ID: "2", Name: "Bob", Age: 25})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler := handleListUsers(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var users []User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

// TestHandleCreateUser_Success
func TestHandleCreateUser_Success_Answer(t *testing.T) {
	store := NewSimpleUserStore()
	body := `{"id":"1","name":"Alice","age":30}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handleCreateUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	// Verify user is in store
	user, err := store.Get("1")
	if err != nil {
		t.Errorf("user not found in store: %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", user.Name)
	}
}

// TestHandleCreateUser_InvalidJSON
func TestHandleCreateUser_InvalidJSON_Answer(t *testing.T) {
	store := NewSimpleUserStore()
	body := `{invalid json}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handleCreateUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

// TestHandleCreateUser_MissingFields
func TestHandleCreateUser_MissingFields_Answer(t *testing.T) {
	store := NewSimpleUserStore()

	tests := []struct {
		name string
		body string
	}{
		{"missing id", `{"name":"Alice","age":30}`},
		{"missing name", `{"id":"1","age":30}`},
		{"empty id", `{"id":"","name":"Alice","age":30}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler := handleCreateUser(store)
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", rr.Code)
			}
		})
	}
}

// TestHandleCreateUser_DuplicateID
func TestHandleCreateUser_DuplicateID_Answer(t *testing.T) {
	store := NewSimpleUserStore()
	body := `{"id":"1","name":"Alice","age":30}`

	// First creation should succeed
	req1 := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()

	handler := handleCreateUser(store)
	handler.ServeHTTP(rr1, req1)

	if rr1.Code != http.StatusCreated {
		t.Errorf("first create failed: expected 201, got %d", rr1.Code)
	}

	// Second creation with same ID should fail with 409
	req2 := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()

	handler.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusConflict {
		t.Errorf("duplicate create: expected 409, got %d", rr2.Code)
	}
}

// TestHandleGetUser_Success
func TestHandleGetUser_Success_Answer(t *testing.T) {
	store := NewSimpleUserStore()
	store.Create(User{ID: "1", Name: "Alice", Age: 30})

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rr := httptest.NewRecorder()

	handler := handleGetUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var user User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if user.Name != "Alice" || user.Age != 30 {
		t.Errorf("unexpected user data: %+v", user)
	}
}

// TestHandleGetUser_NotFound
func TestHandleGetUser_NotFound_Answer(t *testing.T) {
	store := NewSimpleUserStore()

	req := httptest.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	rr := httptest.NewRecorder()

	handler := handleGetUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

// TestHandleDeleteUser_Success
func TestHandleDeleteUser_Success_Answer(t *testing.T) {
	store := NewSimpleUserStore()
	store.Create(User{ID: "1", Name: "Alice", Age: 30})

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rr := httptest.NewRecorder()

	handler := handleDeleteUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Verify user is deleted
	_, err := store.Get("1")
	if err == nil {
		t.Errorf("user should be deleted but still exists")
	}
}

// TestHandleDeleteUser_NotFound
func TestHandleDeleteUser_NotFound_Answer(t *testing.T) {
	store := NewSimpleUserStore()

	req := httptest.NewRequest(http.MethodDelete, "/users/nonexistent", nil)
	rr := httptest.NewRecorder()

	handler := handleDeleteUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

// TestHTTPEndpoints_TableDriven - Comprehensive table-driven test
func TestHTTPEndpoints_TableDriven_Answer(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
		setupStore     func(*SimpleUserStore) // Optional setup function
	}{
		{
			name:           "GET /users empty",
			method:         http.MethodGet,
			path:           "/users",
			body:           "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /users valid",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"1","name":"Alice","age":30}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "POST /users invalid json",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST /users missing fields",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"name":"Alice"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GET /users/:id not found",
			method:         http.MethodGet,
			path:           "/users/999",
			body:           "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "DELETE /users/:id not found",
			method:         http.MethodDelete,
			path:           "/users/999",
			body:           "",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewSimpleUserStore()
			if tt.setupStore != nil {
				tt.setupStore(store)
			}

			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			rr := httptest.NewRecorder()
			router := setupRouter(store)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// Middleware Testing Examples

// TestMiddleware_Logging_Answer demonstrates testing a logging middleware
func TestMiddleware_Logging_Answer(t *testing.T) {
	// Test that logging middleware executes and logs request details
	store := NewSimpleUserStore()

	// Create a custom response writer to capture what's written
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	// Track if middleware was called
	middlewareCalled := false

	// Create a simple logging middleware
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			// In real scenario, this would log to a logger
			// log.Printf("Request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	// Wrap the handler with middleware
	handler := loggingMiddleware(handleListUsers(store))
	handler.ServeHTTP(rr, req)

	// Verify middleware was called
	if !middlewareCalled {
		t.Error("logging middleware was not called")
	}

	// Verify the handler still works correctly
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

// TestMiddleware_Logging_CapturesOutput_Answer shows how to test actual log output
func TestMiddleware_Logging_CapturesOutput_Answer(t *testing.T) {
	store := NewSimpleUserStore()

	// Captured log entries
	var loggedMessages []string

	// Middleware that logs to our slice instead of stdout
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture log message
			logMessage := fmt.Sprintf("Method: %s, Path: %s", r.Method, r.URL.Path)
			loggedMessages = append(loggedMessages, logMessage)
			next.ServeHTTP(w, r)
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler := loggingMiddleware(handleListUsers(store))
	handler.ServeHTTP(rr, req)

	// Verify log was captured
	if len(loggedMessages) != 1 {
		t.Errorf("expected 1 log message, got %d", len(loggedMessages))
	}

	expectedLog := "Method: GET, Path: /users"
	if loggedMessages[0] != expectedLog {
		t.Errorf("expected log '%s', got '%s'", expectedLog, loggedMessages[0])
	}
}

// TestMiddleware_Logging_ChainedMiddleware_Answer shows testing multiple middleware
func TestMiddleware_Logging_ChainedMiddleware_Answer(t *testing.T) {
	store := NewSimpleUserStore()

	var executionOrder []string

	// First middleware
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "logging_before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "logging_after")
		})
	}

	// Second middleware
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "auth_before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "auth_after")
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	// Chain middleware: logging -> auth -> handler
	handler := loggingMiddleware(authMiddleware(handleListUsers(store)))
	handler.ServeHTTP(rr, req)

	// Verify execution order (LIFO - Last In First Out for middleware)
	expectedOrder := []string{
		"logging_before",
		"auth_before",
		"auth_after",
		"logging_after",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("execution order[%d]: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TestMiddleware_Logging_WithStatusCode_Answer shows a simpler middleware logging example
func TestMiddleware_Logging_WithStatusCode_Answer(t *testing.T) {
	store := NewSimpleUserStore()

	var loggedMethod, loggedPath string
	var middlewareCalled bool

	// Middleware that logs request details
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			loggedMethod = r.Method
			loggedPath = r.URL.Path
			// In production, you'd log: log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	rr := httptest.NewRecorder()

	handler := loggingMiddleware(handleGetUser(store))
	handler.ServeHTTP(rr, req)

	// Verify middleware was called
	if !middlewareCalled {
		t.Error("middleware was not called")
	}

	// Verify middleware logged the correct values
	if loggedMethod != http.MethodGet {
		t.Errorf("expected method GET, logged %s", loggedMethod)
	}

	if loggedPath != "/users/999" {
		t.Errorf("expected path /users/999, logged %s", loggedPath)
	}

	// Verify the actual response status
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

// TestMiddleware_Logging_TableDriven_Answer shows testing middleware with multiple requests
func TestMiddleware_Logging_TableDriven_Answer(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		expectedLog string
	}{
		{
			name:        "GET_users",
			method:      http.MethodGet,
			path:        "/users",
			expectedLog: "GET /users",
		},
		{
			name:        "GET_user_by_id",
			method:      http.MethodGet,
			path:        "/users/1",
			expectedLog: "GET /users/1",
		},
		{
			name:        "POST_users",
			method:      http.MethodPost,
			path:        "/users",
			expectedLog: "POST /users",
		},
		{
			name:        "DELETE_user",
			method:      http.MethodDelete,
			path:        "/users/1",
			expectedLog: "DELETE /users/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewSimpleUserStore()
			var loggedMessage string

			loggingMiddleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					loggedMessage = fmt.Sprintf("%s %s", r.Method, r.URL.Path)
					next.ServeHTTP(w, r)
				})
			}

			var body string
			if tt.method == http.MethodPost {
				body = `{"id":"1","name":"Alice","age":30}`
			}

			var req *http.Request
			if body != "" {
				req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			rr := httptest.NewRecorder()

			// Select appropriate handler based on path
			var handler http.Handler
			if tt.path == "/users" {
				if tt.method == http.MethodGet {
					handler = loggingMiddleware(handleListUsers(store))
				} else {
					handler = loggingMiddleware(handleCreateUser(store))
				}
			} else {
				if tt.method == http.MethodGet {
					handler = loggingMiddleware(handleGetUser(store))
				} else {
					handler = loggingMiddleware(handleDeleteUser(store))
				}
			}

			handler.ServeHTTP(rr, req)

			if loggedMessage != tt.expectedLog {
				t.Errorf("expected log '%s', got '%s'", tt.expectedLog, loggedMessage)
			}
		})
	}
}

// Benchmarks

func BenchmarkHandleCreateUser_Answer(b *testing.B) {
	store := NewSimpleUserStore()

	for i := 0; i < b.N; i++ {
		body := fmt.Sprintf(`{"id":"%d","name":"User%d","age":25}`, i, i)
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := handleCreateUser(store)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkHandleListUsers_Answer(b *testing.B) {
	store := NewSimpleUserStore()
	// Add some users first
	for i := 0; i < 10; i++ {
		store.Create(User{ID: fmt.Sprintf("%d", i), Name: fmt.Sprintf("User%d", i), Age: 25})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()

		handler := handleListUsers(store)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkHandleGetUser_Answer(b *testing.B) {
	store := NewSimpleUserStore()
	store.Create(User{ID: "1", Name: "Alice", Age: 30})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rr := httptest.NewRecorder()

		handler := handleGetUser(store)
		handler.ServeHTTP(rr, req)
	}
}
