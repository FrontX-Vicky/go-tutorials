# Day 2 — Structs, Methods, Interfaces

## Concepts to Learn

### 1. Structs
- A `struct` is a composite data type that groups together fields.
- Example:
  ```go
  type User struct {
      Name  string
      Age   int
      Email string
  }
  ```

### 2. Methods
- Functions associated with a `struct`.
- **Value receiver**: Method gets a copy of the struct.
- **Pointer receiver**: Method gets a pointer to the struct, allowing modification.
- Example:
  ```go
  func (u User) Greet() string {
      return "Hello, " + u.Name
  }

  func (u *User) IncrementAge() {
      u.Age++
  }
  ```

### 3. Interfaces
- Define a set of methods that a type must implement.
- Example:
  ```go
  type Profile interface {
      PrintProfile()
  }

  func (u User) PrintProfile() {
      fmt.Printf("Name: %s, Age: %d\n", u.Name, u.Age)
  }
  ```

### 4. Composition
- Embedding one struct into another.
- Example:
  ```go
  type Employee struct {
      User
      Position string
  }
  ```

## Tasks

1. **Create a `User` struct**:
   - Fields: `Name` (string), `Age` (int), `Skills` ([]string).

2. **Add methods to `User`**:
   - `Greet() string`: Returns a greeting message.
   - `IsAdult() bool`: Returns `true` if `Age` >= 18.
   - `AddSkill(skill string)`: Adds a skill to the `Skills` slice (use a pointer receiver).

3. **Define a `Profile` interface**:
   - Method: `PrintProfile()`.
   - Implement `Profile` for `User`.

4. **(Optional) Create an `Employee` struct**:
   - Embed `User`.
   - Add `Position` (string).
   - Override `PrintProfile()` to include `Position`.

## Extra Challenge
- Write a function `PrintProfiles(profiles []Profile)` that takes a slice of `Profile` and calls `PrintProfile()` on each.
