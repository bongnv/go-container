package set_test

import (
	"testing"

	"github.com/bongnv/go-container/set"
	"github.com/google/go-cmp/cmp"
)

func TestSet(t *testing.T) {
	t.Run("set should work properly", func(t *testing.T) {
		s := set.New[string]()
		s.Insert("a")
		s.Insert("b")
		s.Insert("c")
		if diff := cmp.Diff(s.Len(), 3); diff != "" {
			t.Fatalf("Incorrect size: %v", diff)
		}

		if diff := cmp.Diff(s.Contain("a"), true); diff != "" {
			t.Fatal(diff)
		}

		if diff := cmp.Diff(s.Contain("d"), false); diff != "" {
			t.Fatal(diff)
		}

		s.Delete("b")
		if diff := cmp.Diff(s.Contain("b"), false); diff != "" {
			t.Fatal(diff)
		}

		if diff := cmp.Diff(s.Len(), 2); diff != "" {
			t.Fatalf("Incorrect size: %v", diff)
		}
	})
}
