// Copyright 2020 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package btree

import (
	"cmp"
	"sync"
	"sync/atomic"
)

type BTree[T any] struct {
	isoid        uint64
	mu           *sync.RWMutex
	root         *node[T]
	count        int
	locks        bool
	copyItems    bool
	isoCopyItems bool
	less         func(a, b T) bool
	empty        T
	max          int
	min          int
}

type node[T any] struct {
	isoid    uint64
	count    int
	items    []T
	children *[]*node[T]
}

// PathHint is a utility type used with the *Hint() functions. Hints provide
// faster operations for clustered keys.
type PathHint struct {
	used [8]bool
	path [8]uint8
}

// Options for passing to New when creating a new BTree.
type Options struct {
	// Degree is used to define how many items and children each internal node
	// can contain before it must branch. For example, a degree of 2 will
	// create a 2-3-4 tree, where each node may contains 1-3 items and
	// 2-4 children. See https://en.wikipedia.org/wiki/2–3–4_tree.
	// Default is 32
	Degree int
	// NoLocks will disable locking. Otherwide a sync.RWMutex is used to
	// ensure all operations are safe across multiple goroutines.
	NoLocks bool
}

// New returns a new BTree
func NewBTree[T cmp.Ordered]() *BTree[T] {
	return NewBTreeOptions(cmp.Less[T], Options{Degree: 2})
}

func NewBTreeFunc[T any](less func(a, b T) bool) *BTree[T] {
	return NewBTreeOptions(less, Options{})
}

func NewBTreeOptions[T any](less func(a, b T) bool, opts Options) *BTree[T] {
	tr := new(BTree[T])
	tr.isoid = newIsoID()
	tr.mu = new(sync.RWMutex)
	tr.locks = !opts.NoLocks
	tr.less = less
	tr.init(opts.Degree)
	return tr
}

func (tr *BTree[T]) init(degree int) {
	if tr.min != 0 {
		return
	}
	tr.min, tr.max = degreeToMinMax(degree)
	_, tr.copyItems = ((interface{})(tr.empty)).(copier[T])
	if !tr.copyItems {
		_, tr.isoCopyItems = ((interface{})(tr.empty)).(isoCopier[T])
	}
}

func (tr *BTree[T]) newNode(leaf bool) *node[T] {
	n := &node[T]{isoid: tr.isoid}
	if !leaf {
		n.children = new([]*node[T])
	}
	return n
}

// leaf returns true if the node is a leaf.
func (n *node[T]) leaf() bool {
	return n.children == nil
}

func (tr *BTree[T]) bsearch(n *node[T], key T) (index int, found bool) {
	low, high := 0, len(n.items)
	for low < high {
		h := (low + high) / 2
		if !tr.less(key, n.items[h]) {
			low = h + 1
		} else {
			high = h
		}
	}
	if low > 0 && !tr.less(n.items[low-1], key) {
		return low - 1, true
	}
	return low, false
}

func (tr *BTree[T]) find(n *node[T], key T, hint *PathHint, depth int,
) (index int, found bool) {
	if hint == nil {
		return tr.bsearch(n, key)
	}
	return tr.hintsearch(n, key, hint, depth)
}

func (tr *BTree[T]) hintsearch(n *node[T], key T, hint *PathHint, depth int,
) (index int, found bool) {
	// Best case finds the exact match, updates the hint and returns.
	// Worst case, updates the low and high bounds to binary search between.
	low := 0
	high := len(n.items) - 1
	if depth < 8 && hint.used[depth] {
		index = int(hint.path[depth])
		if index >= len(n.items) {
			// tail item
			if tr.less(n.items[len(n.items)-1], key) {
				index = len(n.items)
				goto path_match
			}
			index = len(n.items) - 1
		}
		if tr.less(key, n.items[index]) {
			if index == 0 || tr.less(n.items[index-1], key) {
				goto path_match
			}
			high = index - 1
		} else if tr.less(n.items[index], key) {
			low = index + 1
		} else {
			found = true
			goto path_match
		}
	}

	// Do a binary search between low and high
	// keep on going until low > high, where the guarantee on low is that
	// key >= items[low - 1]
	for low <= high {
		mid := low + ((high+1)-low)/2
		// if key >= n.items[mid], low = mid + 1
		// which implies that key >= everything below low
		if !tr.less(key, n.items[mid]) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	// if low > 0, n.items[low - 1] >= key,
	// we have from before that key >= n.items[low - 1]
	// therefore key = n.items[low - 1],
	// and we have found the entry for key.
	// Otherwise we must keep searching for the key in index `low`.
	if low > 0 && !tr.less(n.items[low-1], key) {
		index = low - 1
		found = true
	} else {
		index = low
		found = false
	}

path_match:
	if depth < 8 {
		hint.used[depth] = true
		var pathIndex uint8
		if n.leaf() && found {
			pathIndex = uint8(index + 1)
		} else {
			pathIndex = uint8(index)
		}
		if pathIndex != hint.path[depth] {
			hint.path[depth] = pathIndex
			for i := depth + 1; i < 8; i++ {
				hint.used[i] = false
			}
		}
	}
	return index, found
}

// SetHint sets or replace a value for a key using a path hint
func (tr *BTree[T]) SetHint(item T, hint *PathHint) (prev T, replaced bool) {
	if tr.locks {
		tr.mu.Lock()
		prev, replaced = tr.setHint(item, hint)
		tr.mu.Unlock()
	} else {
		prev, replaced = tr.setHint(item, hint)
	}
	return prev, replaced
}

func (tr *BTree[T]) setHint(item T, hint *PathHint) (prev T, replaced bool) {
	if tr.root == nil {
		tr.init(0)
		tr.root = tr.newNode(true)
		tr.root.items = append([]T{}, item)
		tr.root.count = 1
		tr.count = 1
		return tr.empty, false
	}
	prev, replaced, split := tr.nodeSet(&tr.root, item, hint, 0)
	if split {
		left := tr.isoLoad(&tr.root, true)
		right, median := tr.nodeSplit(left)
		tr.root = tr.newNode(false)
		*tr.root.children = make([]*node[T], 0, tr.max+1)
		*tr.root.children = append([]*node[T]{}, left, right)
		tr.root.items = append([]T{}, median)
		tr.root.updateCount()
		return tr.setHint(item, hint)
	}
	if replaced {
		return prev, true
	}
	tr.count++
	return tr.empty, false
}

// Upsert replaces or inserts if the item doesn't exist.
// It returns the replaced item and whether it's replaced or not.
func (tr *BTree[T]) Upsert(item T) (T, bool) {
	return tr.SetHint(item, nil)
}

func (tr *BTree[T]) nodeSplit(n *node[T]) (right *node[T], median T) {
	i := tr.max / 2
	median = n.items[i]

	// right node
	right = tr.newNode(n.leaf())
	right.items = n.items[i+1:]
	if !n.leaf() {
		*right.children = (*n.children)[i+1:]
	}
	right.updateCount()

	// left node
	n.items[i] = tr.empty
	n.items = n.items[:i:i]
	if !n.leaf() {
		*n.children = (*n.children)[: i+1 : i+1]
	}
	n.updateCount()

	return right, median
}

func (n *node[T]) updateCount() {
	n.count = len(n.items)
	if !n.leaf() {
		for i := 0; i < len(*n.children); i++ {
			n.count += (*n.children)[i].count
		}
	}
}

// Copy the node for safe isolation.
func (tr *BTree[T]) copy(n *node[T]) *node[T] {
	n2 := new(node[T])
	n2.isoid = tr.isoid
	n2.count = n.count
	n2.items = make([]T, len(n.items), cap(n.items))
	copy(n2.items, n.items)
	if tr.copyItems {
		for i := 0; i < len(n2.items); i++ {
			n2.items[i] = ((interface{})(n2.items[i])).(copier[T]).Copy()
		}
	} else if tr.isoCopyItems {
		for i := 0; i < len(n2.items); i++ {
			n2.items[i] = ((interface{})(n2.items[i])).(isoCopier[T]).IsoCopy()
		}
	}
	if !n.leaf() {
		n2.children = new([]*node[T])
		*n2.children = make([]*node[T], len(*n.children), tr.max+1)
		copy(*n2.children, *n.children)
	}
	return n2
}

// isoLoad loads the provided node and, if needed, performs a copy-on-write.
func (tr *BTree[T]) isoLoad(cn **node[T], mut bool) *node[T] {
	if mut && (*cn).isoid != tr.isoid {
		*cn = tr.copy(*cn)
	}
	return *cn
}

func (tr *BTree[T]) nodeSet(cn **node[T], item T,
	hint *PathHint, depth int,
) (prev T, replaced, split bool) {
	if (*cn).isoid != tr.isoid {
		*cn = tr.copy(*cn)
	}
	n := *cn
	var i int
	var found bool
	if hint == nil {
		i, found = tr.bsearch(n, item)
	} else {
		i, found = tr.hintsearch(n, item, hint, depth)
	}
	if found {
		prev = n.items[i]
		n.items[i] = item
		return prev, true, false
	}
	if n.leaf() {
		if len(n.items) == tr.max {
			return tr.empty, false, true
		}
		n.items = append(n.items, tr.empty)
		copy(n.items[i+1:], n.items[i:])
		n.items[i] = item
		n.count++
		return tr.empty, false, false
	}
	prev, replaced, split = tr.nodeSet(&(*n.children)[i], item, hint, depth+1)
	if split {
		if len(n.items) == tr.max {
			return tr.empty, false, true
		}
		right, median := tr.nodeSplit((*n.children)[i])
		*n.children = append(*n.children, nil)
		copy((*n.children)[i+1:], (*n.children)[i:])
		(*n.children)[i+1] = right
		n.items = append(n.items, tr.empty)
		copy(n.items[i+1:], n.items[i:])
		n.items[i] = median
		return tr.nodeSet(&n, item, hint, depth)
	}
	if !replaced {
		n.count++
	}
	return prev, replaced, false
}

func (tr *BTree[T]) Scan(iter func(item T) bool) {
	tr.scan(iter, false)
}

func (tr *BTree[T]) ScanMut(iter func(item T) bool) {
	tr.scan(iter, true)
}

func (tr *BTree[T]) scan(iter func(item T) bool, mut bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil {
		return
	}
	tr.nodeScan(&tr.root, iter, mut)
}

func (tr *BTree[T]) nodeScan(cn **node[T], iter func(item T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	if n.leaf() {
		for i := 0; i < len(n.items); i++ {
			if !iter(n.items[i]) {
				return false
			}
		}
		return true
	}
	for i := 0; i < len(n.items); i++ {
		if !tr.nodeScan(&(*n.children)[i], iter, mut) {
			return false
		}
		if !iter(n.items[i]) {
			return false
		}
	}
	return tr.nodeScan(&(*n.children)[len(*n.children)-1], iter, mut)
}

// Get a value for key
func (tr *BTree[T]) Get(key T) (T, bool) {
	return tr.getHint(key, nil, false)
}

func (tr *BTree[T]) GetMut(key T) (T, bool) {
	return tr.getHint(key, nil, true)
}

// GetHint gets a value for key using a path hint
func (tr *BTree[T]) GetHint(key T, hint *PathHint) (value T, ok bool) {
	return tr.getHint(key, hint, false)
}

func (tr *BTree[T]) GetHintMut(key T, hint *PathHint) (value T, ok bool) {
	return tr.getHint(key, hint, true)
}

// GetHint gets a value for key using a path hint
func (tr *BTree[T]) getHint(key T, hint *PathHint, mut bool) (T, bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, mut)
	depth := 0
	for {
		i, found := tr.find(n, key, hint, depth)
		if found {
			return n.items[i], true
		}
		if n.children == nil {
			return tr.empty, false
		}
		n = tr.isoLoad(&(*n.children)[i], mut)
		depth++
	}
}

// Len returns the number of items in the tree
func (tr *BTree[T]) Len() int {
	return tr.count
}

// Delete a value for a key and returns the deleted value.
// Returns false if there was no value by that key found.
func (tr *BTree[T]) Delete(key T) (T, bool) {
	return tr.DeleteHint(key, nil)
}

// DeleteHint deletes a value for a key using a path hint and returns the
// deleted value.
// Returns false if there was no value by that key found.
func (tr *BTree[T]) DeleteHint(key T, hint *PathHint) (T, bool) {
	if tr.lock(true) {
		defer tr.unlock(true)
	}
	return tr.deleteHint(key, hint)
}

func (tr *BTree[T]) deleteHint(key T, hint *PathHint) (T, bool) {
	if tr.root == nil {
		return tr.empty, false
	}
	prev, deleted := tr.delete(&tr.root, false, key, hint, 0)
	if !deleted {
		return tr.empty, false
	}
	if len(tr.root.items) == 0 && !tr.root.leaf() {
		tr.root = (*tr.root.children)[0]
	}
	tr.count--
	if tr.count == 0 {
		tr.root = nil
	}
	return prev, true
}

func (tr *BTree[T]) delete(cn **node[T], max bool, key T,
	hint *PathHint, depth int,
) (T, bool) {
	n := tr.isoLoad(cn, true)
	var i int
	var found bool
	if max {
		i, found = len(n.items)-1, true
	} else {
		i, found = tr.find(n, key, hint, depth)
	}
	if n.leaf() {
		if found {
			// found the items at the leaf, remove it and return.
			prev := n.items[i]
			copy(n.items[i:], n.items[i+1:])
			n.items[len(n.items)-1] = tr.empty
			n.items = n.items[:len(n.items)-1]
			n.count--
			return prev, true
		}
		return tr.empty, false
	}

	var prev T
	var deleted bool
	if found {
		if max {
			i++
			prev, deleted = tr.delete(&(*n.children)[i], true, tr.empty, nil, 0)
		} else {
			prev = n.items[i]
			maxItem, _ := tr.delete(&(*n.children)[i], true, tr.empty, nil, 0)
			deleted = true
			n.items[i] = maxItem
		}
	} else {
		prev, deleted = tr.delete(&(*n.children)[i], max, key, hint, depth+1)
	}
	if !deleted {
		return tr.empty, false
	}
	n.count--
	if len((*n.children)[i].items) < tr.min {
		tr.nodeRebalance(n, i)
	}
	return prev, true
}

// nodeRebalance rebalances the child nodes following a delete operation.
// Provide the index of the child node with the number of items that fell
// below minItems.
func (tr *BTree[T]) nodeRebalance(n *node[T], i int) {
	if i == len(n.items) {
		i--
	}

	// ensure copy-on-write
	left := tr.isoLoad(&(*n.children)[i], true)
	right := tr.isoLoad(&(*n.children)[i+1], true)

	if len(left.items)+len(right.items) < tr.max {
		// Merges the left and right children nodes together as a single node
		// that includes (left,item,right), and places the contents into the
		// existing left node. Delete the right node altogether and move the
		// following items and child nodes to the left by one slot.

		// merge (left,item,right)
		left.items = append(left.items, n.items[i])
		left.items = append(left.items, right.items...)
		if !left.leaf() {
			*left.children = append(*left.children, *right.children...)
		}
		left.count += right.count + 1

		// move the items over one slot
		copy(n.items[i:], n.items[i+1:])
		n.items[len(n.items)-1] = tr.empty
		n.items = n.items[:len(n.items)-1]

		// move the children over one slot
		copy((*n.children)[i+1:], (*n.children)[i+2:])
		(*n.children)[len(*n.children)-1] = nil
		(*n.children) = (*n.children)[:len(*n.children)-1]
	} else if len(left.items) > len(right.items) {
		// move left -> right over one slot

		// Move the item of the parent node at index into the right-node first
		// slot, and move the left-node last item into the previously moved
		// parent item slot.
		right.items = append(right.items, tr.empty)
		copy(right.items[1:], right.items)
		right.items[0] = n.items[i]
		right.count++
		n.items[i] = left.items[len(left.items)-1]
		left.items[len(left.items)-1] = tr.empty
		left.items = left.items[:len(left.items)-1]
		left.count--

		if !left.leaf() {
			// move the left-node last child into the right-node first slot
			*right.children = append(*right.children, nil)
			copy((*right.children)[1:], *right.children)
			(*right.children)[0] = (*left.children)[len(*left.children)-1]
			(*left.children)[len(*left.children)-1] = nil
			(*left.children) = (*left.children)[:len(*left.children)-1]
			left.count -= (*right.children)[0].count
			right.count += (*right.children)[0].count
		}
	} else {
		// move left <- right over one slot

		// Same as above but the other direction
		left.items = append(left.items, n.items[i])
		left.count++
		n.items[i] = right.items[0]
		copy(right.items, right.items[1:])
		right.items[len(right.items)-1] = tr.empty
		right.items = right.items[:len(right.items)-1]
		right.count--

		if !left.leaf() {
			*left.children = append(*left.children, (*right.children)[0])
			copy(*right.children, (*right.children)[1:])
			(*right.children)[len(*right.children)-1] = nil
			*right.children = (*right.children)[:len(*right.children)-1]
			left.count += (*left.children)[len(*left.children)-1].count
			right.count -= (*left.children)[len(*left.children)-1].count
		}
	}
}

// Ascend the tree within the range [pivot, last]
// Pass nil for pivot to scan all item in ascending order
// Return false to stop iterating
func (tr *BTree[T]) Ascend(pivot T, iter func(item T) bool) {
	tr.ascend(pivot, iter, false)
}

func (tr *BTree[T]) AscendMut(pivot T, iter func(item T) bool) {
	tr.ascend(pivot, iter, true)
}

func (tr *BTree[T]) ascend(pivot T, iter func(item T) bool, mut bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil {
		return
	}
	tr.nodeAscend(&tr.root, pivot, nil, 0, iter, mut)
}

// The return value of this function determines whether we should keep iterating
// upon this functions return.
func (tr *BTree[T]) nodeAscend(cn **node[T], pivot T, hint *PathHint,
	depth int, iter func(item T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	i, found := tr.find(n, pivot, hint, depth)
	if !found {
		if !n.leaf() {
			if !tr.nodeAscend(&(*n.children)[i], pivot, hint, depth+1, iter,
				mut) {
				return false
			}
		}
	}
	// We are either in the case that
	// - node is found, we should iterate through it starting at `i`,
	//   the index it was located at.
	// - node is not found, and TODO: fill in.
	for ; i < len(n.items); i++ {
		if !iter(n.items[i]) {
			return false
		}
		if !n.leaf() {
			if !tr.nodeScan(&(*n.children)[i+1], iter, mut) {
				return false
			}
		}
	}
	return true
}

func (tr *BTree[T]) ReverseScan(iter func(item T) bool) {
	tr.reverse(iter, false)
}

func (tr *BTree[T]) ReverseScanMut(iter func(item T) bool) {
	tr.reverse(iter, true)
}

func (tr *BTree[T]) reverse(iter func(item T) bool, mut bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil {
		return
	}
	tr.nodeReverse(&tr.root, iter, mut)
}

func (tr *BTree[T]) nodeReverse(cn **node[T], iter func(item T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	if n.leaf() {
		for i := len(n.items) - 1; i >= 0; i-- {
			if !iter(n.items[i]) {
				return false
			}
		}
		return true
	}
	if !tr.nodeReverse(&(*n.children)[len(*n.children)-1], iter, mut) {
		return false
	}
	for i := len(n.items) - 1; i >= 0; i-- {
		if !iter(n.items[i]) {
			return false
		}
		if !tr.nodeReverse(&(*n.children)[i], iter, mut) {
			return false
		}
	}
	return true
}

// Descend the tree within the range [pivot, first]
// Pass nil for pivot to scan all item in descending order
// Return false to stop iterating
func (tr *BTree[T]) Descend(pivot T, iter func(item T) bool) {
	tr.descend(pivot, iter, false)
}

func (tr *BTree[T]) DescendMut(pivot T, iter func(item T) bool) {
	tr.descend(pivot, iter, true)
}

func (tr *BTree[T]) descend(pivot T, iter func(item T) bool, mut bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil {
		return
	}
	tr.nodeDescend(&tr.root, pivot, nil, 0, iter, mut)
}

func (tr *BTree[T]) nodeDescend(cn **node[T], pivot T, hint *PathHint,
	depth int, iter func(item T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	i, found := tr.find(n, pivot, hint, depth)
	if !found {
		if !n.leaf() {
			if !tr.nodeDescend(&(*n.children)[i], pivot, hint, depth+1, iter,
				mut) {
				return false
			}
		}
		i--
	}
	for ; i >= 0; i-- {
		if !iter(n.items[i]) {
			return false
		}
		if !n.leaf() {
			if !tr.nodeReverse(&(*n.children)[i], iter, mut) {
				return false
			}
		}
	}
	return true
}

// Load is for bulk loading pre-sorted items
func (tr *BTree[T]) Load(item T) (T, bool) {
	if tr.lock(true) {
		defer tr.unlock(true)
	}
	if tr.root == nil {
		return tr.setHint(item, nil)
	}
	n := tr.isoLoad(&tr.root, true)
	for {
		n.count++ // optimistically update counts
		if n.leaf() {
			if len(n.items) < tr.max {
				if tr.less(n.items[len(n.items)-1], item) {
					n.items = append(n.items, item)
					tr.count++
					return tr.empty, false
				}
			}
			break
		}
		n = tr.isoLoad(&(*n.children)[len(*n.children)-1], true)
	}
	// revert the counts
	n = tr.root
	for {
		n.count--
		if n.leaf() {
			break
		}
		n = (*n.children)[len(*n.children)-1]
	}
	return tr.setHint(item, nil)
}

// Min returns the minimum item in tree.
// Returns nil if the treex has no items.
func (tr *BTree[T]) Min() (T, bool) {
	return tr.minMut(false)
}

func (tr *BTree[T]) MinMut() (T, bool) {
	return tr.minMut(true)
}

func (tr *BTree[T]) minMut(mut bool) (T, bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, mut)
	for {
		if n.leaf() {
			return n.items[0], true
		}
		n = tr.isoLoad(&(*n.children)[0], mut)
	}
}

// Max returns the maximum item in tree.
// Returns nil if the tree has no items.
func (tr *BTree[T]) Max() (T, bool) {
	return tr.maxMut(false)
}

func (tr *BTree[T]) MaxMut() (T, bool) {
	return tr.maxMut(true)
}

func (tr *BTree[T]) maxMut(mut bool) (T, bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, mut)
	for {
		if n.leaf() {
			return n.items[len(n.items)-1], true
		}
		n = tr.isoLoad(&(*n.children)[len(*n.children)-1], mut)
	}
}

// DeleteMin removes the minimum item in tree and returns it.
// Returns nil if the tree has no items.
func (tr *BTree[T]) DeleteMin() (T, bool) {
	if tr.lock(true) {
		defer tr.unlock(true)
	}
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, true)
	var item T
	for {
		n.count-- // optimistically update counts
		if n.leaf() {
			item = n.items[0]
			if len(n.items) == tr.min {
				break
			}
			copy(n.items[:], n.items[1:])
			n.items[len(n.items)-1] = tr.empty
			n.items = n.items[:len(n.items)-1]
			tr.count--
			if tr.count == 0 {
				tr.root = nil
			}
			return item, true
		}
		n = tr.isoLoad(&(*n.children)[0], true)
	}
	// revert the counts
	n = tr.root
	for {
		n.count++
		if n.leaf() {
			break
		}
		n = (*n.children)[0]
	}
	return tr.deleteHint(item, nil)
}

// DeleteMax removes the maximum item in tree and returns it.
// Returns nil if the tree has no items.
func (tr *BTree[T]) DeleteMax() (T, bool) {
	if tr.lock(true) {
		defer tr.unlock(true)
	}
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, true)
	var item T
	for {
		n.count-- // optimistically update counts
		if n.leaf() {
			item = n.items[len(n.items)-1]
			if len(n.items) == tr.min {
				break
			}
			n.items[len(n.items)-1] = tr.empty
			n.items = n.items[:len(n.items)-1]
			tr.count--
			if tr.count == 0 {
				tr.root = nil
			}
			return item, true
		}
		n = tr.isoLoad(&(*n.children)[len(*n.children)-1], true)
	}
	// revert the counts
	n = tr.root
	for {
		n.count++
		if n.leaf() {
			break
		}
		n = (*n.children)[len(*n.children)-1]
	}
	return tr.deleteHint(item, nil)
}

// GetAt returns the value at index.
// Return nil if the tree is empty or the index is out of bounds.
func (tr *BTree[T]) GetAt(index int) (T, bool) {
	return tr.getAt(index, false)
}

func (tr *BTree[T]) GetAtMut(index int) (T, bool) {
	return tr.getAt(index, true)
}

func (tr *BTree[T]) getAt(index int, mut bool) (T, bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil || index < 0 || index >= tr.count {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, mut)
	for {
		if n.leaf() {
			return n.items[index], true
		}
		i := 0
		for ; i < len(n.items); i++ {
			if index < (*n.children)[i].count {
				break
			} else if index == (*n.children)[i].count {
				return n.items[i], true
			}
			index -= (*n.children)[i].count + 1
		}
		n = tr.isoLoad(&(*n.children)[i], mut)
	}
}

// DeleteAt deletes the item at index.
// Return nil if the tree is empty or the index is out of bounds.
func (tr *BTree[T]) DeleteAt(index int) (T, bool) {
	if tr.lock(true) {
		defer tr.unlock(true)
	}
	if tr.root == nil || index < 0 || index >= tr.count {
		return tr.empty, false
	}
	var pathbuf [8]uint8 // track the path
	path := pathbuf[:0]
	var item T
	n := tr.isoLoad(&tr.root, true)
outer:
	for {
		n.count-- // optimistically update counts
		if n.leaf() {
			// the index is the item position
			item = n.items[index]
			if len(n.items) == tr.min {
				path = append(path, uint8(index))
				break outer
			}
			copy(n.items[index:], n.items[index+1:])
			n.items[len(n.items)-1] = tr.empty
			n.items = n.items[:len(n.items)-1]
			tr.count--
			if tr.count == 0 {
				tr.root = nil
			}
			return item, true
		}
		i := 0
		for ; i < len(n.items); i++ {
			if index < (*n.children)[i].count {
				break
			} else if index == (*n.children)[i].count {
				item = n.items[i]
				path = append(path, uint8(i))
				break outer
			}
			index -= (*n.children)[i].count + 1
		}
		path = append(path, uint8(i))
		n = tr.isoLoad(&(*n.children)[i], true)
	}
	// revert the counts
	var hint PathHint
	n = tr.root
	for i := 0; i < len(path); i++ {
		if i < len(hint.path) {
			hint.path[i] = uint8(path[i])
			hint.used[i] = true
		}
		n.count++
		if !n.leaf() {
			n = (*n.children)[uint8(path[i])]
		}
	}
	return tr.deleteHint(item, &hint)
}

// Height returns the height of the tree.
// Returns zero if tree has no items.
func (tr *BTree[T]) Height() int {
	if tr.lock(false) {
		defer tr.unlock(false)
	}
	var height int
	if tr.root != nil {
		n := tr.root
		for {
			height++
			if n.leaf() {
				break
			}
			n = (*n.children)[0]
		}
	}
	return height
}

// Walk iterates over all items in tree, in order.
// The items param will contain one or more items.
func (tr *BTree[T]) Walk(iter func(item []T) bool) {
	tr.walk(iter, false)
}

func (tr *BTree[T]) WalkMut(iter func(item []T) bool) {
	tr.walk(iter, true)
}

func (tr *BTree[T]) walk(iter func(item []T) bool, mut bool) {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	if tr.root == nil {
		return
	}
	tr.nodeWalk(&tr.root, iter, mut)
}

func (tr *BTree[T]) nodeWalk(cn **node[T], iter func(item []T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	if n.leaf() {
		if !iter(n.items) {
			return false
		}
	} else {
		for i := 0; i < len(n.items); i++ {
			if !tr.nodeWalk(&(*n.children)[i], iter, mut) {
				return false
			}
			if !iter(n.items[i : i+1]) {
				return false
			}
		}
		if !tr.nodeWalk(&(*n.children)[len(n.items)], iter, mut) {
			return false
		}
	}
	return true
}

// Copy the tree. This is a copy-on-write operation and is very fast because
// it only performs a shadowed copy.
func (tr *BTree[T]) Copy() *BTree[T] {
	return tr.IsoCopy()
}

func (tr *BTree[T]) IsoCopy() *BTree[T] {
	if tr.lock(true) {
		defer tr.unlock(true)
	}
	tr.isoid = newIsoID()
	tr2 := new(BTree[T])
	*tr2 = *tr
	tr2.mu = new(sync.RWMutex)
	tr2.isoid = newIsoID()
	return tr2
}

func (tr *BTree[T]) lock(write bool) bool {
	if tr.locks {
		if write {
			tr.mu.Lock()
		} else {
			tr.mu.RLock()
		}
	}
	return tr.locks
}

func (tr *BTree[T]) unlock(write bool) {
	if write {
		tr.mu.Unlock()
	} else {
		tr.mu.RUnlock()
	}
}

// Values returns all the items in order.
func (tr *BTree[T]) Values() []T {
	return tr.items(false)
}

func (tr *BTree[T]) ItemsMut() []T {
	return tr.items(true)
}

func (tr *BTree[T]) items(mut bool) []T {
	if tr.lock(mut) {
		defer tr.unlock(mut)
	}
	items := make([]T, 0, tr.Len())
	if tr.root != nil {
		items = tr.nodeItems(&tr.root, items, mut)
	}
	return items
}

func (tr *BTree[T]) nodeItems(cn **node[T], items []T, mut bool) []T {
	n := tr.isoLoad(cn, mut)
	if n.leaf() {
		return append(items, n.items...)
	}
	for i := 0; i < len(n.items); i++ {
		items = tr.nodeItems(&(*n.children)[i], items, mut)
		items = append(items, n.items[i])
	}
	return tr.nodeItems(&(*n.children)[len(*n.children)-1], items, mut)
}

// Clear will delete all items.
func (tr *BTree[T]) Clear() {
	if tr.lock(true) {
		defer tr.unlock(true)
	}
	tr.root = nil
	tr.count = 0
}

var gisoid uint64

func newIsoID() uint64 {
	return atomic.AddUint64(&gisoid, 1)
}

type copier[T any] interface {
	Copy() T
}

type isoCopier[T any] interface {
	IsoCopy() T
}

func degreeToMinMax(deg int) (min, max int) {
	if deg <= 0 {
		deg = 32
	} else if deg == 1 {
		deg = 2 // must have at least 2
	}
	max = deg*2 - 1 // max items per node. max children is +1
	min = max / 2
	return min, max
}
