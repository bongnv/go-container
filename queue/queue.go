// Package queue provides an implementation of the queue data structure in Go.
package queue

import (
	"github.com/bongnv/go-container/list"
)

// New creates a new queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{
		container: list.New[T](),
	}
}

// Queue is an implementation of queue.
type Queue[T any] struct {
	container *list.List[T]
}

// Size returns the size of the queue.
func (s Queue[T]) Len() int {
	return s.container.Len()
}

// Push pushes a value into the queue.
func (s *Queue[T]) Push(value T) {
	s.container.PushBack(value)
}

// Pop pops a value from the queue.
func (s *Queue[T]) Pop() T {
	return s.container.Delete(s.container.Front())
}

// Front returns the value at the front of the queue.
func (s *Queue[T]) Front() T {
	return s.container.Front().Value
}

// Back returns the value at the back of the queue.
func (s *Queue[T]) Back() T {
	return s.container.Back().Value
}
