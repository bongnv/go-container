package stack_test

import (
	"testing"

	"github.com/bongnv/go-container/stack"
)

func TestStack(t *testing.T) {
	t.Run("stack should work properly", func(t *testing.T) {
		h := stack.New[int]()
		h.Push(1)
		if h.Top() != 1 {
			t.Fatalf("expected 1 but got %v", h.Top())
		}
		h.Push(2)
		if v := h.Pop(); v != 2 {
			t.Fatalf("exected 2 but got %v", v)
		}

		h.Push(3)
		if h.Top() != 3 {
			t.Fatalf("expected 3 but got %v", h.Top())
		}

		if h.Len() != 2 {
			t.Fatalf("expected 2 but got %v", h.Len())
		}
	})
}
