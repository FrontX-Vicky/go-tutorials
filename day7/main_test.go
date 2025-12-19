package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Day 7 Test Template: HTTP Handler & Integration Testing

// 1. Test GET /users Endpoint
func TestHandleListUsers_EmptyStore(t *testing.T) {
	// TODO: Create store, call GET /users, check status, check empty JSON array, check Content-Type
	store := NewSimpleUserStore()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler := handleListUsers(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var users []User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func TestHandleListUsers_WithUsers(t *testing.T) {
	// TODO: Add users to store, call GET /users, check status, check user list, check Content-Type
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
		t.Errorf("Expected 2 users, got %d", len(users))
	}

}

// 2. Test POST /users Endpoint
func TestHandleCreateUser_Valid(t *testing.T) {
	// TODO: Create valid user, call POST /users, check status 201, check response JSON
	store := NewSimpleUserStore()
	userJSON := `{"id":"1","name":"Alice","age":30}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handleCreateUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	user, err := store.Get("1")
	if err != nil {
		t.Errorf("user not found in store: %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", user.Name)
	}

}

func TestHandleCreateUser_InvalidJSON(t *testing.T) {
	// TODO: Send malformed JSON, check status 400
	store := NewSimpleUserStore()
	invalidJSON := `{"id":"1","name":"Alice","age":30` // Missing closing brace

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handleCreateUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHandleCreateUser_DuplicateID(t *testing.T) {
	// TODO: Create user, then try to create with same ID, check status 409
	store := NewSimpleUserStore()
	// store.Create(User{ID: "1", Name: "Alice", Age: 30})
	body := `{"id":"1","name":"Alice","age":30}`

	req1 := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()

	handler := handleCreateUser(store)
	handler.ServeHTTP(rr1, req1)

	if rr1.Code != http.StatusCreated {
		t.Errorf("first create failed: expected 201, got %d", rr1.Code)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()

	handler.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusConflict {
		t.Errorf("expected status 409 on duplicate, got %d", rr2.Code)
	}
}

func TestHandleCreateUser_MissingFields(t *testing.T) {
	// TODO: Table-driven: missing id, missing name, empty id, check status 400
	store := NewSimpleUserStore()
	tests := []struct {
		name string
		body string
	}{
		{"Missing Id", `{"name":"Alice", "age":30}`},
		{"Missing Name", `{"id":"1", "age":30}`},
		{"Empty Id", `{"id":"", "name":"alice", "age":30}`},
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

// 3. Test GET /users/:id Endpoint
func TestHandleGetUser_Existing(t *testing.T) {
	// TEST SETUP PHASE - Preparing the test environment
	// ================================================

	// STEP 1: Create an empty in-memory user store
	// Data flow: NewSimpleUserStore() creates a new map[string]User inside the store
	store := NewSimpleUserStore()

	// STEP 2: Create a User struct with test data
	// Data flow: User struct is created in memory with ID="1", Name="Alice", Age=30
	user := User{ID: "1", Name: "Alice", Age: 30}

	// STEP 3: Add the user to the store (simulates database storage)
	// Data flow: store.Create(user) -> adds user to store.users["1"] = User{...}
	// Now the store has ONE user stored in memory that can be retrieved
	store.Create(user)

	// REQUEST CREATION PHASE - Simulating an HTTP GET request
	// ========================================================

	// STEP 4: Create a fake HTTP GET request to "/users/1"
	// Data flow: httptest.NewRequest creates a mock http.Request object
	// - Method: GET
	// - URL: "/users/1" (the "1" is user.ID which will be extracted by the handler)
	// - Body: nil (GET requests don't have a body)
	req := httptest.NewRequest(http.MethodGet, "/users/"+user.ID, nil)

	// STEP 5: Create a response recorder to capture what the handler writes
	// Data flow: httptest.NewRecorder creates a mock http.ResponseWriter
	// This will record: status code, headers, and response body
	rr := httptest.NewRecorder()

	// STEP 6: Get the HTTP handler function that will process the request
	// Data flow: handleGetUser(store) returns an http.HandlerFunc
	// This handler has access to the store (closure) and will:
	// - Extract user ID from the URL path
	// - Look up the user in the store
	// - Return the user as JSON
	handler := handleGetUser(store)

	// REQUEST SERVING PHASE - The actual request processing
	// ======================================================

	// STEP 7: Execute the handler (this is where the magic happens!)
	// Data flow:
	// 1. handler.ServeHTTP(rr, req) is called
	// 2. Handler extracts "1" from URL "/users/1"
	// 3. Handler calls store.Get("1")
	// 4. Store looks up users["1"] and returns User{ID:"1", Name:"Alice", Age:30}
	// 5. Handler encodes User to JSON: {"id":"1","name":"Alice","age":30}
	// 6. Handler writes: status 200, Content-Type header, JSON body to rr (response recorder)
	handler.ServeHTTP(rr, req)

	// RESPONSE VERIFICATION PHASE - Checking the handler's response
	// =============================================================

	// STEP 8: Verify the HTTP status code is 200 OK
	// Data flow: rr.Code contains the status code written by the handler
	// Expected: 200 (user was found and returned successfully)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// STEP 9: Decode the JSON response body into a User struct
	// Data flow:
	// 1. rr.Body contains the response bytes: {"id":"1","name":"Alice","age":30}
	// 2. json.NewDecoder(rr.Body) creates a JSON decoder
	// 3. Decode(&got) parses the JSON and fills the 'got' variable with User data
	// 4. Now 'got' should contain User{ID:"1", Name:"Alice", Age:30}
	var got User
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	// At this point, the test verifies:
	// ✓ Handler returned status 200
	// ✓ Response body is valid JSON
	// ✓ JSON can be decoded into a User struct
	// (Missing: verification that got.ID == "1", got.Name == "Alice", etc.)
}

func TestHandleGetUser_NotFound(t *testing.T) {
	// TODO: Call GET /users/:id for non-existent user, check status 404
	store := NewSimpleUserStore()

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rr := httptest.NewRecorder()

	handler := handleGetUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandleGetUser_InvalidID(t *testing.T) {
	// TODO: Call GET /users/ (empty id), check status 400
	store := NewSimpleUserStore()

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	rr := httptest.NewRecorder()

	handler := handleGetUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

// 4. Test DELETE /users/:id Endpoint
func TestHandleDeleteUser_Existing(t *testing.T) {
	// TODO: Add user, call DELETE /users/:id, check status 200, verify user deleted
	store := NewSimpleUserStore()
	user := User{ID: "1", Name: "Alice", Age: 30}
	store.Create(user)
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)

	rr := httptest.NewRecorder()

	handler := handleDeleteUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleDeleteUser_NotFound(t *testing.T) {
	// TODO: Call DELETE /users/:id for non-existent user, check status 404
	store := NewSimpleUserStore()
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rr := httptest.NewRecorder()

	handler := handleDeleteUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected %d got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandleDeleteUser_InvalidID(t *testing.T) {
	// TODO: Call DELETE /users/ (empty id), check status 400
	store := NewSimpleUserStore()
	req := httptest.NewRequest(http.MethodDelete, "/users/", nil)
	rr := httptest.NewRecorder()

	handler := handleDeleteUser(store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, rr.Code)
	}
}

// 5. Table-Driven HTTP Tests Example
func TestHTTPEndpoints_TableDriven(t *testing.T) {
	// TODO: Use table-driven pattern to test multiple endpoints and methods
	// See README for example structure

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
		setupStore     func(*SimpleUserStore)
	}{
		{
			name:           "GET_/users_empty",
			method:         http.MethodGet,
			path:           "/users",
			body:           "",
			expectedStatus: http.StatusOK, // 200
		},
		{
			name:           "POST_/users_valid",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"1", "name":"Alice", "age":30}`,
			expectedStatus: http.StatusCreated, // 201
		},
		{
			name:           "GET_/users/1_not_found",
			method:         http.MethodGet,
			path:           "/users/1",
			body:           "",
			expectedStatus: http.StatusNotFound, // 404
		},
		{
			name:           "DELETE_/users/1_invalid_id",
			method:         http.MethodDelete,
			path:           "/users/1",
			body:           "",
			expectedStatus: http.StatusNotFound, // 404
		},
		{
			name:           "DELETE_/users/_invalid_id",
			method:         http.MethodDelete,
			path:           "/users/",
			body:           "",
			expectedStatus: http.StatusBadRequest, // 400
		},
		{
			name:           "POST_/users_duplicate_id",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"1", "name":"Alice", "age":30}`,
			expectedStatus: http.StatusConflict, // 409
			setupStore: func(store *SimpleUserStore) {
				store.Create(User{ID: "1", Name: "Alice", Age: 30})
			},
		},
		{
			name:           "GET_/users/1_existing",
			method:         http.MethodGet,
			path:           "/users/1",
			body:           "",
			expectedStatus: http.StatusOK, // 200
			setupStore: func(store *SimpleUserStore) {
				store.Create(User{ID: "1", Name: "Alice", Age: 30})
			},
		},
		{
			name:           "DELETE_/users/1_existing",
			method:         http.MethodDelete,
			path:           "/users/1",
			body:           "",
			expectedStatus: http.StatusOK, // 200
			setupStore: func(store *SimpleUserStore) {
				store.Create(User{ID: "1", Name: "Alice", Age: 30})
			},
		},
		{
			name:           "DELETE_/users/_invalid_id",
			method:         http.MethodDelete,
			path:           "/users/",
			body:           "",
			expectedStatus: http.StatusBadRequest, // 400
		},
		{
			name:           "POST_/users_invalid_JSON",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"1", "name":"Alice", "age":30`, // Invalid JSON (missing closing brace)
			expectedStatus: http.StatusBadRequest,                 // 400
		},
		{
			name:           "POST_/users_missing_fields",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"2", "name":"", "age":0}`, // Missing name and age fields
			expectedStatus: http.StatusBadRequest,            // 400
		},
		{
			name:           "POST_/users_empty_id",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"", "name":"Alice", "age":30}`, // Empty id field
			expectedStatus: http.StatusBadRequest,                 // 400
		},
		{
			name:           "GET_/users/_non-existent",
			method:         http.MethodGet,
			path:           "/users/999",
			body:           "",
			expectedStatus: http.StatusNotFound, // 404
		},
		{
			name:           "DELETE_/users/_non-existent",
			method:         http.MethodDelete,
			path:           "/users/999",
			body:           "",
			expectedStatus: http.StatusNotFound, // 404
		},
		{
			name:           "GET_/users/_invalid_id",
			method:         http.MethodGet,
			path:           "/users/abc",
			body:           "",
			expectedStatus: http.StatusBadRequest, // 400
		},
		{
			name:           "DELETE_/users/_invalid_id_format",
			method:         http.MethodDelete,
			path:           "/users/abc",
			body:           "",
			expectedStatus: http.StatusBadRequest, // 400
		},
		{
			name:           "POST_/users_negative_age",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"3", "name":"Bob", "age":-5}`, // Negative age
			expectedStatus: http.StatusBadRequest,                // 400
		},
		{
			name:           "POST_/users_large_age",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"4", "name":"Charlie", "age":150}`, // Large age
			expectedStatus: http.StatusBadRequest,                     // 400
		},
		// {
		// 	name:           "GET_/users/_special_characters_in_id",
		// 	method:         http.MethodGet,
		// 	path:           "/users/!@#$%",
		// 	body:  s         "",
		// 	expectedStatus: http.StatusBadRequest, // 400
		// },
		// {
		// 	name:           "DELETE_/users/_special_characters_in_id",
		// 	method:         http.MethodDelete,
		// 	path:           "/users/!@#$%",
		// 	body:           "",
		// 	expectedStatus: http.StatusBadRequest, // 400
		// },
		{
			name:           "POST_/users_extremely_large_payload",
			method:         http.MethodPost,
			path:           "/users",
			body:           `{"id":"5", "name":"Dave", "age":40, "extra":"` + strings.Repeat("x", 10000) + `"}`, // Extremely large payload
			expectedStatus: http.StatusBadRequest,                                                               // 400
		},
		// {
		// 	name:           "GET_/users/_SQL_injection_attempt",
		// 	method:         http.MethodGet,
		// 	path:           "/users/1; DROP TABLE users",
		// 	body:           "",
		// 	expectedStatus: http.StatusBadRequest, // 400
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewSimpleUserStore()

			if tt.setupStore != nil {
				tt.setupStore(store)
			}

			var req *http.Request
			if tt.method == http.MethodPost {
				// escape path for invalid chars like %
				if strings.ContainsAny(tt.path, "%!@#$^&*()") {
					tt.path = url.PathEscape(tt.path)
				}
				// prevent SQL injection-like paths
				// if strings.ContainsAny(tt.path, ";") {
				// 	tt.path = url.PathEscape(tt.path)
				// }
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

// 6. Middleware Integration (if implemented)
func TestMiddleware_Logging(t *testing.T) {
	// TODO: If logging middleware exists, test that it executes
	store := NewSimpleUserStore()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	middlewareCalled := false

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	handler := loggingMiddleware(handleListUsers(store))
	handler.ServeHTTP(rr, req)

	if !middlewareCalled {
		t.Errorf("logging middleware was not called")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestMiddleware_Timeout(t *testing.T) {
	// Test timeout middleware that cancels long-running requests

	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a long-running process (200ms)
		select {
		case <-time.After(200 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"completed"}`))
		case <-r.Context().Done():
			// Context was cancelled (timeout)
			w.WriteHeader(http.StatusRequestTimeout) // 408
			w.Write([]byte(`{"error":"request timed out"}`))
			return
		}
	})

	timeoutMiddleware := func(timeout time.Duration) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx, cancel := context.WithTimeout(r.Context(), timeout)
				defer cancel()

				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}
	}

	t.Run("TimeoutExceeded", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		rr := httptest.NewRecorder()
		// Handler takes 200ms, but timeout is 50ms - should timeout
		handler := timeoutMiddleware(50 * time.Millisecond)(slowHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusRequestTimeout {
			t.Errorf("expected status 408, got %d", rr.Code)
		}

		if !strings.Contains(rr.Body.String(), "request timed out") {
			t.Errorf("expected timeout error message, got %s", rr.Body.String())
		}
	})

	t.Run("CompletedBeforeTimeout", func(t *testing.T) {
		fastHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Fast operation (10ms)
			time.Sleep(10 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"completed"}`))
		})

		req := httptest.NewRequest(http.MethodGet, "/fast", nil)
		rr := httptest.NewRecorder()
		// Handler takes 10ms, timeout is 100ms - should complete
		handler := timeoutMiddleware(100 * time.Millisecond)(fastHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}

		if !strings.Contains(rr.Body.String(), "completed") {
			t.Errorf("expected success message, got %s", rr.Body.String())
		}
	})
}

func TestMiddleware_RateLimiter(t *testing.T) {
	// TODO: If rate limiter exists, test that it blocks after threshold
	store := NewSimpleUserStore()

	rateLimiterMiddleware := func(maxRequests int) func(http.Handler) http.Handler {
		var requestCount int 
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
				requestCount++
				if requestCount > maxRequests {
					http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
			})
		}
	}

	t.Run("blocks after threshold", func(t *testing.T) {
		handler := rateLimiterMiddleware(3)(handleListUsers(store))

		for i := 1; i <=3; i++ {
			request := httptest.NewRequest(http.MethodGet, "/users", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, request)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d on request %d", rr.Code, i)
			}
		} 

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusTooManyRequests {
			t.Errorf("expected status 429, got %d after exceeding limit", rr.Code)
		}

		if !strings.Contains(rr.Body.String(),  `{"error":"rate limit exceeded"}`) {
			t.Errorf("expected error message 'Rate limit exceeded', got %s", rr.Body.String())
		}
	})

}

// 7. Error Response Format
func TestErrorResponses_Format(t *testing.T) {
	// TODO: Test error message format for 400, 404, 500

	store := NewSimpleUserStore()

	tests := []struct {
		name           string
		setupRequest   func() (*http.Request, http.Handler)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "400_Bad_request_invalid_json",
			setupRequest: func() (*http.Request, http.Handler) {
				req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"invalid json`))
				req.Header.Set("Content-Type", "application/json")
				return req, handleCreateUser(store)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError: "invalid",
		},
			{
			name: "404_Not_Found",
			setupRequest: func() (*http.Request, http.Handler) {
				req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
				return req, handleGetUser(store)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "not found",
		},
		{
			name: "400_Bad_Request_Empty_ID",
			setupRequest: func() (*http.Request, http.Handler) {
				req := httptest.NewRequest(http.MethodGet, "/users/", nil)
				return req, handleGetUser(store)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid",
		},
		{
			name: "409_Conflict_Duplicate",
			setupRequest: func() (*http.Request, http.Handler) {
				// Pre-create a user
				store.Create(User{ID: "100", Name: "Existing", Age: 30})
				req := httptest.NewRequest(http.MethodPost, "/users",
					strings.NewReader(`{"id":"100","name":"Duplicate","age":25}`))
				req.Header.Set("Content-Type", "application/json")
				return req, handleCreateUser(store)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "already exists",
		},

	}

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// store := NewSimpleUserStore()
			req, handler := tt.setupRequest()
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			body := strings.ToLower(rr.Body.String())
			if tt.expectedStatus >= 400 && !strings.Contains(body, tt.expectedError) {
				t.Logf("Response body should contain '%s', got: %s", tt.expectedError, rr.Body.String())
			}

			contentType := rr.Header().Get("Content-Type")
			if tt.expectedStatus >= 400 && contentType != "" {
				if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "text/plain") {
					t.Logf("expected Content-Type application/json, got %s", contentType)
				}
			}


		})
	}
}

// 8. Benchmark Handlers
func BenchmarkHandleCreateUser(b *testing.B) {
	// TODO: Benchmark POST /users handler
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

func BenchmarkHandleListUsers(b *testing.B) {
	// TODO: Benchmark GET /users handler
	store := NewSimpleUserStore()

	for i := 0; i < 1000; i++ {
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

func BenchmarkHandleGetUser(b *testing.B) {
	// TODO: Benchmark GET /users/:id handler
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