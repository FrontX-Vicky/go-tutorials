package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Helper functions to test

func Add(a, b int) int {
	return a + b
}

func IsPalindrome(s string) bool {
	s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// UserStore from Day 5 (simplified for testing)

type User struct {
	ID   string
	Name string
	Age  int
}

type UserStore struct {
	users map[string]User
	mu    sync.RWMutex
}

type Response struct {
	Message string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]User),
	}
}

func (us *UserStore) Create(ctx context.Context, user User) Response {
	select {
	case <-ctx.Done():
		return Response{Message: "operation cancelled", ErrorCode: 2}
	default:
	}

	time.Sleep(100 * time.Millisecond)
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.users[user.ID]; exists {
		return Response{Message: fmt.Sprintf("user with id %s already exists", user.ID), ErrorCode: 1}
	}
	us.users[user.ID] = user
	return Response{Message: "user created successfully", ErrorCode: 0}
}

func (us *UserStore) Get(ctx context.Context, id string) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}

	us.mu.RLock()
	defer us.mu.RUnlock()

	user, exists := us.users[id]
	if !exists {
		return User{}, fmt.Errorf("user with id %s not found", id)
	}
	return user, nil
}

func (us *UserStore) List(ctx context.Context) ([]User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	us.mu.RLock()
	defer us.mu.RUnlock()

	users := []User{}
	for _, user := range us.users {
		users = append(users, user)
	}
	return users, nil
}

func (us *UserStore) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.users[id]; !exists {
		return fmt.Errorf("user with id %s not found", id)
	}
	delete(us.users, id)
	return nil
}

func main() {
	fmt.Println("Day 6: Testing - Run 'go test -v' to execute tests")
}
