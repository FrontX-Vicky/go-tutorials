# Day 9 - Generics and Reusable Data Structures

## Concepts to Learn

### 1) Type Parameters
- Define generic types and functions with `[T any]`.
- Type inference lets you omit type arguments in most calls.
- Use generic methods on generic structs.

### 2) Constraints
- `any` means any type.
- `comparable` allows `==` and `!=` (maps/sets/Contains).
- Custom constraints enable ordered comparisons.

### 3) Zero Values in Generics
- Use `var zero T` when you need the zero value for any type.
- Common for empty Pop/Dequeue results.

### 4) Performance Considerations
- Generic code is type-checked at compile time.
- Keep allocations predictable (pre-size slices when possible).

---

## Today's Exercise

You are given generic data structures and helpers in `main.go`.
Your job is to write tests for each component in `main_test.go`.

### Components to Test

1) **Stack**
- Push, Pop, Peek, Len
- Empty behavior should return ok=false

2) **Queue**
- Enqueue/Dequeue FIFO order
- Len decreases correctly
- Empty behavior should return ok=false

3) **Set**
- Add/Remove/Has/Len
- Duplicates should not increase size

4) **Slice Utilities**
- MapSlice, FilterSlice, ReduceSlice
- Contains and Unique
- Min/Max for ordered types

---

## Test Requirements

- Use table-driven tests where it makes sense.
- Check empty inputs and edge cases.
- Keep tests small and focused on behavior.

---

## Commands

```bash
# Run your exercise tests
go test -v

# Run answer tests to see expected behavior
go test -v -run "_Answer$"
```

---

## Files

- `main.go` - Generic data structures and helpers
- `main_test.go` - Exercise tests (TODO)
- `answers_test.go` - Answer key for tests

Good luck with Day 9!
