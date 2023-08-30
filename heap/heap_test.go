package heap_test

import (
	"cmp"
	"testing"

	"github.com/bongnv/go-container/heap"
)

func TestHeap(t *testing.T) {
	t.Run("min heap should work properly", func(t *testing.T) {
		h := heap.New[int](cmp.Less)
		h.Push(1)
		h.Push(2)
		h.Push(3)

		if h.Top() != 1 {
			t.Fatalf("expected 1 but got %v", h.Top())
		}

		if v := h.Pop(); v != 1 {
			t.Fatalf("exected 1 but got %v", v)
		}

		if h.Top() != 2 {
			t.Fatalf("expected 2 but got %v", h.Top())
		}

		if h.Len() != 2 {
			t.Fatalf("expected 2 but got %v", h.Len())
		}
	})

	t.Run("max heap should work properly", func(t *testing.T) {
		h := heap.New[int](greater)
		h.Push(1)
		h.Push(2)
		h.Push(3)

		if h.Top() != 3 {
			t.Fatalf("expected 3 but got %v", h.Top())
		}

		if v := h.Pop(); v != 3 {
			t.Fatalf("exected 3 but got %v", v)
		}

		if h.Top() != 2 {
			t.Fatalf("expected 2 but got %v", h.Top())
		}

		if h.Len() != 2 {
			t.Fatalf("expected 2 but got %v", h.Len())
		}
	})

	t.Run("heap should work fine with duplicates", func(t *testing.T) {
		h := heap.New[int](cmp.Less)
		h.Push(1)
		h.Push(2)
		h.Push(3)
		h.Push(1)

		if h.Top() != 1 {
			t.Fatalf("expected 1 but got %v", h.Top())
		}

		if v := h.Pop(); v != 1 {
			t.Fatalf("exected 1 but got %v", v)
		}

		if h.Top() != 1 {
			t.Fatalf("expected 1 but got %v", h.Top())
		}

		if h.Len() != 3 {
			t.Fatalf("expected 3 but got %v", h.Len())
		}
	})

	t.Run("heap should work with custom data structure", func(t *testing.T) {
		h := heap.New[*Custom](func(x, y *Custom) bool {
			return x.Value < y.Value
		})

		origin := &Custom{}
		h.Push(origin)
		h.Push(&Custom{Value: 1})
		h.Push(&Custom{Value: 2})
		h.Push(&Custom{Value: 3})

		if h.Top().Value != 0 {
			t.Fatalf("expected 0 but got %v", h.Top().Value)
		}

		origin.Value = 2
		h.Fix(origin)
		if h.Top().Value != 1 {
			t.Fatalf("expected 1 but got %v", h.Top().Value)
		}

		origin.Value = 0
		h.Fix(origin)
		if h.Top().Value != 0 {
			t.Fatalf("expected 1 but got %v", h.Top().Value)
		}
	})
}

func greater[T cmp.Ordered](x, y T) bool {
	return x > y
}

type Custom struct {
	Value int
}
