package rbtree

import (
	"cmp"

	"github.com/bongnv/go-container/algorithm"
)

// ItemIterator is a function to iterate through items.
type ItemIterator[T any] func(i T) bool

// Tree is a Left-Leaning Red-Black (LLRB) implementation of 2-3 trees
type LLRB[T any] struct {
	count int
	root  *Node[T]
	less  algorithm.LessFunc[T]
}

// Node represents a node in LLRB.
type Node[T any] struct {
	Item        T
	Left, Right *Node[T] // Pointers to left and right child nodes
	Black       bool     // If set, the color of the link (incoming from the parent) is black
	// In the LLRB, new nodes are always red, hence the zero-value for node
}

// New allocates a new tree
func New[T cmp.Ordered]() *LLRB[T] {
	return NewFunc[T](cmp.Less[T])
}

// NewFunc creates a new LLRB tree using less.
func NewFunc[T any](less algorithm.LessFunc[T]) *LLRB[T] {
	return &LLRB[T]{
		less: less,
	}
}

// SetRoot sets the root node of the tree.
// It is intended to be used by functions that deserialize the tree.
func (t *LLRB[T]) SetRoot(r *Node[T]) {
	t.root = r
}

// Root returns the root node of the tree.
// It is intended to be used by functions that serialize the tree.
func (t *LLRB[T]) Root() *Node[T] {
	return t.root
}

// Len returns the number of nodes in the tree.
func (t *LLRB[T]) Len() int { return t.count }

// Has returns true if the tree contains an element whose order is the same as that of key.
func (t *LLRB[T]) Has(key T) bool {
	_, found := t.Get(key)
	return found
}

// Get retrieves an element from the tree whose order is the same as that of key.
func (t *LLRB[T]) Get(key T) (item T, present bool) {
	h := t.root
	for h != nil {
		switch {
		case t.less(key, h.Item):
			h = h.Left
		case t.less(h.Item, key):
			h = h.Right
		default:
			return h.Item, true
		}
	}
	return
}

// Min returns the minimum element in the tree.
func (t *LLRB[T]) Min() (item T, present bool) {
	h := t.root
	if h == nil {
		return
	}
	for h.Left != nil {
		h = h.Left
	}
	return h.Item, true
}

// Max returns the maximum element in the tree.
func (t *LLRB[T]) Max() (item T, present bool) {
	h := t.root
	if h == nil {
		return
	}
	for h.Right != nil {
		h = h.Right
	}
	return h.Item, true
}

// Upsert inserts item into the tree. If an existing
// element has the same order, it is removed from the tree and returned.
func (t *LLRB[T]) Upsert(item T) (replacedItem T, replaced bool) {
	t.root, replacedItem, replaced = t.replaceOrInsert(t.root, item)
	t.root.Black = true
	if !replaced {
		t.count++
	}
	return replacedItem, replaced
}

func (t *LLRB[T]) replaceOrInsert(h *Node[T], item T) (node *Node[T], replacedTtem T, replaced bool) {
	if h == nil {
		node = newNode[T](item)
		return
	}

	h = walkDownRot23(h)

	if t.less(item, h.Item) { // BUG
		h.Left, replacedTtem, replaced = t.replaceOrInsert(h.Left, item)
	} else if t.less(h.Item, item) {
		h.Right, replacedTtem, replaced = t.replaceOrInsert(h.Right, item)
	} else {
		replacedTtem, h.Item, replaced = h.Item, item, true
	}

	h = walkUpRot23(h)

	return h, replacedTtem, replaced
}

// Insert inserts item into the tree. If an existing
// element has the same order, both elements remain in the tree.
func (t *LLRB[T]) Insert(item T) {
	t.root = t.insertNoReplace(t.root, item)
	t.root.Black = true
	t.count++
}

func (t *LLRB[T]) insertNoReplace(h *Node[T], item T) *Node[T] {
	if h == nil {
		return newNode(item)
	}

	h = walkDownRot23(h)

	if t.less(item, h.Item) {
		h.Left = t.insertNoReplace(h.Left, item)
	} else {
		h.Right = t.insertNoReplace(h.Right, item)
	}

	return walkUpRot23(h)
}

// Rotation driver routines for 2-3 algorithm

func walkDownRot23[T any](h *Node[T]) *Node[T] { return h }

func walkUpRot23[T any](h *Node[T]) *Node[T] {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

// Rotation driver routines for 2-3-4 algorithm

func walkDownRot234[T any](h *Node[T]) *Node[T] {
	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

func walkUpRot234[T any](h *Node[T]) *Node[T] {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	return h
}

// DeleteMin deletes the minimum element in the tree and returns the
// deleted item or nil otherwise.
func (t *LLRB[T]) DeleteMin() (deletedItem T, deleted bool) {
	t.root, deletedItem, deleted = deleteMin(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted {
		t.count--
	}
	return deletedItem, deleted
}

// deleteMin code for LLRB 2-3 trees
func deleteMin[T any](h *Node[T]) (node *Node[T], deletedItem T, deleted bool) {
	if h == nil {
		return nil, deletedItem, false
	}
	if h.Left == nil {
		return nil, h.Item, true
	}

	if !isRed(h.Left) && !isRed(h.Left.Left) {
		h = moveRedLeft(h)
	}

	h.Left, deletedItem, deleted = deleteMin(h.Left)

	return fixUp(h), deletedItem, deleted
}

// DeleteMax deletes the maximum element in the tree and returns
// the deleted item or nil otherwise
func (t *LLRB[T]) DeleteMax() (deletedItem T, deleted bool) {
	t.root, deletedItem, deleted = deleteMax(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted {
		t.count--
	}
	return deletedItem, deleted
}

func deleteMax[T any](h *Node[T]) (node *Node[T], deletedItem T, deleted bool) {
	if h == nil {
		return nil, deletedItem, false
	}
	if isRed(h.Left) {
		h = rotateRight(h)
	}
	if h.Right == nil {
		return nil, h.Item, true
	}
	if !isRed(h.Right) && !isRed(h.Right.Left) {
		h = moveRedRight(h)
	}
	h.Right, deletedItem, deleted = deleteMax(h.Right)

	return fixUp(h), deletedItem, deleted
}

// Delete deletes an item from the tree whose key equals key.
// The deleted item is return, otherwise nil is returned.
func (t *LLRB[T]) Delete(key T) (deletedItem T, deleted bool) {
	t.root, deletedItem, deleted = t.delete(t.root, key)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted {
		t.count--
	}
	return deletedItem, deleted
}

func (t *LLRB[T]) delete(h *Node[T], item T) (node *Node[T], deletedItem T, deleted bool) {
	if h == nil {
		return nil, deletedItem, false
	}
	if t.less(item, h.Item) {
		if h.Left == nil { // item not present. Nothing to delete
			return h, deletedItem, false
		}
		if !isRed(h.Left) && !isRed(h.Left.Left) {
			h = moveRedLeft(h)
		}
		h.Left, deletedItem, deleted = t.delete(h.Left, item)
	} else {
		if isRed(h.Left) {
			h = rotateRight(h)
		}
		// If @item equals @h.Item and no right children at @h
		if !t.less(h.Item, item) && h.Right == nil {
			return nil, h.Item, true
		}
		// PETAR: Added 'h.Right != nil' below
		if h.Right != nil && !isRed(h.Right) && !isRed(h.Right.Left) {
			h = moveRedRight(h)
		}
		// If @item equals @h.Item, and (from above) 'h.Right != nil'
		if !t.less(h.Item, item) {
			var subDeleted T
			h.Right, subDeleted, deleted = deleteMin(h.Right)
			deletedItem, h.Item = h.Item, subDeleted
		} else { // Else, @item is bigger than @h.Item
			h.Right, deletedItem, deleted = t.delete(h.Right, item)
		}
	}

	return fixUp(h), deletedItem, deleted
}

// Internal node manipulation routines

func newNode[T any](item T) *Node[T] { return &Node[T]{Item: item} }

func isRed[T any](h *Node[T]) bool {
	if h == nil {
		return false
	}
	return !h.Black
}

func rotateLeft[T any](h *Node[T]) *Node[T] {
	x := h.Right
	if x.Black {
		panic("rotating a black link")
	}
	h.Right = x.Left
	x.Left = h
	x.Black = h.Black
	h.Black = false
	return x
}

func rotateRight[T any](h *Node[T]) *Node[T] {
	x := h.Left
	if x.Black {
		panic("rotating a black link")
	}
	h.Left = x.Right
	x.Right = h
	x.Black = h.Black
	h.Black = false
	return x
}

// REQUIRE: Left and Right children must be present
func flip[T any](h *Node[T]) {
	h.Black = !h.Black
	h.Left.Black = !h.Left.Black
	h.Right.Black = !h.Right.Black
}

// REQUIRE: Left and Right children must be present
func moveRedLeft[T any](h *Node[T]) *Node[T] {
	flip(h)
	if isRed(h.Right.Left) {
		h.Right = rotateRight(h.Right)
		h = rotateLeft(h)
		flip(h)
	}
	return h
}

// REQUIRE: Left and Right children must be present
func moveRedRight[T any](h *Node[T]) *Node[T] {
	flip(h)
	if isRed(h.Left.Left) {
		h = rotateRight(h)
		flip(h)
	}
	return h
}

func fixUp[T any](h *Node[T]) *Node[T] {
	if isRed(h.Right) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

func (t *LLRB[T]) AscendRange(greaterOrEqual, lessThan T, iterator ItemIterator[T]) {
	t.ascendRange(t.root, greaterOrEqual, lessThan, iterator)
}

func (t *LLRB[T]) ascendRange(h *Node[T], inf, sup T, iterator ItemIterator[T]) bool {
	if h == nil {
		return true
	}
	if !t.less(h.Item, sup) {
		return t.ascendRange(h.Left, inf, sup, iterator)
	}
	if t.less(h.Item, inf) {
		return t.ascendRange(h.Right, inf, sup, iterator)
	}

	if !t.ascendRange(h.Left, inf, sup, iterator) {
		return false
	}
	if !iterator(h.Item) {
		return false
	}
	return t.ascendRange(h.Right, inf, sup, iterator)
}

// AscendGreaterOrEqual will call iterator once for each element greater or equal to
// pivot in ascending order. It will stop whenever the iterator returns false.
func (t *LLRB[T]) AscendGreaterOrEqual(pivot T, iterator ItemIterator[T]) {
	t.ascendGreaterOrEqual(t.root, pivot, iterator)
}

func (t *LLRB[T]) ascendGreaterOrEqual(h *Node[T], pivot T, iterator ItemIterator[T]) bool {
	if h == nil {
		return true
	}
	if !t.less(h.Item, pivot) {
		if !t.ascendGreaterOrEqual(h.Left, pivot, iterator) {
			return false
		}
		if !iterator(h.Item) {
			return false
		}
	}
	return t.ascendGreaterOrEqual(h.Right, pivot, iterator)
}

// AscendLessThan will call iterator once for each element lower than
// pivot in ascending order. It will stop whenever the iterator returns false.
func (t *LLRB[T]) AscendLessThan(pivot T, iterator ItemIterator[T]) {
	t.ascendLessThan(t.root, pivot, iterator)
}

func (t *LLRB[T]) ascendLessThan(h *Node[T], pivot T, iterator ItemIterator[T]) bool {
	if h == nil {
		return true
	}
	if !t.ascendLessThan(h.Left, pivot, iterator) {
		return false
	}
	if t.less(h.Item, pivot) {
		if !iterator(h.Item) {
			return false
		}
		return t.ascendLessThan(h.Right, pivot, iterator)
	}
	return true
}

// DescendLessOrEqual will call iterator once for each element less than the
// pivot in descending order. It will stop whenever the iterator returns false.
func (t *LLRB[T]) DescendLessOrEqual(pivot T, iterator ItemIterator[T]) {
	t.descendLessOrEqual(t.root, pivot, iterator)
}

func (t *LLRB[T]) descendLessOrEqual(h *Node[T], pivot T, iterator ItemIterator[T]) bool {
	if h == nil {
		return true
	}
	if t.less(h.Item, pivot) || !t.less(pivot, h.Item) {
		if !t.descendLessOrEqual(h.Right, pivot, iterator) {
			return false
		}
		if !iterator(h.Item) {
			return false
		}
	}
	return t.descendLessOrEqual(h.Left, pivot, iterator)
}

// Scan will call iterator once for each element in ascending order.
// It will stop whenever the iterator returns false.
func (t *LLRB[T]) Scan(iterator ItemIterator[T]) {
	t.ascend(t.root, iterator)
}

func (t *LLRB[T]) ascend(h *Node[T], iterator ItemIterator[T]) bool {
	if h == nil {
		return true
	}
	if !t.ascend(h.Left, iterator) {
		return false
	}
	if !iterator(h.Item) {
		return false
	}
	return t.ascend(h.Right, iterator)
}

// ReverseScan will call iterator once for each element in descending order.
// It will stop whenever the iterator returns false.
func (t *LLRB[T]) ReverseScan(iterator ItemIterator[T]) {
	t.descend(t.root, iterator)
}

func (t *LLRB[T]) descend(h *Node[T], iterator ItemIterator[T]) bool {
	if h == nil {
		return true
	}
	if !t.descend(h.Right, iterator) {
		return false
	}
	if !iterator(h.Item) {
		return false
	}
	return t.descend(h.Left, iterator)
}

// Values returns all values from the tree in order.
func (t *LLRB[T]) Values() []T {
	allValues := make([]T, 0, t.Len())
	t.ascend(t.root, func(value T) bool {
		allValues = append(allValues, value)
		return true
	})
	return allValues
}
