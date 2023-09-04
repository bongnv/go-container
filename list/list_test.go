package list_test

import (
	"testing"

	"github.com/bongnv/go-container/list"
)

func TestList(t *testing.T) {
	t.Run("should add and remove nodes properly", func(t *testing.T) {
		l := list.New[string]()
		l.PushBack("bong")
		l.PushBack("c")
		l.PushFront("a")
		expectList(t, l, "a", "bong", "c")
		toRemove := l.Back()
		l.PushBack("d")
		l.Delete(toRemove)
		expectList(t, l, "a", "bong", "d")
	})
}

func expectList(t *testing.T, l *list.List[string], elements ...string) {
	if l.Len() != len(elements) {
		t.Errorf("Expected size %v but got %v", len(elements), l.Len())
		return
	}

	for e, i := l.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		if e.Value != elements[i] {
			t.Errorf("Expected %v but got %v", elements[i], e.Value)
			return
		}
	}
}
