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
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TODO: Define UserStore struct with RWMutex and map

type UserStore struct {
	users map[string]User
	mu    sync.RWMutex
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

	// ANSWER: If you forget to unlock:
	// 1. DEADLOCK: The mutex stays locked forever
	// 2. Other goroutines trying Lock() or RLock() will block indefinitely
	// 3. Your entire application will freeze/hang
	// WHY defer? It ensures unlock happens even if:
	//   - Function panics
	//   - Multiple return paths exist
	//   - You forget manual unlock
	us.mu.Lock()
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
	// COMPOSITION: Server "has-a" UserStore (not "is-a")
	// Benefits:
	//   - Clean separation of concerns (HTTP logic vs data logic)
	//   - Easy to swap stores (e.g., database vs in-memory)
	//   - Testability: can mock store in tests
	// WHY POINTER?
	//   - UserStore contains a mutex; copying mutexes is forbidden
	//   - All handlers share the SAME store instance (not copies)
	//   - Memory efficient: passing pointer (8 bytes) vs entire struct
	store *UserStore
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", s.handleUsers)
	mux.HandleFunc("/users/", s.handleUserByID)
	mux.HandleFunc("/healthz", s.handleHealthz)

	var h http.Handler = mux

	h = Logging(h)
	// SYNTAX EXPLANATION: RequestTimeout(1 * time.Second)(h)
	// Step 1: RequestTimeout(1 * time.Second) -> returns a Middleware function
	// Step 2: That Middleware function is called with (h) -> returns wrapped Handler
	// It's a function that returns a function (higher-order function)
	// Equivalent to:
	//   middleware := RequestTimeout(1 * time.Second)
	//   h = middleware(h)
	h = RequestTimeout(1 * time.Second)(h)
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
		// WHY THIS BLOCK?
		// 1. s.store.Create() can return errors (duplicate ID, context timeout)
		// 2. We need to handle those errors and send appropriate HTTP response
		// 3. If create fails, client should know (409 Conflict for duplicate)
		// 4. Without this check, errors would be silently ignored
		if err := s.store.Create(r.Context(), u); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		// time.Sleep(100 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// WHAT IT RETURNS?
		// json.NewEncoder(w).Encode(u) returns an error (or nil)
		// We use _ to explicitly ignore it (acceptable for simple demos)
		// In production, you'd log/handle: if err := json.NewEncoder(w).Encode(u); err != nil {...}
		_ = json.NewEncoder(w).Encode(u)
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

// WHAT IS THIS?
// Middleware is a type alias for a function that:
//   - Takes an http.Handler (the "next" handler in chain)
//   - Returns an http.Handler (wrapped version with added behavior)
//
// Pattern: wrap handlers to add cross-cutting concerns (logging, auth, timeouts)
// Example flow: Request -> RateLimiter -> RequestTimeout -> Logging -> Your Handler
type Middleware func(http.Handler) http.Handler

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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HOW DOES IT CALCULATE TOTAL TIME TAKEN BY API?
		// 1. Record start time BEFORE calling next handler
		start := time.Now()
		// 2. Wrap ResponseWriter to capture status code
		ww := &statusWriter{ResponseWriter: w, status: 200}
		// 3. Call next handler (this is where actual work happens - could take seconds)
		next.ServeHTTP(ww, r)
		// 4. Calculate duration = current time - start time
		// time.Since(start) internally does time.Now().Sub(start)
		dur := time.Since(start)
		// 5. Log everything (method, path, status captured by wrapper, duration)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, ww.status, dur)
	})
}

// TODO: Implement RequestTimeout middleware
func RequestTimeout(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TODO: Implement RateLimiter middleware
// HOW THIS WORKS (Token Bucket Algorithm):
// 1. Create a buffered channel with 'burst' capacity (acts as token bucket)
// 2. Fill bucket initially with 'burst' tokens (allows burst traffic)
// 3. Start goroutine that refills bucket at 'rate' tokens/second
// 4. For each request: try to take a token from bucket
//   - If token available: proceed
//   - If bucket empty: request is rate-limited (blocks until token available or context canceled)

// LINE-BY-LINE EXPLANATION:

func RateLimiter(rate int, burst int) Middleware {
	// FUNCTION SIGNATURE:
	// - Takes: rate (tokens/second), burst (max simultaneous requests)
	// - Returns: Middleware (a function that wraps handlers)
	// - Example: RateLimiter(20, 10) means 20 req/sec with burst of 10

	// VALIDATION: Ensure positive values
	if rate <= 0 {
		rate = 1 // Default: 1 token per second
	}
	if burst <= 0 {
		burst = 1 // Default: burst of 1
	}

	// CREATE TOKEN BUCKET:
	// tokens := make(chan struct{}, burst)
	// - Buffered channel with capacity 'burst'
	// - Acts as the "bucket" holding available tokens
	// - struct{} uses 0 bytes (memory efficient)
	// - When full: burst requests can proceed immediately
	// - When empty: requests must wait for refill
	tokens := make(chan struct{}, burst)

	// INITIAL FILL: Pre-populate bucket with 'burst' tokens
	// This allows burst traffic at startup without waiting
	// WHAT THIS DOES: struct{}{} is an empty struct (0 bytes)
	// WHY THIS SYNTAX?
	//   - chan struct{} is idiomatic for signaling (we don't need data, just a signal)
	//   - struct{}{} creates an instance of empty struct
	//   - More memory-efficient than chan bool or chan int
	//   - Common Go pattern for "fire and forget" signals
	for i := 0; i < burst; i++ {
		tokens <- struct{}{} // Send 'burst' tokens into channel
	}

	// START REFILL GOROUTINE:
	// This goroutine runs in background for the lifetime of the server
	go func() {
		// CREATE TICKER:
		// t := time.NewTicker(time.Second / time.Duration(rate))
		// - Fires every (1 second / rate) interval
		// - Example: rate=20 -> ticker fires every 50ms (1000ms/20)
		// - Example: rate=5 -> ticker fires every 200ms (1000ms/5)
		t := time.NewTicker(time.Second / time.Duration(rate))
		defer t.Stop() // Clean up ticker when function exits

		// REFILL LOOP:
		// for range t.C {...}
		// - Executes every time ticker fires
		// - t.C is a channel that receives time.Time values on each tick
		// - We ignore the time value (don't need it)
		for range t.C {
			// TRY TO ADD TOKEN:
			select {
			case tokens <- struct{}{}:
				// SUCCESS: Token added to bucket
				// This happens when bucket is not full (< burst capacity)
			default:
				// BUCKET FULL: Do nothing, drop this token
				// This prevents bucket from exceeding burst capacity
				// Similar to water overflowing from a full bucket
			}
		}
	}()

	// RETURN MIDDLEWARE FUNCTION:
	// This function wraps the next handler with rate limiting logic
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// FOR EACH REQUEST:
			// Try to consume a token before proceeding
			select {
			// CASE 1: Request canceled/timed out
			case <-r.Context().Done():
				// Context canceled (client disconnected, timeout, etc.)
				http.Error(w, r.Context().Err().Error(), http.StatusGatewayTimeout)
				return

			// CASE 2: Token available (proceed with request)
			case <-tokens:
				// TOKEN CONSUMED: Remove one token from bucket
				// - If bucket has tokens: proceeds immediately
				// - If bucket empty: BLOCKS here until refill goroutine adds token
				// - Blocking means request waits in queue (rate limiting in action)
				next.ServeHTTP(w, r) // Forward to next handler
			}
		})
	}
}

// FLOW EXAMPLE (rate=20, burst=10):
// 1. Server starts: bucket has 10 tokens
// 2. 10 requests arrive instantly: all proceed (burst capacity)
// 3. 11th request: bucket empty, blocks waiting for refill
// 4. After 50ms (1000ms/20): refill goroutine adds 1 token
// 5. 11th request proceeds, consumes that token
// 6. Pattern continues: 20 requests/second sustained rate
//
// WHY THIS WORKS:
// - Allows bursts without penalizing occasional spikes
// - Prevents sustained overload (enforces long-term rate)
// - Fair: requests queue in order (FIFO via channel)

// 4. Graceful Shutdown
// TODO: Set up signal handling and call Server.Shutdown with context

// 5. Client & Tests (Extra Challenge)
// TODO: Write a client function that calls your server with context timeout
// TODO: Write table-driven tests for UserStore methods

func main() {
	// TODO: Initialize UserStore
	us := NewUserStore()

	// WHY context.Background()?
	// - Used at the TOP LEVEL of your application (main function)
	// - No parent context exists to derive from
	// - Creates a root context that's never canceled
	// - In HTTP handlers, you'd use r.Context() instead (request-scoped)
	// - context.Background() vs r.Context():
	//     Background: for testing, initialization, top-level operations
	//     r.Context(): for HTTP handlers (auto-canceled on client disconnect/timeout)
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

	// EXPLAIN THIS LINE BY EACH WORD:
	// srv :=
	//   Create and assign to variable 'srv'
	// &http.Server{...}
	//   Create a pointer to http.Server struct (& means "address of")
	// Addr: ":8080"
	//   Server listens on all interfaces (0.0.0.0) on port 8080
	//   Empty string before : means "all network interfaces"
	// Handler: (&Server{store: us}).routes()
	//   Step 1: &Server{store: us} -> create Server instance with our UserStore
	//   Step 2: .routes() -> call routes() method which returns middleware-wrapped handler
	//   This handler processes ALL incoming HTTP requests
	srv := &http.Server{
		Addr:    ":8080",
		Handler: (&Server{store: us}).routes(),
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
	// WHAT defer cancel() DOES:
	// 1. context.WithTimeout creates a context that auto-cancels after 5 seconds
	// 2. It returns (ctx, cancel) where cancel is a function to release resources
	// 3. defer cancel() ensures cleanup happens when main() exits
	// 4. Why needed? Prevents resource leaks (timers, goroutines)
	// 5. Always call cancel() for contexts you create with With* functions
	// 6. defer ensures it runs even if srv.Shutdown panics
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
