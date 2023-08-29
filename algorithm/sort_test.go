package algorithm_test

import (
	"cmp"
	"testing"

	"github.com/bongnv/go-container/algorithm"
)

func TestSort(t *testing.T) {
	t.Run("should be fine if the array is sorted", func(t *testing.T) {
		vals := []int{1, 2, 3}
		algorithm.Sort(vals, cmp.Less)
		if vals[0] != 1 || vals[1] != 2 || vals[2] != 3 {
			t.Fatalf("the array isn't sorted")
		}
	})

	t.Run("should sort if the array isn't sorted", func(t *testing.T) {
		vals := []int{3, 1, 2}
		algorithm.Sort(vals, cmp.Less)
		if vals[0] != 1 || vals[1] != 2 || vals[2] != 3 {
			t.Fatalf("the array isn't sorted")
		}
	})
}
