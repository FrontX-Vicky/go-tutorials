package main

import (
	"reflect"
	"testing"
)

// Day 9: Generics - Answer Key
// Reference implementations for all exercises

// ============================================
// Test 1: Stack - Basic Operations
// ============================================

func TestStack_Basic_Answer(t *testing.T) {
	stack := NewStack[int]()
	if stack.Len() != 0 {
		t.Fatalf("expected empty stack, got len=%d", stack.Len())
	}

	stack.Push(1)
	stack.Push(2)
	if stack.Len() != 2 {
		t.Fatalf("expected len=2, got len=%d", stack.Len())
	}

	if top, ok := stack.Peek(); !ok || top != 2 {
		t.Fatalf("expected peek=2 ok=true, got %v ok=%v", top, ok)
	}

	if val, ok := stack.Pop(); !ok || val != 2 {
		t.Fatalf("expected pop=2 ok=true, got %v ok=%v", val, ok)
	}

	if val, ok := stack.Pop(); !ok || val != 1 {
		t.Fatalf("expected pop=1 ok=true, got %v ok=%v", val, ok)
	}

	if _, ok := stack.Pop(); ok {
		t.Fatalf("expected pop on empty to return ok=false")
	}
}

// ============================================
// Test 2: Queue - FIFO Behavior
// ============================================

func TestQueue_Basic_Answer(t *testing.T) {
	queue := NewQueue[string]()
	queue.Enqueue("a")
	queue.Enqueue("b")
	queue.Enqueue("c")

	if queue.Len() != 3 {
		t.Fatalf("expected len=3, got len=%d", queue.Len())
	}

	if val, ok := queue.Dequeue(); !ok || val != "a" {
		t.Fatalf("expected dequeue=a ok=true, got %q ok=%v", val, ok)
	}

	if val, ok := queue.Dequeue(); !ok || val != "b" {
		t.Fatalf("expected dequeue=b ok=true, got %q ok=%v", val, ok)
	}

	if val, ok := queue.Dequeue(); !ok || val != "c" {
		t.Fatalf("expected dequeue=c ok=true, got %q ok=%v", val, ok)
	}

	if queue.Len() != 0 {
		t.Fatalf("expected len=0, got len=%d", queue.Len())
	}

	if _, ok := queue.Dequeue(); ok {
		t.Fatalf("expected dequeue on empty to return ok=false")
	}
}

// ============================================
// Test 3: Set - Uniqueness
// ============================================

func TestSet_Basic_Answer(t *testing.T) {
	set := NewSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(2)
	set.Add(3)

	if set.Len() != 3 {
		t.Fatalf("expected len=3, got len=%d", set.Len())
	}

	if !set.Has(2) {
		t.Fatalf("expected set to contain 2")
	}

	if set.Has(4) {
		t.Fatalf("did not expect set to contain 4")
	}

	set.Remove(2)
	if set.Has(2) {
		t.Fatalf("expected 2 to be removed")
	}

	if set.Len() != 2 {
		t.Fatalf("expected len=2, got len=%d", set.Len())
	}
}

// ============================================
// Test 4: Map, Filter, Reduce
// ============================================

func TestSliceUtilities_MapFilterReduce_Answer(t *testing.T) {
	numbers := []int{1, 2, 3, 4}

	squared := MapSlice(numbers, func(n int) int { return n * n })
	expectedSquares := []int{1, 4, 9, 16}
	if !reflect.DeepEqual(squared, expectedSquares) {
		t.Fatalf("expected %v, got %v", expectedSquares, squared)
	}

	evens := FilterSlice(numbers, func(n int) bool { return n%2 == 0 })
	expectedEvens := []int{2, 4}
	if !reflect.DeepEqual(evens, expectedEvens) {
		t.Fatalf("expected %v, got %v", expectedEvens, evens)
	}

	sum := ReduceSlice(numbers, 0, func(acc, n int) int { return acc + n })
	if sum != 10 {
		t.Fatalf("expected sum=10, got %d", sum)
	}
}

// ============================================
// Test 5: Contains and Unique
// ============================================

func TestSliceUtilities_ContainsUnique_Answer(t *testing.T) {
	words := []string{"go", "rust", "go", "zig"}

	if !Contains(words, "go") {
		t.Fatalf("expected Contains to find 'go'")
	}

	if Contains(words, "java") {
		t.Fatalf("did not expect Contains to find 'java'")
	}

	unique := Unique(words)
	expected := []string{"go", "rust", "zig"}
	if !reflect.DeepEqual(unique, expected) {
		t.Fatalf("expected %v, got %v", expected, unique)
	}
}

// ============================================
// Test 6: Min and Max
// ============================================

func TestSliceUtilities_MinMax_Answer(t *testing.T) {
	ints := []int{5, 2, 9, 1}
	min, ok := Min(ints)
	if !ok || min != 1 {
		t.Fatalf("expected min=1 ok=true, got %v ok=%v", min, ok)
	}

	max, ok := Max(ints)
	if !ok || max != 9 {
		t.Fatalf("expected max=9 ok=true, got %v ok=%v", max, ok)
	}

	words := []string{"beta", "alpha", "gamma"}
	minWord, ok := Min(words)
	if !ok || minWord != "alpha" {
		t.Fatalf("expected min=alpha ok=true, got %v ok=%v", minWord, ok)
	}

	maxWord, ok := Max(words)
	if !ok || maxWord != "gamma" {
		t.Fatalf("expected max=gamma ok=true, got %v ok=%v", maxWord, ok)
	}

	var empty []int
	if _, ok := Min(empty); ok {
		t.Fatalf("expected Min on empty slice to return ok=false")
	}
	if _, ok := Max(empty); ok {
		t.Fatalf("expected Max on empty slice to return ok=false")
	}
}
