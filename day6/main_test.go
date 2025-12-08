package main

import (
	"context"
	"fmt"
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
		want  string
	}{
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

	ctx := context.Background()

	test := []User{
		{ID: "1", Name: "Alice", Age: 30},
	}
	t.Run("create", func(t *testing.T) {
		// Test logic here
		for _, user := range test {
			response := CreatedUserStore.Create(ctx, user)
			if response.ErrorCode != 0 {
				t.Errorf("Create(%v) unexpected error: %v", user, response.Message)
			}
		}
	})
	// TODO: Test duplicate ID error
	t.Run("Duplicate ID", func(t *testing.T) {
		// Test logic here
		for _, user := range test {
			response := CreatedUserStore.Create(ctx, user)
			if response.ErrorCode == 0 || response.ErrorCode == 2 {
				t.Errorf("Create(%v) unexpected error: %v", user, response.Message)
			}
		}
	})

	// TODO: Test context cancellation
	t.Run("Context Cancellation", func(t *testing.T) { // Why this is not working as expected?
		CreatedUserStore := NewUserStore()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately before calling Create
		response := CreatedUserStore.Create(ctx, User{ID: "1", Name: "Alice", Age: 30})
		if response.ErrorCode != 2 {
			t.Errorf("Create with cancelled context: expected error code 2, got %d", response.ErrorCode)
		}
	})
}

func TestUserStore_Get(t *testing.T) {

	CreatedUserStore := NewUserStore()
	ctx, cancel := createTimeoutContext(1)
	user := User{ID: "1", Name: "Alice", Age: 30}
	defer cancel()
	CreatedUserStore.Create(ctx, user)
	// TODO: Test success case
	t.Run("Success Case", func(t *testing.T) {
		userGot, err := CreatedUserStore.Get(context.Background(), "1")
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if userGot != user {
			t.Errorf("Get returned %+v; want %+v", userGot, user)
		}
	})
	// TODO: Test not found error

	t.Run("Not Found Error", func(t *testing.T) {
		_, err := CreatedUserStore.Get(context.Background(), "10")
		if err == nil {
			t.Errorf("Get for non-existent user should fail, got nil error")
		}
	})
	// TODO: Test context cancellation
	t.Run("Context Cancellation", func(t *testing.T) { // Why this is not working as expected?
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		_, err := CreatedUserStore.Get(ctx, "1")
		if err == nil {
			t.Errorf("Get with cancelled context should fail, got nil error")
		}
	})
}

func TestUserStore_List(t *testing.T) {
	// TODO: Test empty store
	UserStore := NewUserStore()
	users, err := UserStore.List(context.Background())
	if err != nil {
		t.Errorf("list is not empty")
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
	// TODO: Test with multiple users
	UserStore.Create(context.Background(), User{ID: "1", Name: "Alice", Age: 30})
	UserStore.Create(context.Background(), User{ID: "2", Name: "Bob", Age: 25})
	t.Run("Multiple Users", func(t *testing.T) {
		users, err := UserStore.List(context.Background())
		if err != nil {
			t.Errorf("Unexpected error %d", err)
		}
		if len(users) != 2 {
			t.Errorf("In users expected count is 2 but found %d", len(users))
		}
	})
	// TODO: Test context cancellation
}

func TestUserStore_Delete(t *testing.T) {
	// TODO: Test success case
	UserStore := NewUserStore()
	ctx_bg := context.Background()
	UserStore.Create(ctx_bg, User{ID: "1", Name: "Alice", Age: 30})
	t.Run("Success Case", func(t *testing.T) {
		err := UserStore.Delete(ctx_bg, "1")
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})
	// TODO: Test not found error
	t.Run("Not Found Error", func(t *testing.T) {
		err := UserStore.Delete(ctx_bg, "10")
		if err == nil {
			t.Errorf("User Deleted hence case faild")
		}
	})
	// TODO: Test context cancellation
	t.Run("Context Cancellation", func(t *testing.T) { // Why this is not working as expected?
		UserStore := NewUserStore()
		UserStore.Create(context.Background(), User{ID: "1", Name: "Alice", Age: 30})
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		err := UserStore.Delete(ctx, "1")
		if err == nil {
			t.Errorf("Delete with cancelled context should fail, got nil error")
		}
	})
}

// TODO: Write benchmarks // why benchmarks are useful?

func BenchmarkAdd(b *testing.B) {
	// TODO: Benchmark Add function
	for i := 0; i < b.N; i++ {
		Add(i, i)
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	// TODO: Benchmark with long palindrome
	for i := 0; i < b.N; i++ {
		IsPalindrome("A man a plan a canal Panama")
	}
}

func BenchmarkUserStore_Create(b *testing.B) {
	// TODO: Benchmark store creation
	store := NewUserStore()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		user := User{ID: fmt.Sprintf("%d", i), Name: "User", Age: 20}
		store.Create(ctx, user)
	}
}

func BenchmarkUserStore_Get(b *testing.B) {
	// TODO: Setup: create store with users
	store := NewUserStore()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		user := User{ID: fmt.Sprintf("%d", i), Name: "User", Age: 20}
		store.Create(ctx, user)
	}
	// TODO: Benchmark Get operation
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("%d", i)
		store.Get(ctx, id)
	}
}

// Helper function example
func createTestUser(id, name string, age int) User {
	return User{ID: id, Name: name, Age: age}
}

// Context with timeout helper
func createTimeoutContext(ms int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(ms)*time.Millisecond)
}
