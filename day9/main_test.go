package main

import (
	"reflect"
	"testing"
)

// Day 9: Generics - Exercises
// Run with: go test -v

// ============================================
// Test 1: Stack - Basic Operations
// ============================================
// TODO: Test Push, Pop, Peek, Len, and empty behavior

func TestStack_Basic(t *testing.T) {
	// t.Skip("TODO: Implement TestStack_Basic")

	// 1. Create a stack and verify initial length
	// 2. Push values and verify Len
	// 3. Peek returns last item without removing
	// 4. Pop returns last item and removes it
	// 5. Pop on empty stack returns ok=false

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
		t.Fatalf("expected peek=2 ok=true, got %v ok =%v", top, ok)
	}

	if val, ok := stack.Pop(); !ok || val != 2 {
		t.Fatalf("expected pop=2 ok=true, got %v ok=%v", val, ok)
	}

	if val, ok := stack.Pop(); !ok || val != 1 {
		t.Fatalf("expeated pop=1 ok=true, got %v ok=%v", val, ok)
	}

	if _, ok := stack.Pop(); ok {
		t.Fatalf("expected pop on empty to return ok=false")
	}

}

// ============================================
// Test 2: Queue - FIFO Behavior
// ============================================
// TODO: Test Enqueue, Dequeue, Len, and empty behavior

func TestQueue_Basic(t *testing.T) {
	// t.Skip("TODO: Implement TestQueue_Basic")

	// 1. Create queue and Enqueue 3 items
	// 2. Dequeue should return in FIFO order
	// 3. Len should decrease on each Dequeue
	// 4. Dequeue on empty queue returns ok=false

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

}

// ============================================
// Test 3: Set - Uniqueness
// ============================================
// TODO: Test Add, Remove, Has, Len

func TestSet_Basic(t *testing.T) {
	//t.Skip("TODO: Implement TestSet_Basic")

	// 1. Create set and add duplicates
	// 2. Len should count unique values only
	// 3. Has should report membership
	// 4. Remove should delete values

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
// TODO: Test MapSlice, FilterSlice, ReduceSlice

func TestSliceUtilities_MapFilterReduce(t *testing.T) {
	t.Skip("TODO: Implement TestSliceUtilities_MapFilterReduce")

	// 1. Map ints to strings or squares
	// 2. Filter only even numbers
	// 3. Reduce to sum
}

// ============================================
// Test 5: Contains and Unique
// ============================================
// TODO: Test Contains and Unique ordering

func TestSliceUtilities_ContainsUnique(t *testing.T) {
	t.Skip("TODO: Implement TestSliceUtilities_ContainsUnique")

	// 1. Contains should return true/false for values
	// 2. Unique should keep first occurrence order
	// 3. Use reflect.DeepEqual for slice comparison
	_ = reflect.DeepEqual
}

// ============================================
// Test 6: Min and Max
// ============================================
// TODO: Test Min/Max with ints and strings

func TestSliceUtilities_MinMax(t *testing.T) {
	t.Skip("TODO: Implement TestSliceUtilities_MinMax")

	// 1. Min/Max on ints
	// 2. Min/Max on strings
	// 3. Empty slice should return ok=false
}
