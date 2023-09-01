package algorithm_test

import (
	"cmp"
	"testing"

	"github.com/bongnv/go-container/algorithm"
	gocmp "github.com/google/go-cmp/cmp"
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

func TestSortOrdered(t *testing.T) {
	testCases := map[string]struct {
		input    []int
		expected []int
	}{
		"should be fine if the array is sorted": {
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		"should sort if the array isn't sorted": {
			input:    []int{3, 1, 2},
			expected: []int{1, 2, 3},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			algorithm.SortOrdered(tc.input)
			if diff := gocmp.Diff(tc.expected, tc.input); diff != "" {
				t.Fatalf("the array isn't sorted: %s", diff)
			}
		})
	}
}
