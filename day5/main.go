package main

// --- Day 5 Coding Template ---
// Follow the TODOs and function signatures below to implement each concept step by step.
// For each section, write your code in the marked region and test as you go.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"

	"net/http"
	"os/signal"
	"sync"
	"time"
)

// 1. User model and store
// TODO: Define User struct (ID, Name, Age)

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}

// TODO: Define UserStore struct with RWMutex and map

type UserStore struct {
	users map[string]User
	mu sync.RWMutex
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]User),
	}
}

// TODO: Implement Create, Get, List, Delete methods (context-aware)
func (us *UserStore) Create(ctx context.Context, user User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	us.mu.Lock() // what will happen if we forget to unlock? and what will happen with other readers/writers?
	defer us.mu.Unlock()
	
	if _, exists := us.users[user.ID]; exists {
		return fmt.Errorf("User ID already Exists!")
	}
	us.users[user.ID] = user
	fmt.Println("User Created:", user.ID)
	return nil
}

func (us *UserStore) Get(ctx context.Context, id string) (User, error) {

	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}

	us.mu.RLock()
	defer us.mu.RUnlock()	
	user, exists := us.users[id]
	if !exists {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (us *UserStore) List(ctx context.Context) ([]User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	us.mu.RLock()
	defer us.mu.RUnlock()	
	users := []User{}
	for _, user := range us.users {
		users = append(users, user)
	}
	return users, nil
}

func (us *UserStore) Delete(ctx context.Context, id string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	us.mu.Lock()
	defer us.mu.Unlock()	

	if _, exists := us.users[id]; !exists {
		return fmt.Errorf("User with id %s not found", id)
	}

	delete(us.users, id)
	fmt.Printf("User Deleted with id %s\n", id)
	return nil
}	

// 2. HTTP Handlers

type Server struct {
	store *UserStore // Why use composition here? why pointer?
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", s.handleUsers)
	mux.HandleFunc("/users/", s.handleUserByID)
	mux.HandleFunc("/healthz", s.handleHealthz)

	var h http.Handler = mux

	h = Logging(h)
	h = RequestTimeout(1 * time.Second)(h) // what is this syntax? what does it do?
	h = RateLimiter(20, 10)(h)
	return h
}


// TODO: Implement handler for POST /users (create user)
// TODO: Implement handler for GET /users (list users)
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if u.ID == "" || u.Name == "" || u.Age <= 0 {
			http.Error(w, "id, name, age required", http.StatusBadRequest)
			return
		}
		if err := s.store.Create(r.Context(), u); err != nil { // Why this block?
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		// time.Sleep(100 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(u) // what it returns?
	case http.MethodGet:
		users, err := s.store.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusRequestTimeout)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(users)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// TODO: Implement handler for GET /users/{id} (get user)
// TODO: Implement handler for DELETE /users/{id} (delete user)
func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	// print r in console with all its formatting
	// fmt.Printf("%+v\n", r.Context())
	if id == "" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet: 
		users, err := s.store.Get(r.Context(), id)
		if err != nil {
			if err == context.DeadlineExceeded {
				http.Error(w, err.Error(), http.StatusRequestTimeout)
			} else {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
			return
		}
		_ = json.NewEncoder(w).Encode(users)
	case http.MethodDelete:
		if err := s.store.Delete(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusRequestTimeout)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// TODO: Implement handler for GET /healthz (simulate dependency check with context)
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(100 * time.Millisecond):
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	case <-r.Context().Done():
		http.Error(w, "health check canceled", http.StatusServiceUnavailable)
	}
}

// 3. Middleware

type Middleware func (http.Handler) http.Handler // what is this?

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// TODO: Implement Logging middleware
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		dur := time.Since(start)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, ww.status, dur) // how does it calculate total time taken by API?
	})
}
// TODO: Implement RequestTimeout middleware
func RequestTimeout(timeout time.Duration) Middleware {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TODO: Implement RateLimiter middleware
func RateLimiter(rate int, burst int) Middleware { // Explain how this works?
	if rate <= 0 {
		rate = 1
	}
	if burst <= 0 {
		burst = 1
	}
	tokens := make(chan struct{}, burst)
	for i := 0; i < burst; i++ {
		tokens <- struct{}{} // explain what this does? and why this syntax?
	}
	go func() {
		t := time.NewTicker(time.Second / time.Duration(rate))
		defer t.Stop()
		for range t.C {
			select {
			case tokens <- struct{}{}:
			default:
			}
		}
	}()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			select {
			case <-r.Context().Done():
				http.Error(w, r.Context().Err().Error(), http.StatusGatewayTimeout)
				return 
			case <-tokens:
				next.ServeHTTP(w, r)
			}
		})
	}
}

// 4. Graceful Shutdown
// TODO: Set up signal handling and call Server.Shutdown with context

// 5. Client & Tests (Extra Challenge)
// TODO: Write a client function that calls your server with context timeout
// TODO: Write table-driven tests for UserStore methods

func main() {
	// TODO: Initialize UserStore
	us := NewUserStore()
	// user := User{
	// 	ID: "1", 
	// 	Name: "Alice", 
	// 	Age: 30,
	// }
	// us.Create(context.Background(), user) // why context.Background()?
	// user = User{
	// 	ID: "2", 
	// 	Name: "Bob", 
	// 	Age: 25,
	// }
	// us.Create(context.Background(), user) // why context.Background()?
	// getUser, err := us.Get(context.Background(), "1")
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }else {
	// 	fmt.Println("User fetched:", getUser)
	// }

	// userList, err := us.List(context.Background())
	// if err != nil {
	// 	fmt.Errorf("Error listing users: %v", err)
	// } else {
	// 	fmt.Println("User list:", userList)
	// }
	// err = us.Delete(context.Background(), "1")
	// if err != nil {
	// 	fmt.Printf("Error deleting user: %v\n", err)
	// }else {
	// 	fmt.Println("User deleted")
	// }
	// userList, err = us.List(context.Background())
	// if err != nil {
	// 	fmt.Errorf("Error listing users: %v", err)
	// } else {
	// 	fmt.Println("User list:", userList)
	// }

	// TODO: Set up HTTP server and routes

	srv := &http.Server{
		Addr:    ":8080",
		Handler: (&Server{store: us}).routes(), // Explain this lines by each word what it does ?
	}

	// TODO: Wrap handlers with middleware

	// TODO: Start server in goroutine
	// Start server
	go func() {
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// TODO: Wait for signal and gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Explain what it does?

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
