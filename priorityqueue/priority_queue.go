// Package heap provides an implementation of priority queue structure in Go.
package priorityqueue

import (
	"container/heap"
)

// Less is a function that returns whether x < y or not.
type Less[T any] func(x, y T) bool

// PriorityQueue represents a priority queue.
type PriorityQueue[T comparable] struct {
	container heapContainer[T]
}

// New creates a new heap of T.
func New[T comparable](less Less[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		container: heapContainer[T]{
			less: less,
		},
	}
}

// Push pushes a value into the queue.
func (h *PriorityQueue[T]) Push(value T) {
	heap.Push(&h.container, value)
}

// Pop pops a value from the queue.
func (h *PriorityQueue[T]) Pop() T {
	val := heap.Pop(&h.container).(T)
	return val
}

// Top returns the value at the top of the queue.
func (h *PriorityQueue[T]) Top() T {
	return h.container.nodes[0]
}

// Size returns the size of the queue.
func (h *PriorityQueue[T]) Len() int {
	return len(h.container.nodes)
}

type heapContainer[T any] struct {
	nodes []T
	less  Less[T]
}

func (hc heapContainer[T]) Len() int {
	return len(hc.nodes)
}

func (hc heapContainer[T]) Less(i, j int) bool {
	return hc.less(hc.nodes[i], hc.nodes[j])
}

func (hc heapContainer[T]) Swap(i, j int) {
	hc.nodes[i], hc.nodes[j] = hc.nodes[j], hc.nodes[i]
}

func (hc *heapContainer[T]) Push(x any) {
	hc.nodes = append(hc.nodes, x.(T))
}

func (hc *heapContainer[T]) Pop() any {
	n := len(hc.nodes)
	item := hc.nodes[n-1]
	hc.nodes = hc.nodes[0 : n-1]
	return item
}
