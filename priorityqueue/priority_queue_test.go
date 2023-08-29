package priorityqueue_test

import (
	"cmp"
	"testing"

	"github.com/bongnv/go-container/priorityqueue"
)

func TestPriorityQueue(t *testing.T) {
	t.Run("min heap should work properly", func(t *testing.T) {
		h := priorityqueue.New[int](cmp.Less)
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

		if h.Size() != 2 {
			t.Fatalf("expected 2 but got %v", h.Size())
		}
	})

	t.Run("max heap should work properly", func(t *testing.T) {
		h := priorityqueue.New[int](greater)
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

		if h.Size() != 2 {
			t.Fatalf("expected 2 but got %v", h.Size())
		}
	})

	t.Run("heap should work fine with duplicates", func(t *testing.T) {
		h := priorityqueue.New[int](cmp.Less)
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

		if h.Size() != 3 {
			t.Fatalf("expected 3 but got %v", h.Size())
		}
	})

	t.Run("heap should work with custom data structure", func(t *testing.T) {
		h := priorityqueue.New[*Custom](func(x, y *Custom) bool {
			return x.Value < y.Value
		})

		h.Push(&Custom{Value: 2})
		if h.Top().Value != 2 {
			t.Fatalf("expected 2 but got %v", h.Top().Value)
		}

		h.Push(&Custom{Value: 1})
		h.Push(&Custom{Value: 3})

		if h.Top().Value != 1 {
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
