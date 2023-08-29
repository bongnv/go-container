package algorithm

import "sort"

// Less is a function that returns whether x < y or not.
type Less[T any] func(x, y T) bool

// Sort sorts an array using less.
func Sort[T any](values []T, less Less[T]) {
	sort.Sort(&sortableContainer[T]{
		values: values,
		less:   less,
	})
}

type sortableContainer[T any] struct {
	values []T
	less   Less[T]
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
