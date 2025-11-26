# Day 5 — Context, HTTP Server, Graceful Shutdown, Middleware

## Concepts to Learn

### 1) Context for Cancellation and Timeouts
- `context.Context` carries deadlines, cancellation signals, and request-scoped values.
- Create contexts: `context.WithCancel`, `context.WithTimeout`, `context.WithDeadline`.
- Always propagate context to I/O calls and downstream functions. Check `<-ctx.Done()`.

### 2) HTTP Server Basics
- Handlers implement `http.Handler` or `func(http.ResponseWriter, *http.Request)`.
- JSON encoding/decoding using `encoding/json`.
- Validation and error responses (status codes, body).

### 3) Graceful Shutdown
- Capture OS signals and call `Server.Shutdown(ctx)`.
- `Shutdown` stops accepting new connections and waits for in-flight requests (until context deadline).

### 4) Middleware Pattern
- Wrap a handler with another function to add cross-cutting concerns (logging, metrics, rate limiting).
- Signature: `func(next http.Handler) http.Handler`.

### 5) Sync Primitives for Safety
- Protect shared state with `sync.Mutex` / `sync.RWMutex`.
- Prefer composition: a store type with methods that lock internally.

## Tasks

1. Implement an in-memory `UserStore`:
   - Data model: `User{ID string, Name string, Age int}`.
   - Methods (context-aware): `Create(ctx, u)`, `Get(ctx, id)`, `List(ctx)`, `Delete(ctx, id)`.
   - Use a `sync.RWMutex` to protect the map.

2. Build HTTP endpoints:
   - `POST /users` → create user (JSON body)
   - `GET /users/{id}` → get by id
   - `GET /users` → list all
   - `DELETE /users/{id}` → delete
   - All handlers should be context-aware and return proper status codes.

3. Add graceful shutdown:
   - Listen for `os.Interrupt` / `syscall.SIGTERM`.
   - On signal: create a timeout context (e.g., 5s), call `srv.Shutdown(ctx)`.

4. Add middleware:
   - `Logging` (method, path, status, duration)
   - `RequestTimeout` (wrap handlers with `context.WithTimeout` per request)

5. Add simple rate limiting (per-process is fine):
   - Use a token bucket with `time.Ticker` and a buffered channel.
   - Middleware should `select` on token channel vs. timeout.

## Extra Challenge
- Add a `GET /healthz` that checks an internal dependency via context (simulate with `time.Sleep` and `select` on `ctx.Done()`).
- Add a client function that calls your server with its own context timeout and prints the result.
- Write 1-2 table-driven tests for the store methods using the standard library `testing` package.
