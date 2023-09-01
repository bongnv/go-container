package algorithm_test

import (
	"testing"

	"github.com/bongnv/go-container/algorithm"
	gocmp "github.com/google/go-cmp/cmp"
)

func TestSearch(t *testing.T) {
	testCases := map[string]struct {
		input    []int
		expected int
		target   int
	}{
		"should be correct if the target is found": {
			input:    []int{1, 2, 3},
			target:   2,
			expected: 1,
		},
		"should return the smallest index when the target isn't found": {
			input:    []int{1, 2, 4},
			target:   3,
			expected: 2,
		},
		"should return the smallest index when the target is smaller than all elements": {
			input:    []int{1, 2, 4},
			target:   0,
			expected: 0,
		},
		"should return the size of the array when the target is bigger than all elements": {
			input:    []int{1, 2, 4},
			target:   5,
			expected: 3,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			index := algorithm.Search(tc.input, tc.target)
			if diff := gocmp.Diff(tc.expected, index); diff != "" {
				t.Fatalf("wrong index is returned: %s", diff)
			}
		})
	}
}
