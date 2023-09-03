package rbtree_test

import (
	"math/rand"
	"testing"

	"github.com/bongnv/go-container/rbtree"
	"github.com/google/go-cmp/cmp"
)

func TestCases(t *testing.T) {
	tree := rbtree.New[int]()
	tree.Upsert(1)
	tree.Upsert(1)
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has(1) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(1) {
		t.Errorf("not expecting to find key=1")
	}
}

func TestReverseInsertOrder(t *testing.T) {
	tree := rbtree.New[int]()
	n := 100
	for i := 0; i < n; i++ {
		tree.Upsert(n - i)
	}
	i := 0
	tree.AscendGreaterOrEqual(0, func(item int) bool {
		i++
		if item != i {
			t.Errorf("bad order: got %d, expect %d", item, i)
		}
		return true
	})
}

func TestRange(t *testing.T) {
	tree := rbtree.New[string]()
	order := []string{
		"ab", "aba", "abc", "a", "aa", "aaa", "b", "a-", "a!",
	}
	for _, i := range order {
		tree.Upsert(i)
	}
	k := 0
	tree.AscendRange("ab", "ac", func(item string) bool {
		if k > 3 {
			t.Fatalf("returned more items than expected")
		}
		i1 := order[k]
		i2 := item
		if i1 != i2 {
			t.Errorf("expecting %s, got %s", i1, i2)
		}
		k++
		return true
	})
}

func TestRandomInsertOrder(t *testing.T) {
	tree := rbtree.New[int]()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(perm[i])
	}
	j := 0
	tree.AscendGreaterOrEqual(0, func(item int) bool {
		if item != j {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestRandomReplace(t *testing.T) {
	tree := rbtree.New[int]()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(perm[i])
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		if replacedItem, replaced := tree.Upsert(perm[i]); !replaced || replacedItem != perm[i] {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	tree := rbtree.New[int]()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(perm[i])
	}
	for i := 0; i < n; i++ {
		tree.Delete(i)
	}
}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	tree := rbtree.New[int]()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(perm[i])
	}
	if _, deleted := tree.Delete(200); deleted {
		t.Errorf("deleted non-existent item")
	}
	if _, deleted := tree.Delete(-2); deleted {
		t.Errorf("deleted non-existent item")
	}
	for i := 0; i < n; i++ {
		if u, deleted := tree.Delete(i); !deleted || u != i {
			t.Errorf("delete failed")
		}
	}
	if _, deleted := tree.Delete(200); deleted {
		t.Errorf("deleted non-existent item")
	}
	if _, deleted := tree.Delete(-2); deleted {
		t.Errorf("deleted non-existent item")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	tree := rbtree.New[int]()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(perm[i])
	}
	for i := 1; i < n-1; i++ {
		tree.Delete(i)
	}
	j := 0
	tree.AscendGreaterOrEqual(0, func(item int) bool {
		switch j {
		case 0:
			if item != 0 {
				t.Errorf("expecting 0")
			}
		case 1:
			if item != n-1 {
				t.Errorf("expecting %d", n-1)
			}
		}
		j++
		return true
	})
}

func TestInsertNoReplace(t *testing.T) {
	tree := rbtree.New[int]()
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			tree.Insert(perm[i])
		}
	}
	j := 0
	tree.AscendGreaterOrEqual(0, func(item int) bool {
		if item != j/2 {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestLLRB(t *testing.T) {
	testCases := map[string]struct {
		scenario      func(t *rbtree.LLRB[int])
		expectedOrder []int
	}{
		"should be ordered properly": {
			scenario: func(t *rbtree.LLRB[int]) {
				t.Insert(1)
				t.Insert(0)
				t.Insert(2)
				t.Insert(2)
				t.Insert(4)
			},
			expectedOrder: []int{0, 1, 2, 2, 4},
		},
		"should be ordered properly after deleting": {
			scenario: func(t *rbtree.LLRB[int]) {
				t.Insert(1)
				t.Insert(0)
				t.Insert(2)
				t.Insert(4)
				t.Delete(2)
			},
			expectedOrder: []int{0, 1, 4},
		},
		"should be ordered properly after deleting 3 items": {
			scenario: func(t *rbtree.LLRB[int]) {
				t.Insert(1)
				t.Insert(0)
				t.Insert(0)
				t.Insert(2)
				t.Insert(4)
				t.Insert(4)
				t.Delete(2)
				t.Delete(4)
				t.Delete(0)
			},
			expectedOrder: []int{0, 1, 4},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			rbt := rbtree.New[int]()
			tc.scenario(rbt)
			allItems := make([]int, 0, rbt.Len())
			rbt.Scan(func(i int) bool {
				allItems = append(allItems, i)
				return true
			})
			if diff := cmp.Diff(allItems, tc.expectedOrder); diff != "" {
				t.Errorf("unexpected order (+got, -wanted): %v", diff)
			}
		})
	}
}

func TestLLRB_ReverseScan(t *testing.T) {
	testCases := map[string]struct {
		scenario      func(t *rbtree.LLRB[int])
		expectedOrder []int
	}{
		"should be ordered properly": {
			scenario: func(t *rbtree.LLRB[int]) {
				t.Insert(1)
				t.Insert(0)
				t.Insert(2)
				t.Insert(2)
				t.Insert(4)
			},
			expectedOrder: []int{4, 2, 2, 1, 0},
		},
		"should be ordered properly after deleting": {
			scenario: func(t *rbtree.LLRB[int]) {
				t.Insert(1)
				t.Insert(0)
				t.Insert(2)
				t.Insert(4)
				t.Delete(2)
			},
			expectedOrder: []int{4, 1, 0},
		},
		"should be ordered properly after deleting 2 items": {
			scenario: func(t *rbtree.LLRB[int]) {
				t.Insert(1)
				t.Insert(1)
				t.Insert(0)
				t.Insert(2)
				t.Insert(4)
				t.Delete(2)
				t.Delete(1)
			},
			expectedOrder: []int{4, 1, 0},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			rbt := rbtree.New[int]()
			tc.scenario(rbt)
			allItems := make([]int, 0, rbt.Len())
			rbt.ReverseScan(func(i int) bool {
				allItems = append(allItems, i)
				return true
			})
			if diff := cmp.Diff(allItems, tc.expectedOrder); diff != "" {
				t.Errorf("unexpected order (+got, -wanted): %v", diff)
			}
		})
	}
}
