package main

// import (
// 	"context"
// 	"testing"
// 	"time"
// )

// func TestAdd(t *testing.T) {
// 	tests := []struct {
// 		a, b int
// 		want int
// 	}{
// 		{1, 2, 3},
// 		{-1, 1, 0},
// 		{0, 0, 0},
// 	}
// 	for _, tt := range tests {
// 		got := Add(tt.a, tt.b)
// 		if got != tt.want {
// 			t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
// 		}
// 	}
// }

// func TestIsPalindrome(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		want  bool
// 	}{
// 		{"madam", true},
// 		{"racecar", true},
// 		{"hello", false},
// 		{"", true},
// 	}
// 	for _, tt := range tests {
// 		got := IsPalindrome(tt.input)
// 		if got != tt.want {
// 			t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, got, tt.want)
// 		}
// 	}
// }

// func TestReverse(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		want  string
// 	}{
// 		{"hello", "olleh"},
// 		{"Go", "oG"},
// 		{"", ""},
// 	}
// 	for _, tt := range tests {
// 		got := Reverse(tt.input)
// 		if got != tt.want {
// 			t.Errorf("Reverse(%q) = %q; want %q", tt.input, got, tt.want)
// 		}
// 	}
// }

// func TestValidateEmail(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		want  bool
// 	}{
// 		{"test@example.com", true},
// 		{"invalid-email", false},
// 		{"user@domain", false},
// 		{"user@domain.com", true},
// 	}
// 	for _, tt := range tests {
// 		got := ValidateEmail(tt.input)
// 		if got != tt.want {
// 			t.Errorf("ValidateEmail(%q) = %v; want %v", tt.input, got, tt.want)
// 		}
// 	}
// }

// func createTestUser(id, name, email string) User {
// 	return User{ID: id, Name: name, Email: email}
// }

// func createTimeoutContext(d time.Duration) (context.Context, context.CancelFunc) {
// 	return context.WithTimeout(context.Background(), d)
// }

// func TestUserStore_Create_Get_List_Delete(t *testing.T) {
// 	store := NewUserStore()
// 	ctx := context.Background()
// 	user := createTestUser("1", "Alice", "alice@example.com")

// 	// Test Create
// 	if err := store.Create(ctx, user); err != nil {
// 		t.Fatalf("Create failed: %v", err)
// 	}

// 	// Test Get
// 	got, err := store.Get(ctx, "1")
// 	if err != nil {
// 		t.Fatalf("Get failed: %v", err)
// 	}
// 	if got != user {
// 		t.Errorf("Get returned %+v; want %+v", got, user)
// 	}

// 	// Test List
// 	users, err := store.List(ctx)
// 	if err != nil {
// 		t.Fatalf("List failed: %v", err)
// 	}
// 	if len(users) != 1 || users[0] != user {
// 		t.Errorf("List returned %+v; want [%+v]", users, user)
// 	}

// 	// Test Delete
// 	if err := store.Delete(ctx, "1"); err != nil {
// 		t.Fatalf("Delete failed: %v", err)
// 	}
// 	_, err = store.Get(ctx, "1")
// 	if err == nil {
// 		t.Errorf("Get after Delete should fail, got nil error")
// 	}
// }

// func BenchmarkAdd(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		Add(i, i)
// 	}
// }

// func BenchmarkIsPalindrome(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		IsPalindrome("racecar")
// 	}
// }

// func BenchmarkUserStore_Create_Get(b *testing.B) {
// 	store := NewUserStore()
// 	for i := 0; i < b.N; i++ {
// 		id := string(rune(i))
// 		user := createTestUser(id, "Name", "email@example.com")
// 		ctx := context.Background()
// 		_ = store.Create(ctx, user)
// 		_, _ = store.Get(ctx, id)
// 	}
// }
