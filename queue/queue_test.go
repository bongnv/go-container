package queue_test

import (
	"testing"

	"github.com/bongnv/go-container/queue"
)

func TestQueue(t *testing.T) {
	t.Run("queue should work properly", func(t *testing.T) {
		h := queue.New[int]()
		h.Push(1)
		if h.Front() != 1 {
			t.Fatalf("expected 1 but got %v", h.Front())
		}
		h.Push(2)
		if v := h.Pop(); v != 1 {
			t.Fatalf("exected 1 but got %v", v)
		}

		h.Push(3)
		if h.Back() != 3 {
			t.Fatalf("expected 3 but got %v", h.Back())
		}

		if h.Size() != 2 {
			t.Fatalf("expected 2 but got %v", h.Size())
		}

		if h.Front() != 2 {
			t.Fatalf("expected 2 but got %v", h.Front())
		}
	})
}
