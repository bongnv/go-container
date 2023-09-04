package heap_test

import (
	"testing"

	"github.com/bongnv/go-container/heap"
	"github.com/google/go-cmp/cmp"
)

type Custom struct {
	Key int
	Val string
}

func TestHeap(t *testing.T) {
	testCases := map[string]struct {
		scenario     func(h *heap.Heap[*Custom])
		expectedData *Custom
	}{
		"should maintain the heap properly": {
			scenario: func(h *heap.Heap[*Custom]) {
				h.Push(&Custom{3, "three"})
				h.Push(&Custom{1, "one"})
				h.Push(&Custom{2, "two"})
			},
			expectedData: &Custom{1, "one"},
		},
		"should maintain the heap properly after updating": {
			scenario: func(h *heap.Heap[*Custom]) {
				h.Push(&Custom{3, "three"})
				one := h.Push(&Custom{1, "one"})
				h.Push(&Custom{2, "two"})
				one.Value.Key = 4
				h.Fix(one)
			},
			expectedData: &Custom{2, "two"},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			h := heap.NewFunc[*Custom](func(x, y *Custom) bool {
				return x.Key < y.Key
			})
			tc.scenario(h)
			ret := h.Top().Value
			if diff := cmp.Diff(ret, tc.expectedData); diff != "" {
				t.Errorf("Unexpected result, (+got|-wanted): %s", diff)
			}
		})
	}
}
