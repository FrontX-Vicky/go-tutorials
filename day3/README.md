# Day 3 — Error Handling, Pointers, Slices & Maps

## Concepts to Learn

### 1. Error Handling

Go doesn't use exceptions—functions return errors as values.

- **Pattern**: Return `(result, error)` from functions that can fail.
- **Check errors explicitly**:
  ```go
  result, err := someFunction()
  if err != nil {
      // handle error
      return err
  }
  // use result
  ```
- **Custom errors**:
  ```go
  import "errors"
  
  func Divide(a, b float64) (float64, error) {
      if b == 0 {
          return 0, errors.New("division by zero")
      }
      return a / b, nil
  }
  ```

### 2. Pointers Deep Dive

- **Why pointers?** To modify data in place without copying large structs, and to represent "optional" values (nil).
- **Value vs pointer receiver**: Use pointer receivers when you need to modify the receiver or avoid copying large structs.
- **nil pointers**: Always check for nil before dereferencing.
  ```go
  var user *User
  if user != nil {
      fmt.Println(user.Name)
  }
  ```

### 3. Slices

- **Dynamic arrays**: Slices grow as needed.
- **Common operations**:
  ```go
  slice := []int{1, 2, 3}
  slice = append(slice, 4)        // add element
  slice = slice[1:]               // remove first element
  slice = slice[:len(slice)-1]    // remove last element
  ```
- **Slices are references**: Modifying a slice affects the underlying array.

### 4. Maps

- **Key-value storage**:
  ```go
  ages := make(map[string]int)
  ages["Alice"] = 30
  age, exists := ages["Alice"]  // check if key exists
  delete(ages, "Alice")         // remove key
  ```
- **Iterate over maps**:
  ```go
  for key, value := range ages {
      fmt.Printf("%s: %d\n", key, value)
  }
  ```

## Tasks

### 1. Build a User Registry with Error Handling

Create a `UserRegistry` struct that stores users by ID:

- `AddUser(id int, user User) error`: Add a user. Return error if ID already exists.
- `GetUser(id int) (*User, error)`: Get a user by ID. Return error if not found.
- `UpdateUser(id int, user User) error`: Update user. Return error if not found.
- `DeleteUser(id int) error`: Delete user. Return error if not found.
- `ListUsers() []User`: Return all users.

### 2. Implement a Simple Calculator

Create functions that return errors for invalid operations:

- `Divide(a, b float64) (float64, error)`: Return error on division by zero.
- `Sqrt(x float64) (float64, error)`: Return error for negative numbers.

### 3. Work with Slices

Write a function `RemoveDuplicates(slice []int) []int` that removes duplicate values from a slice.

### 4. Work with Maps

Write a function `CountWords(text string) map[string]int` that counts word frequency in a string.

## Extra Challenge

Implement a `Safe` wrapper that handles panics:

```go
func Safe(fn func() error) error {
    // Use defer and recover to catch panics
    // Return the error if no panic, or convert panic to error
}
```

Test it with a function that might panic.
