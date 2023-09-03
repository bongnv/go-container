package heap_test

import (
	"testing"

	"github.com/bongnv/go-container/heap"
)

type Custom struct {
	Key int
	Val string
}

func TestHeap(t *testing.T) {
	testCases := map[string]struct {
		scenario         func(h *heap.Heap[*Custom])
		expectedTopKey   int
		expectedTopValue string
	}{
		"should maintain the heap properly": {
			scenario: func(h *heap.Heap[*Custom]) {
				h.Push(&Custom{3, "three"})
				h.Push(&Custom{1, "one"})
				h.Push(&Custom{2, "two"})
			},
			expectedTopKey:   1,
			expectedTopValue: "one",
		},
		"should maintain the heap properly after updating": {
			scenario: func(h *heap.Heap[*Custom]) {
				one := &Custom{1, "one"}
				h.Push(&Custom{3, "three"})
				h.Push(one)
				h.Push(&Custom{2, "two"})
				one.Key = 4
				h.Fix(one)
			},
			expectedTopKey:   2,
			expectedTopValue: "two",
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			h := heap.NewFunc[*Custom](func(x, y *Custom) bool {
				return x.Key < y.Key
			})
			tc.scenario(h)
			ret := h.Top()
			if ret.Key != tc.expectedTopKey {
				t.Errorf("expected key %v, but got: %v", tc.expectedTopKey, ret.Key)
			}

			if ret.Val != tc.expectedTopValue {
				t.Errorf("expected value %v, but got: %v", tc.expectedTopValue, ret.Val)
			}
		})
	}
}
