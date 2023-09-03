// Package heap provides an implementation of heap data structure in Go.
package heap

import (
	"cmp"
	"container/heap"

	"github.com/bongnv/go-container/algorithm"
)

// Heap represents a heap.
type Heap[T comparable] struct {
	container  heapContainer[T]
	nodeLookUp map[T]*heapNode[T]
}

// New creates a new heap of T.
func New[T cmp.Ordered]() *Heap[T] {
	return NewFunc[T](cmp.Less[T])
}

// NewFunc creates a new heap of T using less.
func NewFunc[T comparable](less algorithm.LessFunc[T]) *Heap[T] {
	return &Heap[T]{
		nodeLookUp: map[T]*heapNode[T]{},
		container: heapContainer[T]{
			less: less,
		},
	}
}

// Push pushes a value into the heap.
func (h *Heap[T]) Push(value T) {
	newNode := &heapNode[T]{
		value: value,
	}
	h.nodeLookUp[value] = newNode
	heap.Push(&h.container, newNode)
}

// Pop pops a value from the heap.
func (h *Heap[T]) Pop() T {
	val := heap.Pop(&h.container).(*heapNode[T]).value
	delete(h.nodeLookUp, val)
	return val
}

// Top returns the value at the top of the heap.
func (h *Heap[T]) Top() T {
	return h.container.nodes[0].value
}

// Fix fixes the position of value in the heap data structure.
// It should be called after its data changes.
func (h *Heap[T]) Fix(value T) {
	heap.Fix(&h.container, h.nodeLookUp[value].index)
}

// Size returns the size of the queue.
func (h *Heap[T]) Len() int {
	return len(h.container.nodes)
}

type heapNode[T any] struct {
	value T
	index int
}

type heapContainer[T any] struct {
	nodes []*heapNode[T]
	less  algorithm.LessFunc[T]
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
