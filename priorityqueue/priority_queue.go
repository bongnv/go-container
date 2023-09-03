// Package heap provides an implementation of priority queue structure in Go.
package priorityqueue

import (
	"cmp"
	"container/heap"

	"github.com/bongnv/go-container/algorithm"
)

// PriorityQueue represents a priority queue.
type PriorityQueue[T any] struct {
	container heapContainer[T]
}

// New creates a new priority queue of T.
func New[T cmp.Ordered]() *PriorityQueue[T] {
	return NewFunc[T](cmp.Less[T])
}

// NewNewFunc creates a new priority queue of T using Less function.
func NewFunc[T any](less algorithm.LessFunc[T]) *PriorityQueue[T] {
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
	less  algorithm.LessFunc[T]
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
