package main

import (
	"context"
	"testing"
	"time"
)

// TODO: Write table-driven tests for helper functions

func TestAdd(t *testing.T) { //whats the point to add t *testing.T package here?
	// TODO: Create test cases slice
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -1, -1, -2},
		{"zero", 0, 5, 5},
	}
	// results := []int{}
	// TODO: Loop through cases with t.Run
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // whats the point to add t *testing.T package here again?
			result := Add(tt.a, tt.b)
			// results = append(results, result)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
	// TODO: Compare result with expected
	// for i, tt := range tests {
	// 	if results[i] != tt.expected {
	// 		t.Errorf("Test %s failed: expected %d, got %d", tt.name, tt.expected, results[i])
	// 	}
	// }
}

func TestIsPalindrome(t *testing.T) {
	// TODO: Test cases: "racecar", "hello", "A man a plan a canal Panama"
	tests := []struct {
		input string
		want  bool
	}{
		{"racecar", true},
		{"hello", false},
		{"A man a plan a canal Panama", true},
		{"", true},
		{"a", true},
	}
	// TODO: Test empty string, single character
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) { // why to use t.Run here?
			got := IsPalindrome(tt.input)
			if got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	// TODO: Test cases: "hello" -> "olleh", "Go" -> "oG", ""
	tests := []struct {
		input string
		want string}{
			{"hello", "olleh"},
			{"Go", "oG"},
			{"", ""},
		}
	
	// TODO: Test Unicode strings
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Reverse(tt.input)
			if got != tt.want {
				t.Errorf("Reverse(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	// TODO: Test valid emails: "test@example.com"
	tests := []struct {
		input string
		want  bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
		{"user@domain", false},
		{"user@domain.com", true},
	}
	// TODO: Test invalid: "", "notanemail", "missing@domain"
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			got := err == nil
			if got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TODO: Write table-driven tests for UserStore

func TestUserStore_Create(t *testing.T) {
	// TODO: Test success case

	CreatedUserStore := NewUserStore()
	test := []User{
		{ID: "1", Name: "Alice", Age: 30},
	}
	t.Run("create", func(t *testing.T) {
		// Test logic here
		for _, user := range test {
			response := CreatedUserStore.Create(context.Background(), user)
			if response.ErrorCode != 0 {
				t.Errorf("Create(%v) unexpected error: %v", user, response.Message)
			}
		}
	})
	// TODO: Test duplicate ID error
	t.Run("Duplicate ID", func(t *testing.T) {
		// Test logic here
		for _, user := range test {
			response := CreatedUserStore.Create(context.Background(), user)
			if response.ErrorCode != 0 {
				t.Errorf("Create(%v) unexpected error: %v", user, response.Message)
			}
		}
	})
	
	// TODO: Test context cancellation
	t.Run("Context Cancellation", func(t *testing.T) {
		CreatedUserStore := NewUserStore()
		ctx, cancel := createTimeoutContext(1)
		defer cancel()
		response := CreatedUserStore.Create(ctx, User{ID: "1", Name: "Alice", Age: 30})
		if response.ErrorCode != 2 {
			t.Errorf("Create with cancelled context: expected error code 2, got %d", response.ErrorCode)
		}
	})
}

func TestUserStore_Get(t *testing.T) {
	// TODO: Test success case
	// TODO: Test not found error
	// TODO: Test context cancellation
}

func TestUserStore_List(t *testing.T) {
	// TODO: Test empty store
	// TODO: Test with multiple users
	// TODO: Test context cancellation
}

func TestUserStore_Delete(t *testing.T) {
	// TODO: Test success case
	// TODO: Test not found error
	// TODO: Test context cancellation
}

// TODO: Write benchmarks

func BenchmarkAdd(b *testing.B) {
	// TODO: Benchmark Add function
}

func BenchmarkIsPalindrome(b *testing.B) {
	// TODO: Benchmark with long palindrome
}

func BenchmarkUserStore_Create(b *testing.B) {
	// TODO: Benchmark store creation
}

func BenchmarkUserStore_Get(b *testing.B) {
	// TODO: Setup: create store with users
	// TODO: Benchmark Get operation
}

// Helper function example
func createTestUser(id, name string, age int) User {
	return User{ID: id, Name: name, Age: age}
}

// Context with timeout helper
func createTimeoutContext(ms int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(ms)*time.Millisecond)
}
