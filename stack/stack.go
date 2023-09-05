// Package stack provides an implementation of the stack data structure in Go.
package stack

import (
	"github.com/bongnv/go-container/list"
)

// New creates a new stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{
		container: list.New[T](),
	}
}

// Stack is an implementation of stack.
type Stack[T any] struct {
	container *list.List[T]
}

// Size returns the size of the stack.
func (s Stack[T]) Len() int {
	return s.container.Len()
}

// Push pushes a value into the stack.
func (s *Stack[T]) Push(value T) {
	s.container.PushBack(value)
}

// Pop pops a value from the queue.
func (s *Stack[T]) Pop() T {
	return s.container.Delete(s.container.Back())
}

// Top returns the value at the top of the queue.
func (s *Stack[T]) Top() T {
	return s.container.Back().Value
}
