package main

import (
	"fmt"
)

// TODO: Define the User struct (reuse from Day 2 or simplify)
type User struct {
	Name   string
	Age    int
	Skills []string
}

// TODO: Define the UserRegistry struct with a map to store users by ID

// TODO: Implement AddUser method - return error if ID already exists

// TODO: Implement GetUser method - return error if user not found

// TODO: Implement UpdateUser method - return error if user not found

// TODO: Implement DeleteUser method - return error if user not found

// TODO: Implement ListUsers method - return slice of all users

// TODO: Implement Divide function - return error on division by zero

// TODO: Implement Sqrt function - return error for negative numbers (use math.Sqrt)

// TODO: Implement RemoveDuplicates function

// TODO: Implement CountWords function

// TODO (Optional): Implement Safe function that handles panics

func main() {
	fmt.Println("=== Day 3: Error Handling, Pointers, Slices & Maps ===\n")

	// TODO: Test UserRegistry
	// - Create registry
	// - Add users
	// - Try to add duplicate (should error)
	// - Get user
	// - Update user
	// - Delete user
	// - List all users

	// TODO: Test Calculator functions
	// - Test Divide with valid and invalid inputs
	// - Test Sqrt with valid and invalid inputs

	// TODO: Test RemoveDuplicates
	// - Test with slice containing duplicates

	// TODO: Test CountWords
	// - Test with sample text

	// TODO (Optional): Test Safe function

	fmt.Println("\nDay 3 tasks completed!")
}
