# Learning Log

## Day 1 — Go Basics
**Date**: November 17, 2025

**Concepts Covered**:
- Basic syntax and structure
- Variables and types
- Control flow (if, for, switch)
- Functions

**Files Created**:
- `day1/main.go`

**Key Takeaways**:
- Go's simple syntax and explicit error handling
- No semicolons needed
- Package system and imports

---

## Day 2 — Structs, Methods, Interfaces
**Date**: November 19, 2025

**Concepts Covered**:
- Structs as composite types
- Methods with value and pointer receivers
- Interfaces for polymorphism
- Struct composition (embedding)

**Files Created/Updated**:
- `day2/README.md`
- `day2/main.go`

**Key Takeaways**:
- Use pointer receivers when you need to modify the receiver or avoid copying
- Interfaces are implemented implicitly (no `implements` keyword)
- Composition over inheritance through embedding
- Interface satisfaction is checked at compile time

**Tasks Completed**:
- Created `User` struct with `Name`, `Age`, and `Skills` fields
- Implemented methods: `Greet()`, `IsAdult()`, `AddSkill()`
- Defined and implemented `Profile` interface
- Created `Employee` struct with embedded `User`
- Implemented `PrintProfiles()` helper for polymorphic behavior

---

## Day 3 — Error Handling, Pointers, Slices & Maps
**Date**: November 19-20, 2025

**Concepts Covered**:
- Error handling patterns (no exceptions)
- Pointers and nil checks
- Slice operations and behavior
- Map operations and iteration

**Files Created**:
- `day3/README.md`
- `day3/main.go`

**Key Takeaways**:
- Go uses error values instead of exceptions—always check `if err != nil`
- Must initialize maps with `make()` before use (nil maps panic on write)
- Map-based deduplication (O(n)) is more idiomatic than nested loops
- `strings.Fields()` is the standard way to split text into words
- `defer` + `recover()` can catch panics and convert them to errors

**Tasks Completed**:
- Built `UserRegistry` with CRUD operations (AddUser, GetUser, UpdateUser, DeleteUser, ListUsers)
- Implemented `Divide` and `Sqrt` functions with proper error handling
- Created `RemoveDuplicates` using map for O(n) performance
- Implemented `CountWords` using `strings.Fields()` and map
- Implemented `Safe` wrapper to handle panics gracefully

---

## Day 4 — Concurrency: Goroutines, Channels & Sync
**Date**: November 20, 2025

**Concepts Covered**:
- Goroutines (lightweight threads)
- Channels (unbuffered and buffered)
- Select statement for multiplexing
- WaitGroups for synchronization
- Concurrency patterns (worker pool, pipeline, fan-out/fan-in)

**Files Created**:
- `day4/README.md`
- `day4/main.go`

**Key Takeaways**:
- (To be filled after completing Day 4)

**Tasks to Complete**:
- Launch goroutines and coordinate with channels
- Understand buffered vs unbuffered channels
- Implement worker pool pattern
- Use select with timeout
- (Optional) Build a pipeline with multiple stages
