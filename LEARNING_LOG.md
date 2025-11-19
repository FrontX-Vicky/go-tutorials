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
**Date**: November 19, 2025

**Concepts Covered**:
- Error handling patterns (no exceptions)
- Pointers and nil checks
- Slice operations and behavior
- Map operations and iteration

**Files Created**:
- `day3/README.md`
- `day3/main.go`

**Key Takeaways**:
- (To be filled after completing Day 3)

**Tasks to Complete**:
- Build a `UserRegistry` with CRUD operations and error handling
- Implement calculator functions with error returns
- Work with slices: remove duplicates
- Work with maps: count word frequency
- (Optional) Implement panic recovery with `Safe` wrapper
