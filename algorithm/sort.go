package algorithm

import (
	"cmp"
	"sort"
)

// LessFunc is a function that returns whether x < y or not.
type LessFunc[T any] func(x, y T) bool

// SortFunc sorts an array using less.
func SortFunc[T any](values []T, less LessFunc[T]) {
	sort.Sort(&sortableContainer[T]{
		values: values,
		less:   less,
	})
}

type sortableContainer[T any] struct {
	values []T
	less   LessFunc[T]
}

func (sc sortableContainer[T]) Len() int {
	return len(sc.values)
}

func (sc sortableContainer[T]) Less(i, j int) bool {
	return sc.less(sc.values[i], sc.values[j])
}

func (sc *sortableContainer[T]) Swap(i, j int) {
	sc.values[i], sc.values[j] = sc.values[j], sc.values[i]
}

// Sort sorts an array of values from ordered types like int, float, etc....
func Sort[T cmp.Ordered](values []T) {
	SortFunc(values, cmp.Less[T])
}
