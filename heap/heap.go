// Package heap provides an implementation of heap data structure in Go.
package heap

import (
	"cmp"
	"container/heap"

	"github.com/bongnv/go-container/algorithm"
)

// Element is an element in the heap.
type Element[T any] struct {
	Value T
	index int
}

// Heap represents a heap.
type Heap[T comparable] struct {
	container heapContainer[T]
}

// New creates a new heap of T.
func New[T cmp.Ordered]() *Heap[T] {
	return NewFunc[T](cmp.Less[T])
}

// NewFunc creates a new heap of T using less.
func NewFunc[T comparable](less algorithm.LessFunc[T]) *Heap[T] {
	return &Heap[T]{
		container: heapContainer[T]{
			less: less,
		},
	}
}

// Push pushes a value into the heap.
// It returns the created element for the provided value.
func (h *Heap[T]) Push(value T) *Element[T] {
	newNode := &Element[T]{
		Value: value,
	}
	heap.Push(&h.container, newNode)
	return newNode
}

// Pop pops a value from the heap.
func (h *Heap[T]) Pop() T {
	val := heap.Pop(&h.container).(*Element[T]).Value
	return val
}

// Top returns the element at the top of the heap.
func (h *Heap[T]) Top() *Element[T] {
	return h.container.nodes[0]
}

// Fix fixes the position of value in the heap data structure.
// It should be called after its data changes.
func (h *Heap[T]) Fix(e *Element[T]) {
	heap.Fix(&h.container, e.index)
}

// Size returns the size of the queue.
func (h *Heap[T]) Len() int {
	return len(h.container.nodes)
}

type heapContainer[T any] struct {
	nodes []*Element[T]
	less  algorithm.LessFunc[T]
}

func (hc heapContainer[T]) Len() int {
	return len(hc.nodes)
}

func (hc heapContainer[T]) Less(i, j int) bool {
	return hc.less(hc.nodes[i].Value, hc.nodes[j].Value)
}

func (hc heapContainer[T]) Swap(i, j int) {
	hc.nodes[i], hc.nodes[j] = hc.nodes[j], hc.nodes[i]
	hc.nodes[i].index = i
	hc.nodes[j].index = j
}

func (hc *heapContainer[T]) Push(x any) {
	n := len(hc.nodes)
	item := x.(*Element[T])
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
