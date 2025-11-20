package main

import (
	"fmt"
	"math"
	"strings"
)

// TODO: Define the User struct (reuse from Day 2 or simplify)
type User struct {
	Name   string
	Age    int
}

// TODO: Define the UserRegistry struct with a map to store users by ID
type UserRegistry struct {
	users map[string]User
}

// TODO: Implement AddUser method - return error if ID already exists
func (ur *UserRegistry) AddUser(id string, user User) error {
	if _, exists := ur.users[id]; exists {
		return fmt.Errorf("User with ID %s already exists", id)
	}
	ur.users[id] = user
	return nil
}

// TODO: Implement GetUser method - return error if user not found
func (ur *UserRegistry) GetUser(id string) (User, error) {
	if user, exists := ur.users[id]; exists {
		return user, nil
	} 
	return User{}, fmt.Errorf("User with Id %s not found", id)
}

// TODO: Implement UpdateUser method - return error if user not found
func (ur *UserRegistry) UpdateUser(id string, updated User) error {
	if _, exists := ur.users[id]; exists {
		ur.users[id] = updated
		return nil
	}
	return fmt.Errorf("User with Id %s not found", id)
}

// TODO: Implement DeleteUser method - return error if user not found
func (ur *UserRegistry) DeleteUser(id string) error {
	if _, exists := ur.users[id]; exists {
		delete(ur.users, id)
		return nil
	}
	return fmt.Errorf("User with Id %s not found", id)
}

// TODO: Implement ListUsers method - return slice of all users
func (ur *UserRegistry) ListUsers() []User {
	users := []User{}
	for _, user := range ur.users {
		users = append(users, user)
	}
	return users
}

// TODO: Implement Divide function - return error on division by zero
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("Cannot divide by 0")
	}
	return a / b, nil
}

// TODO: Implement Sqrt function - return error for negative numbers (use math.Sqrt)
func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, fmt.Errorf("cannot compute square root of negative numbers")
	}
	return math.Sqrt(x), nil
}

// compare and tell which is better
// TODO: Implement RemoveDuplicates function
func RemoveDuplicates1(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, num := range input {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func RemoveDuplicates2(input []int) []int {
	seen := []int{}
	for _, num := range input {
		for _, v := range seen {
			if v == num {
				goto next
			}
		}
		seen = append(seen, num)
		next:
	}
	return seen
}

// TODO: Implement CountWords function
func CountWords(text string) map[string]int {
	wordCount := make(map[string]int)
	for _, word := range strings.Fields(text) {
		wordCount[word]++
	}
	return wordCount
}

// TODO (Optional): Implement Safe function that handles panics
func Safe(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	f()
	return err
}

func main() {
	fmt.Println("=== Day 3: Error Handling, Pointers, Slices & Maps ===\n")

	// TODO: Test UserRegistry
	// - Create registry
	registry := UserRegistry{users: make(map[string]User)} // explain why make is used here
	// - Add users
	registry.AddUser("1", User{Name: "Alice", Age: 30})
	registry.AddUser("2", User{Name: "Bob", Age: 25})
	// - Try to add duplicate (should error)
	err := registry.AddUser("1", User{Name: "Charlie", Age: 28})
	if err != nil {
		fmt.Println(err)
	}
	// - Get user
	var id = "1"
	user, err := registry.GetUser(id)
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println("User ", id, ":", user)
	}
	id = "10"
	user, err = registry.GetUser(id)
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println("User ", id, ":", user)
	}

	// - Update user
	id = "1"
	err = registry.UpdateUser(id, User{Name : "vicky", Age: 26})
	if err != nil {
		fmt.Println(err)
	}
	id = "10"
	err = registry.UpdateUser(id, User{Name : "vicky", Age: 26})
	if err != nil {
		fmt.Println(err)
	}

	// - Delete user
	id = "2"
	err = registry.DeleteUser(id)
	if err != nil {
		fmt.Println(err)
	}
	id = "10"
	err = registry.DeleteUser(id)
	if err != nil {
		fmt.Println(err)
	}
	// - List all users
	users := registry.ListUsers()
	fmt.Println("All users in registry :")
	for _, user := range users {
		fmt.Println(user)
	}

	// TODO: Test Calculator functions
	// - Test Divide with valid and invalid inputs
	result, err := Divide(10, 2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	result, err = Divide(10, 0)
	if err != nil {
		fmt.Println(err)
	}
	// - Test Sqrt with valid and invalid inputs
	result, err = Sqrt(8)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	result, err = Sqrt(-1)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Test RemoveDuplicates
	// - Test with slice containing duplicates
	nums := []int{1,2,2,3,3,4,5,5,6,6,6,6}
	unniqueNums1 := RemoveDuplicates1(nums)
	fmt.Println("Method1", unniqueNums1)
	
	unniqueNums2 := RemoveDuplicates2(nums)
	fmt.Println("Method2", unniqueNums2)

	// TODO: Test CountWords
	// - Test with sample text
	text := "hello world hello Go Go Go"
	wordCount := CountWords(text)
	fmt.Println(wordCount)

	// TODO (Optional): Test Safe function
	err = Safe(func() {
		fmt.Println("Inside Safe function")
		panic("something went wrong!")
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nDay 3 tasks completed!")
}
