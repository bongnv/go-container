package set

// New creates a new Set.
func New[T comparable]() *Set[T] {
	return &Set[T]{
		container: make(map[T]struct{}),
	}
}

// Set implements the set data structure.
type Set[T comparable] struct {
	container map[T]struct{}
}

// Len returns the size of the set.
func (s Set[T]) Len() int {
	return len(s.container)
}

// Insert inserts a new value into the set.
func (s *Set[T]) Insert(val T) {
	s.container[val] = struct{}{}
}

// Delete deletes a key from a set. If key doesn't exist, it's a no-op.
func (s *Set[T]) Delete(val T) {
	delete(s.container, val)
}

// Contain checks whether the set contains the given value or not.
func (s *Set[T]) Has(val T) bool {
	_, found := s.container[val]
	return found
}

// Scan scans through the set in an arbitrary order.
func (s *Set[T]) Scan(itor func(val T) bool) {
	for val := range s.container {
		if !itor(val) {
			return
		}
	}
}

// Empty returns whether the queue is empty or not.
func (s *Set[T]) Empty() bool {
	return s.Len() == 0
}
