package algorithm_test

import (
	"testing"

	"github.com/bongnv/go-container/algorithm"
	gocmp "github.com/google/go-cmp/cmp"
)

func TestSort(t *testing.T) {
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
			algorithm.Sort(tc.input)
			if diff := gocmp.Diff(tc.expected, tc.input); diff != "" {
				t.Fatalf("the array isn't sorted: %s", diff)
			}
		})
	}
}
