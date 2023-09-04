package orderedmap_test

import (
	"testing"

	"github.com/bongnv/go-container/orderedmap"
	"github.com/google/go-cmp/cmp"
)

func TestOrderedMap_Scan(t *testing.T) {
	testCases := map[string]struct {
		scenario      func(om *orderedmap.OrderedMap[int, string])
		expectedPairs []orderedmap.Pair[int, string]
	}{
		"should maintain the inserted order": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{1, "one"},
				{2, "two"},
				{3, "three"},
			},
		},
		"should be able to MoveBefore correctly": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
				om.MoveBefore(3, 2)
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{1, "one"},
				{3, "three"},
				{2, "two"},
			},
		},
		"should be able to MoveAfter correctly": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
				om.MoveAfter(3, 1)
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{1, "one"},
				{3, "three"},
				{2, "two"},
			},
		},
		"should be able to MoveToFront correctly": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
				om.MoveToFront(2)
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{2, "two"},
				{1, "one"},
				{3, "three"},
			},
		},
		"should be able to MoveToBack correctly": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
				om.MoveToBack(2)
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{1, "one"},
				{3, "three"},
				{2, "two"},
			},
		},
		"should be able to Delete correctly": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
				om.MoveAfter(3, 1)
				om.Set(4, "four")
				om.Delete(2)
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{1, "one"},
				{3, "three"},
				{4, "four"},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			om := orderedmap.New[int, string]()
			tc.scenario(om)
			allPairs := make([]orderedmap.Pair[int, string], 0, om.Len())
			om.Scan(func(key int, val string) bool {
				allPairs = append(allPairs, orderedmap.Pair[int, string]{key, val})
				return true
			})
			if diff := cmp.Diff(allPairs, tc.expectedPairs); diff != "" {
				t.Errorf("Unexpected result (+got,-wanted): %v", diff)
			}
		})
	}
}

func TestOrderedMap_ReverseScan(t *testing.T) {
	testCases := map[string]struct {
		scenario      func(om *orderedmap.OrderedMap[int, string])
		expectedPairs []orderedmap.Pair[int, string]
	}{
		"should maintain the inserted order": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{3, "three"},
				{2, "two"},
				{1, "one"},
			},
		},
		"should be able to MoveBefore correctly": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
				om.MoveBefore(3, 2)
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{2, "two"},
				{3, "three"},
				{1, "one"},
			},
		},
		"should be able to MoveAfter correctly": {
			scenario: func(om *orderedmap.OrderedMap[int, string]) {
				om.Set(1, "one")
				om.Set(2, "two")
				om.Set(3, "three")
				om.MoveAfter(3, 1)
			},
			expectedPairs: []orderedmap.Pair[int, string]{
				{2, "two"},
				{3, "three"},
				{1, "one"},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			om := orderedmap.New[int, string]()
			tc.scenario(om)
			allPairs := make([]orderedmap.Pair[int, string], 0, om.Len())
			om.ReverseScan(func(key int, val string) bool {
				allPairs = append(allPairs, orderedmap.Pair[int, string]{key, val})
				return true
			})
			if diff := cmp.Diff(allPairs, tc.expectedPairs); diff != "" {
				t.Errorf("Unexpected result (+got,-wanted): %v", diff)
			}
		})
	}
}

func TestOrderedMap(t *testing.T) {
	om := orderedmap.New[int, string]()
	om.Set(1, "one")
	two, replaced := om.Set(2, "two")
	if replaced || two != "" {
		t.Errorf("Set doesn't return values properly, replaced: %v, value: %v", replaced, two)
	}

	om.Set(3, "three")
	frontKey, frontVal := om.Front()
	if frontKey != 1 || frontVal != "one" {
		t.Errorf("Invalid front values")
	}

	backKey, backVal := om.Back()
	if backKey != 3 || backVal != "three" {
		t.Errorf("invalid back values")
	}

	replacedVal, present := om.Set(2, "new-two")
	if replacedVal != "two" || !present {
		t.Errorf("Set returns invalid values")
	}

	deletedVal, present := om.Delete(2)
	if deletedVal != "new-two" || !present {
		t.Errorf("Delete returns invalid values")
	}
}
