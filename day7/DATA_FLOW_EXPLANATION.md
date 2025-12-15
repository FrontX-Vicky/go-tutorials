# Complete Data Flow Explanation for TestHandleGetUser_Existing

## Visual Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    TEST EXECUTION FLOW                          │
└─────────────────────────────────────────────────────────────────┘

1. TEST SETUP (Lines 166-174)
   ┌──────────────────────────────────────────────────────┐
   │ store := NewSimpleUserStore()                        │
   │ Creates: map[string]User{}                           │
   │          (empty in-memory storage)                   │
   └──────────────────────────────────────────────────────┘
                           ↓
   ┌──────────────────────────────────────────────────────┐
   │ user := User{ID: "1", Name: "Alice", Age: 30}        │
   │ Creates user object in memory                        │
   └──────────────────────────────────────────────────────┘
                           ↓
   ┌──────────────────────────────────────────────────────┐
   │ store.Create(user)                                   │
   │ Stores: store.users["1"] = User{...}                 │
   │ State: Store now contains 1 user                     │
   └──────────────────────────────────────────────────────┘

2. REQUEST CREATION (Lines 180-189)
   ┌──────────────────────────────────────────────────────┐
   │ req := httptest.NewRequest(GET, "/users/1", nil)     │
   │ Creates mock HTTP request:                           │
   │   - Method: GET                                      │
   │   - URL: "/users/1"                                  │
   │   - Body: nil                                        │
   └──────────────────────────────────────────────────────┘
                           ↓
   ┌──────────────────────────────────────────────────────┐
   │ rr := httptest.NewRecorder()                         │
   │ Creates response recorder to capture:                │
   │   - Status code (rr.Code)                            │
   │   - Headers (rr.Header())                            │
   │   - Body (rr.Body)                                   │
   └──────────────────────────────────────────────────────┘
                           ↓
   ┌──────────────────────────────────────────────────────┐
   │ handler := handleGetUser(store)                      │
   │ Returns closure function with access to 'store'      │
   └──────────────────────────────────────────────────────┘

3. REQUEST SERVING (Line 199 - Critical!)
   ┌──────────────────────────────────────────────────────┐
   │ handler.ServeHTTP(rr, req)                           │
   │ This triggers the actual handler execution...        │
   └──────────────────────────────────────────────────────┘
                           ↓
   ┌─────────────────────────────────────────────────────────────┐
   │           INSIDE handleGetUser FUNCTION                     │
   │  (from main.go lines 156-183)                               │
   ├─────────────────────────────────────────────────────────────┤
   │  Step 1: Check method                                       │
   │    if r.Method != http.MethodGet { ... }                    │
   │    ✓ Passes (method is GET)                                 │
   │                                                              │
   │  Step 2: Extract ID from URL                                │
   │    id := strings.TrimPrefix(r.URL.Path, "/users/")          │
   │    URL Path: "/users/1"                                     │
   │    After TrimPrefix: "1"                                    │
   │    id = "1"                                                 │
   │                                                              │
   │  Step 3: Validate ID                                        │
   │    if id == "" { ... }                                      │
   │    ✓ Passes (id = "1", not empty)                           │
   │                                                              │
   │  Step 4: Get user from store                                │
   │    user, err := store.Get(id)                               │
   │    ┌────────────────────────────────────────┐               │
   │    │ INSIDE store.Get("1")                  │               │
   │    │ - Locks read mutex                     │               │
   │    │ - Looks up store.users["1"]            │               │
   │    │ - Finds: User{ID:"1", Name:"Alice"}    │               │
   │    │ - Returns user, nil error              │               │
   │    └────────────────────────────────────────┘               │
   │    Result: user = User{ID:"1", Name:"Alice", Age:30}        │
   │            err = nil                                        │
   │                                                              │
   │  Step 5: Check for errors                                   │
   │    if err != nil { ... }                                    │
   │    ✓ Passes (err is nil, user was found)                    │
   │                                                              │
   │  Step 6: Set response header                                │
   │    w.Header().Set("Content-Type", "application/json")       │
   │    Writes to rr.Header()                                    │
   │                                                              │
   │  Step 7: Encode and write response                          │
   │    json.NewEncoder(w).Encode(user)                          │
   │    ┌────────────────────────────────────────┐               │
   │    │ Converts User struct to JSON:          │               │
   │    │ {"id":"1","name":"Alice","age":30}     │               │
   │    │ Writes to rr.Body                      │               │
   │    │ Auto-sets status: 200 OK               │               │
   │    └────────────────────────────────────────┘               │
   └─────────────────────────────────────────────────────────────┘
                           ↓
   ┌──────────────────────────────────────────────────────┐
   │ Handler execution complete!                          │
   │ Response recorder (rr) now contains:                 │
   │   - Code: 200                                        │
   │   - Header: {"Content-Type": ["application/json"]}   │
   │   - Body: {"id":"1","name":"Alice","age":30}         │
   └──────────────────────────────────────────────────────┘

4. RESPONSE VERIFICATION (Lines 207-221)
   ┌──────────────────────────────────────────────────────┐
   │ if rr.Code != http.StatusOK { ... }                  │
   │ Checks: rr.Code == 200? ✓ YES                        │
   └──────────────────────────────────────────────────────┘
                           ↓
   ┌──────────────────────────────────────────────────────┐
   │ var got User                                         │
   │ json.NewDecoder(rr.Body).Decode(&got)                │
   │ Parses JSON from rr.Body:                            │
   │   Input: {"id":"1","name":"Alice","age":30}          │
   │   Output: got = User{ID:"1", Name:"Alice", Age:30}   │
   └──────────────────────────────────────────────────────┘
                           ↓
   ┌──────────────────────────────────────────────────────┐
   │ TEST PASSES ✓                                        │
   │ - Status code is 200                                 │
   │ - JSON decoding successful                           │
   └──────────────────────────────────────────────────────┘
```

## Key Points About Data Flow

### 1. **The Store (In-Memory Database)**
```go
store := NewSimpleUserStore()
// Creates: {users: map[string]User{}, mu: sync.RWMutex{}}

store.Create(user)
// After this: {users: map[string]User{"1": User{ID:"1",...}}}
```

### 2. **Request Flow**
The request goes through these objects:
- `req` (mock request) → contains URL "/users/1"
- `handler` (closure) → has access to `store`
- `rr` (response recorder) → captures everything handler writes

### 3. **Handler Execution**
When `handler.ServeHTTP(rr, req)` is called:
1. Handler reads from `req.URL.Path` to get ID
2. Handler calls `store.Get(id)` to fetch user
3. Handler writes to `rr` (status, headers, body)

### 4. **The Critical Line**
```go
handler.ServeHTTP(rr, req)
```
This is where EVERYTHING happens! Before this line:
- Store has data
- Request is ready
- Response recorder is empty

After this line:
- Response recorder has status, headers, and body
- Test can verify the results

### 5. **Why This Pattern?**
This testing pattern allows you to:
- Test HTTP handlers WITHOUT running a real HTTP server
- Control exact input (request)
- Inspect exact output (response)
- Run tests in milliseconds

## Complete Flow Summary

```
Test Data → Store → Handler → Response → Verification
   ↓         ↓        ↓          ↓          ↓
  User    stores   processes  captures   checks
 Created  in map    request    output     output
```

1. **Test creates user and stores it**
2. **Test creates fake HTTP request**
3. **Handler processes request:**
   - Extracts ID from URL
   - Queries store
   - Encodes user to JSON
   - Writes response
4. **Test verifies response is correct**

This is unit testing for HTTP handlers - testing the handler's logic in isolation!
