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
	newNode := &heapNode[T]{
		value: value,
	}
	heap.Push(&h.container, newNode)
}

// Pop pops a value from the queue.
func (h *PriorityQueue[T]) Pop() T {
	val := heap.Pop(&h.container).(*heapNode[T]).value
	return val
}

// Top returns the value at the top of the queue.
func (h *PriorityQueue[T]) Top() T {
	return h.container.nodes[0].value
}

// Size returns the size of the queue.
func (h *PriorityQueue[T]) Len() int {
	return len(h.container.nodes)
}

type heapNode[T any] struct {
	value T
	index int
}

type heapContainer[T any] struct {
	nodes []*heapNode[T]
	less  Less[T]
}

func (hc heapContainer[T]) Len() int {
	return len(hc.nodes)
}

func (hc heapContainer[T]) Less(i, j int) bool {
	return hc.less(hc.nodes[i].value, hc.nodes[j].value)
}

func (hc heapContainer[T]) Swap(i, j int) {
	hc.nodes[i], hc.nodes[j] = hc.nodes[j], hc.nodes[i]
	hc.nodes[i].index = i
	hc.nodes[j].index = j
}

func (hc *heapContainer[T]) Push(x any) {
	n := len(hc.nodes)
	item := x.(*heapNode[T])
	item.index = n
	hc.nodes = append(hc.nodes, item)
}

func (hc *heapContainer[T]) Pop() any {
	n := len(hc.nodes)
	item := hc.nodes[n-1]
	hc.nodes[n-1] = nil // avoid memory leak
	item.index = -1     // for safety
	hc.nodes = hc.nodes[0 : n-1]
	return item
}
