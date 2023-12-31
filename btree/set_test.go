package btree

import (
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	N := 1_000_000
	var tr Set[int]
	for i := 0; i < N; i++ {
		tr.Load(i)
	}
	assert(t, tr.Len() == N)
	for i := 0; i < N; i++ {
		assert(t, tr.Has(i))
	}

	count := 0
	tr.Scan(func(_ int) bool {
		count++
		return true
	})
	assert(t, count == N)
	count = 0
	tr.Ascend(N/2, func(_ int) bool {
		count++
		return true
	})
	assert(t, count == N/2)

	count = 0
	tr.Reverse(func(_ int) bool {
		count++
		return true
	})
	assert(t, count == N)
	count = 0
	tr.Descend(N/2, func(_ int) bool {
		count++
		return true
	})
	assert(t, count == N/2+1)

	for i := 0; i < N; i++ {
		tr.Delete(i)
	}

	dotup := func(v int, ok bool) interface{} {
		if !ok {
			return nil
		}
		return v
	}

	assert(t, tr.Len() == 0)
	assert(t, dotup(tr.Min()) == nil)
	assert(t, dotup(tr.Max()) == nil)
	assert(t, dotup(tr.DeleteMin()) == nil)
	assert(t, dotup(tr.DeleteMax()) == nil)
	for i := 0; i < N; i++ {
		assert(t, !tr.Has(i))
	}
	for i := 0; i < N; i++ {
		tr.Insert(i)
	}
	assert(t, tr.Len() == N)
	for i := 0; i < N; i++ {
		tr.Insert(i)
	}
	assert(t, tr.Len() == N)
	for i := 0; i < N; i++ {
		tr.Load(i)
	}
	assert(t, tr.Len() == N)
	assert(t, dotup(tr.Min()) == 0)
	assert(t, dotup(tr.Max()) == N-1)
	assert(t, dotup(tr.DeleteMin()) == 0)
	assert(t, dotup(tr.DeleteMax()) == N-1)
	tr.Insert(0)
	tr.Insert(N - 1)
	assert(t, dotup(tr.GetAt(0)) == 0)
	assert(t, dotup(tr.GetAt(N)) == nil)
	tr.Insert(N - 1)
	assert(t, tr.Height() > 0)
	assert(t, dotup(tr.DeleteAt(0)) == 0)
	tr.Insert(0)
	assert(t, dotup(tr.DeleteAt(N-1)) == N-1)
	assert(t, dotup(tr.DeleteAt(N)) == nil)
	tr.Insert(N - 1)

	count = 0
	tr.Scan(func(item int) bool {
		count++
		return true
	})

	assert(t, count == N)

	for i := 0; i < N; i++ {
		assert(t, tr.Has(i))
	}
	for i := 0; i < N; i++ {
		tr.Delete(i)
	}
	for i := 0; i < N; i++ {
		assert(t, !tr.Has(i))
	}
	assert(t, tr.base.lt(1, 2))
	assert(t, tr.base.lt(2, 10))
}

func TestSetClear(t *testing.T) {
	var tr Set[int]
	for i := 0; i < 100; i++ {
		tr.Insert(i)
	}
	assert(t, tr.Len() == 100)
	tr.Clear()
	assert(t, tr.Len() == 0)
	for i := 0; i < 100; i++ {
		tr.Insert(i)
	}
	assert(t, tr.Len() == 100)
}

func copySetEntries(m *Set[int]) []int {
	all := m.Keys()
	sort.Ints(all)
	return all
}

func setEntriesEqual(a, b []int) bool {
	return reflect.DeepEqual(a, b)
}

func copySetTest(N int, s1 *Set[int], e11 []int, deep bool) {
	e12 := copySetEntries(s1)
	if !setEntriesEqual(e11, e12) {
		panic("!")
	}

	// Make a copy and compare the values
	s2 := s1.Copy()
	e21 := copySetEntries(s1)
	if !setEntriesEqual(e21, e12) {
		panic("!")
	}

	// Delete every other key
	var e22 []int
	for i, j := range rand.Perm(N) {
		if i&1 == 0 {
			e22 = append(e22, e21[j])
		} else {
			s2.Delete(e21[j])
		}
	}

	if s2.Len() != N/2 {
		panic("!")
	}
	sort.Ints(e22)
	e23 := copySetEntries(s2)
	if !setEntriesEqual(e23, e22) {
		panic("!")
	}
	if !deep {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			copySetTest(N/2, s2, e23, true)
		}()
		go func() {
			defer wg.Done()
			copySetTest(N/2, s2, e23, true)
		}()
		wg.Wait()
	}
	e24 := copySetEntries(s2)
	if !setEntriesEqual(e24, e23) {
		panic("!")
	}
}

func TestSetCopy(t *testing.T) {
	N := 1_000
	// create the initial map

	s1 := new(Set[int])
	for s1.Len() < N {
		s1.Insert(rand.Int())
	}
	e11 := copySetEntries(s1)
	dur := time.Second * 2
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			for time.Since(start) < dur {
				copySetTest(N, s1, e11, false)
			}
		}()
	}
	wg.Wait()
	e12 := copySetEntries(s1)
	if !setEntriesEqual(e11, e12) {
		panic("!")
	}
}
