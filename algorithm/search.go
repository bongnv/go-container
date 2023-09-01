package algorithm

import (
	"cmp"
	"sort"
)

// Search searches for target in a sorted array of values using less
// and return the smallest index i which satisfies !less(values[i], target).
func Search[T any](values []T, target T, less LessFunc[T]) int {
	return sort.Search(len(values), func(i int) bool {
		return !less(values[i], target)
	})
}

// SearchOrdered searches for target in a sorted array of values
// and return the smallest index i which satisfies values[i] >= target.
func SearchOrdered[T cmp.Ordered](values []T, target T) int {
	return sort.Search(len(values), func(i int) bool {
		return !cmp.Less(values[i], target)
	})
}
