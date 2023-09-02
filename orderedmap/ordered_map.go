// Package orderedmap provides an implementation of ordered map structure in Go.
// In this map, keys are maintained in an order.
package orderedmap

import (
	"cmp"
	"errors"

	"github.com/bongnv/go-container/list"
)

// ErrKeyNotFound means the key couldn't be found in the map.
var ErrKeyNotFound = errors.New("orderedmap: key not found")

// Pair is a pair of key and value.
type Pair[K, V any] struct {
	Key   K
	Value V
}

// New creates a new ordered map.
func New[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		nodeOf: map[K]*list.Element[Pair[K, V]]{},
		values: list.New[Pair[K, V]](),
	}
}

// OrderedMap is an implementation of ordered map. It should be initialized with New function.
type OrderedMap[K cmp.Ordered, V any] struct {
	values *list.List[Pair[K, V]]
	nodeOf map[K]*list.Element[Pair[K, V]]
}

// Get returns the value for the provided key and whether the key presents in the map or not.
func (om *OrderedMap[K, V]) Get(key K) (value V, found bool) {
	node, found := om.nodeOf[key]
	if !found {
		return
	}

	return node.Value.Value, true
}

// Set inserts a new key, value into the map or replaces it if the key presents in the map.
func (om *OrderedMap[K, V]) Set(key K, value V) (oldVal V, replaced bool) {
	node, found := om.nodeOf[key]
	if !found {
		om.nodeOf[key] = om.values.PushBack(Pair[K, V]{
			Key:   key,
			Value: value,
		})
		return
	}

	oldVal = node.Value.Value
	om.values.Remove(node)
	om.nodeOf[key] = om.values.PushBack(Pair[K, V]{
		Key:   key,
		Value: value,
	})
	return oldVal, true
}

// Len returns the size of the map.
func (om *OrderedMap[K, V]) Len() int {
	return om.values.Len()
}

// Delete deletes a key. It returns the deleted value.
func (om *OrderedMap[K, V]) Delete(key K) (val V, present bool) {
	node, found := om.nodeOf[key]
	if !found {
		return
	}

	val = node.Value.Value
	om.values.Remove(node)
	delete(om.nodeOf, key)
	return val, true
}

// MoveAfter moves key to a new position after markedKey.
func (om *OrderedMap[K, V]) MoveAfter(key, markedKey K) error {
	node, found := om.nodeOf[key]
	if !found {
		return ErrKeyNotFound
	}
	markedNode, found := om.nodeOf[markedKey]
	if !found {
		return ErrKeyNotFound
	}

	om.values.MoveAfter(node, markedNode)
	return nil
}

// MoveAfter moves key to a new position after markedKey.
func (om *OrderedMap[K, V]) MoveBefore(key, markedKey K) error {
	node, found := om.nodeOf[key]
	if !found {
		return ErrKeyNotFound
	}
	markedNode, found := om.nodeOf[markedKey]
	if !found {
		return ErrKeyNotFound
	}

	om.values.MoveBefore(node, markedNode)
	return nil
}

// Front returns the pair of key and value at the front of the list.
func (om *OrderedMap[K, V]) Front() (K, V) {
	frontNode := om.values.Front()
	return frontNode.Value.Key, frontNode.Value.Value
}

// Back returns the pair of key and value at the back of the list.
func (om *OrderedMap[K, V]) Back() (K, V) {
	frontNode := om.values.Back()
	return frontNode.Value.Key, frontNode.Value.Value
}

// Scan scans through the map in in the stored order.
func (om *OrderedMap[K, V]) Scan(itor func(key K, val V) bool) {
	for node := om.values.Front(); node != nil; node = node.Next() {
		if !itor(node.Value.Key, node.Value.Value) {
			return
		}
	}
}

// ReverseScan scans through the map in in the reverse of the stored order.
func (om *OrderedMap[K, V]) ReverseScan(itor func(key K, val V) bool) {
	for node := om.values.Back(); node != nil; node = node.Prev() {
		if !itor(node.Value.Key, node.Value.Value) {
			return
		}
	}
}
