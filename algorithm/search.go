package algorithm

import (
	"cmp"
	"sort"
)

// SearchFunc searches for target in a sorted array of values using less
// and return the smallest index i which satisfies !less(values[i], target).
func SearchFunc[T any](values []T, target T, less LessFunc[T]) int {
	return sort.Search(len(values), func(i int) bool {
		return !less(values[i], target)
	})
}

// Search searches for target in a sorted array of values
// and return the smallest index i which satisfies values[i] >= target.
func Search[T cmp.Ordered](values []T, target T) int {
	return sort.Search(len(values), func(i int) bool {
		return !cmp.Less(values[i], target)
	})
}
