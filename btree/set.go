package btree

import "cmp"

// NewSet creates a new set with degree = 2.
func NewSet[T cmp.Ordered]() *Set[T] {
	return NewSetDegree[T](2)
}

func NewSetDegree[T cmp.Ordered](degree int) *Set[T] {
	s := &Set[T]{}
	s.base.init(degree)
	return s
}

type Set[K cmp.Ordered] struct {
	base Map[K, struct{}]
}

// Copy
func (tr *Set[K]) Copy() *Set[K] {
	tr2 := new(Set[K])
	tr2.base = *tr.base.Copy()
	return tr2
}

func (tr *Set[K]) IsoCopy() *Set[K] {
	tr2 := new(Set[K])
	tr2.base = *tr.base.IsoCopy()
	return tr2
}

// Insert an item
func (tr *Set[K]) Insert(key K) {
	tr.base.Set(key, struct{}{})
}

func (tr *Set[K]) Scan(iter func(key K) bool) {
	tr.base.Scan(func(key K, value struct{}) bool {
		return iter(key)
	})
}

// Has checks whether a key exists or not.
func (tr *Set[K]) Has(key K) bool {
	_, ok := tr.base.Get(key)
	return ok
}

// Len returns the number of items in the tree
func (tr *Set[K]) Len() int {
	return tr.base.Len()
}

// Delete an item
func (tr *Set[K]) Delete(key K) {
	tr.base.Delete(key)
}

// Ascend the tree within the range [pivot, last]
// Pass nil for pivot to scan all item in ascending order
// Return false to stop iterating
func (tr *Set[K]) Ascend(pivot K, iter func(key K) bool) {
	tr.base.Ascend(pivot, func(key K, value struct{}) bool {
		return iter(key)
	})
}

func (tr *Set[K]) Reverse(iter func(key K) bool) {
	tr.base.Reverse(func(key K, value struct{}) bool {
		return iter(key)
	})
}

// Descend the tree within the range [pivot, first]
// Pass nil for pivot to scan all item in descending order
// Return false to stop iterating
func (tr *Set[K]) Descend(pivot K, iter func(key K) bool) {
	tr.base.Descend(pivot, func(key K, value struct{}) bool {
		return iter(key)
	})
}

// Load is for bulk loading pre-sorted items
func (tr *Set[K]) Load(key K) {
	tr.base.Load(key, struct{}{})
}

// Min returns the minimum item in tree.
// Returns nil if the treex has no items.
func (tr *Set[K]) Min() (K, bool) {
	key, _, ok := tr.base.Min()
	return key, ok
}

// Max returns the maximum item in tree.
// Returns nil if the tree has no items.
func (tr *Set[K]) Max() (K, bool) {
	key, _, ok := tr.base.Max()
	return key, ok
}

// DeleteMin removes the minimum item in tree and returns it.
// Returns nil if the tree has no items.
func (tr *Set[K]) DeleteMin() (K, bool) {
	key, _, ok := tr.base.PopMin()
	return key, ok
}

// PopMax removes the maximum item in tree and returns it.
// Returns nil if the tree has no items.
func (tr *Set[K]) DeleteMax() (K, bool) {
	key, _, ok := tr.base.PopMax()
	return key, ok
}

// GetAt returns the value at index.
// Return nil if the tree is empty or the index is out of bounds.
func (tr *Set[K]) GetAt(index int) (K, bool) {
	key, _, ok := tr.base.GetAt(index)
	return key, ok
}

// DeleteAt deletes the item at index.
// Return nil if the tree is empty or the index is out of bounds.
func (tr *Set[K]) DeleteAt(index int) (K, bool) {
	key, _, ok := tr.base.DeleteAt(index)
	return key, ok
}

// Height returns the height of the tree.
// Returns zero if tree has no items.
func (tr *Set[K]) Height() int {
	return tr.base.Height()
}

// Keys returns all the keys in order.
func (tr *Set[K]) Keys() []K {
	return tr.base.Keys()
}

// Clear will delete all items.
func (tr *Set[K]) Clear() {
	tr.base.Clear()
}
