package main

import (
	"fmt"
)

// Day 9: Generics and Reusable Data Structures

// Ordered is a constraint for types that can be ordered with < and >.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// ============================================
// 1. Stack - Generic LIFO container
// ============================================

type Stack[T any] struct {
	data []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{data: make([]T, 0)}
}

func (s *Stack[T]) Push(value T) {
	s.data = append(s.data, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

// ============================================
// 2. Queue - Generic FIFO container
// ============================================

type Queue[T any] struct {
	data []T
	head int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{data: make([]T, 0), head: 0}
}

func (q *Queue[T]) Enqueue(value T) {
	q.data = append(q.data, value)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if q.head >= len(q.data) {
		q.data = q.data[:0]
		q.head = 0
		var zero T
		return zero, false
	}
	value := q.data[q.head]
	q.head++

	// Compact the slice occasionally to avoid unbounded growth.
	if q.head > 32 && q.head*2 >= len(q.data) {
		q.data = append([]T(nil), q.data[q.head:]...)
		q.head = 0
	}

	return value, true
}

func (q *Queue[T]) Len() int {
	return len(q.data) - q.head
}

// ============================================
// 3. Set - Generic collection of unique values
// ============================================

type Set[T comparable] struct {
	data map[T]struct{}
}

func NewSet[T comparable](values ...T) *Set[T] {
	set := &Set[T]{data: make(map[T]struct{}, len(values))}
	for _, value := range values {
		set.Add(value)
	}
	return set
}

func (s *Set[T]) Add(value T) {
	s.data[value] = struct{}{}
}

func (s *Set[T]) Remove(value T) {
	delete(s.data, value)
}

func (s *Set[T]) Has(value T) bool {
	_, ok := s.data[value]
	return ok
}

func (s *Set[T]) Len() int {
	return len(s.data)
}

func (s *Set[T]) Values() []T {
	values := make([]T, 0, len(s.data))
	for value := range s.data {
		values = append(values, value)
	}
	return values
}

// ============================================
// 4. Generic Slice Utilities
// ============================================

func MapSlice[T any, R any](in []T, mapper func(T) R) []R {
	out := make([]R, len(in))
	for i, value := range in {
		out[i] = mapper(value)
	}
	return out
}

func FilterSlice[T any](in []T, predicate func(T) bool) []T {
	out := make([]T, 0, len(in))
	for _, value := range in {
		if predicate(value) {
			out = append(out, value)
		}
	}
	return out
}

func ReduceSlice[T any, R any](in []T, init R, reducer func(R, T) R) R {
	acc := init
	for _, value := range in {
		acc = reducer(acc, value)
	}
	return acc
}

func Contains[T comparable](in []T, value T) bool {
	for _, v := range in {
		if v == value {
			return true
		}
	}
	return false
}

func Unique[T comparable](in []T) []T {
	seen := make(map[T]struct{}, len(in))
	out := make([]T, 0, len(in))
	for _, value := range in {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func Min[T Ordered](in []T) (T, bool) {
	if len(in) == 0 {
		var zero T
		return zero, false
	}
	min := in[0]
	for _, value := range in[1:] {
		if value < min {
			min = value
		}
	}
	return min, true
}

func Max[T Ordered](in []T) (T, bool) {
	if len(in) == 0 {
		var zero T
		return zero, false
	}
	max := in[0]
	for _, value := range in[1:] {
		if value > max {
			max = value
		}
	}
	return max, true
}

// ============================================
// Main - Demo
// ============================================

func main() {
	fmt.Println("Day 9: Generics and Reusable Data Structures")
	fmt.Println("Run 'go test -v' to execute tests")
	fmt.Println("")

	stack := NewStack[int]()
	stack.Push(10)
	stack.Push(20)
	if top, ok := stack.Peek(); ok {
		fmt.Printf("Stack peek: %d (len=%d)\n", top, stack.Len())
	}

	queue := NewQueue[string]()
	queue.Enqueue("first")
	queue.Enqueue("second")
	if item, ok := queue.Dequeue(); ok {
		fmt.Printf("Queue dequeue: %s (len=%d)\n", item, queue.Len())
	}

	set := NewSet("go", "go", "generics")
	fmt.Printf("Set has 'go': %v (len=%d)\n", set.Has("go"), set.Len())

	numbers := []int{3, 1, 4, 1, 5}
	squared := MapSlice(numbers, func(n int) int { return n * n })
	unique := Unique(numbers)
	min, _ := Min(numbers)
	max, _ := Max(numbers)

	fmt.Printf("Numbers: %v\n", numbers)
	fmt.Printf("Squared: %v\n", squared)
	fmt.Printf("Unique: %v\n", unique)
	fmt.Printf("Min=%d Max=%d\n", min, max)
}
