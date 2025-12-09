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
