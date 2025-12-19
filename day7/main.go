package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	// "net/url"
	"strings"
	"sync"
)

// User represents a user entity
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// SimpleUserStore - in-memory storage for testing
type SimpleUserStore struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewSimpleUserStore() *SimpleUserStore {
	return &SimpleUserStore{
		users: make(map[string]User),
	}
}

func (s *SimpleUserStore) Create(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user already exists")
	}
	s.users[user.ID] = user
	return nil
}

func (s *SimpleUserStore) Get(id string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[id]
	if !exists {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *SimpleUserStore) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := []User{}
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *SimpleUserStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(s.users, id)
	return nil
}

// HTTP Handlers - TODO: Implement these to match the requirements

// handleListUsers - GET /users
// Should return:
// - 200 OK with JSON array of users (even if empty)
// - Content-Type: application/json
func handleListUsers(store *SimpleUserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		// 1. Check request method is GET
		// 2. Get list of users from store
		// 3. Set Content-Type header to application/json
		// 4. Write 200 status
		// 5. Marshal and write users as JSON
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		users := store.List()
		json.NewEncoder(w).Encode(users)
	}
}

// handleCreateUser - POST /users
// Should return:
// - 201 Created on success
// - 400 Bad Request if JSON invalid or required fields missing
// - 409 Conflict if user ID already exists
// - Content-Type: application/json
func handleCreateUser(store *SimpleUserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		// 1. Check request method is POST
		// 2. Read and parse JSON body
		// 3. Validate required fields (ID, Name, Age)
		// 4. Call store.Create()
		// 5. Handle errors: invalid JSON -> 400, duplicate -> 409
		// 6. Write 201 status on success
		// 7. Marshal and write created user as JSON
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // 405
			return
		}

		// Limit request body size to 5KB (enough for reasonable user data)
		const maxBodySize = 5 * 1024 // 5KB
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		var user User
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "request body too large or failed to read", http.StatusBadRequest) // 400
			return
		}
		defer r.Body.Close()

		if len(body) == 0 {
			http.Error(w, "empty request body", http.StatusBadRequest) // 400
			return
		}

		if err := json.Unmarshal(body, &user); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest) // 400
			return
		}

		if user.ID == "" || user.Name == "" || user.Age <= 0 || user.Age > 120 {
			http.Error(w, "missing required fields", http.StatusBadRequest) // 400
			return
		}

		if err := store.Create(user); err != nil {
			if strings.Contains(err.Error(), "already exists") {
				http.Error(w, err.Error(), http.StatusConflict) // 409
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(user)
	}
}

// handleGetUser - GET /users/:id
// Should return:
// - 200 OK with user JSON if found
// - 404 Not Found if user doesn't exist
// - Content-Type: application/json
func handleGetUser(store *SimpleUserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		// 1. Check request method is GET
		// 2. Extract ID from URL path
		// 3. Call store.Get(id)
		// 4. Handle error: not found -> 404
		// 5. Write 200 status
		// 6. Marshal and write user as JSON
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/users/")
		// escape id for invalid chars like %
		if id == "" || !isNumeric(id) { // give error for non-numeric IDs
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		user, err := store.Get(id)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// handleDeleteUser - DELETE /users/:id
// Should return:
// - 200 OK on success
// - 404 Not Found if user doesn't exist
// - Content-Type: application/json
func handleDeleteUser(store *SimpleUserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		// 1. Check request method is DELETE
		// 2. Extract ID from URL path
		// 3. Call store.Delete(id)
		// 4. Handle error: not found -> 404
		// 5. Write 200 status
		// 6. Write JSON response confirming deletion
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" || !isNumeric(id) { // give error for non-numeric IDs
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := store.Delete(id); err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "user deleted"})
	}
}

// setupRouter creates an HTTP mux with all routes
func setupRouter(store *SimpleUserStore) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleListUsers(store)(w, r)
		} else if r.Method == http.MethodPost {
			handleCreateUser(store)(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			handleGetUser(store)(w, r)
		} else if r.Method == http.MethodDelete {
			handleDeleteUser(store)(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("Day 7: HTTP Handler Testing - Run 'go test -v' to execute tests")
}
