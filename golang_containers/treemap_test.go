package containers

import (
	"cmp"
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Red-Black structural validator
//
// validateRB checks all five RB invariants on the subtree rooted at n and
// returns the black-height (number of black nodes on any root→nil path).
// It calls t.Fatalf on the first violation so the test stops immediately with
// a descriptive message rather than panicking deep in the tree.
// ---------------------------------------------------------------------------

func validateRB[K cmp.Ordered, V any](t *testing.T, n *rbNode[K, V], parent *rbNode[K, V]) int {
	t.Helper()
	if n == nil {
		return 1 // nil counts as one black node (invariant 3)
	}

	// Invariant 4: red node must not have a red parent.
	if isRed(n) && isRed(parent) {
		t.Fatalf("RB violation: consecutive red nodes at key %v", n.key)
	}

	// Parent pointer consistency.
	if n.parent != parent {
		t.Fatalf("RB violation: parent pointer mismatch at key %v", n.key)
	}

	// BST ordering.
	if parent != nil {
		if n == parent.left && cmp.Compare(n.key, parent.key) >= 0 {
			t.Fatalf("RB violation: left child key %v >= parent key %v", n.key, parent.key)
		}
		if n == parent.right && cmp.Compare(n.key, parent.key) <= 0 {
			t.Fatalf("RB violation: right child key %v <= parent key %v", n.key, parent.key)
		}
	}

	leftBH := validateRB(t, n.left, n)
	rightBH := validateRB(t, n.right, n)

	// Invariant 5: equal black-height on both sides.
	if leftBH != rightBH {
		t.Fatalf("RB violation: black-height mismatch at key %v (left=%d right=%d)", n.key, leftBH, rightBH)
	}

	if isRed(n) {
		return leftBH
	}
	return leftBH + 1
}

func checkRB[K cmp.Ordered, V any](t *testing.T, m *TreeMap[K, V]) {
	t.Helper()
	// Invariant 2: root is black.
	if m.root != nil && isRed(m.root) {
		t.Fatalf("RB violation: root is red")
	}
	validateRB(t, m.root, nil)
}

// ---------------------------------------------------------------------------
// Group 1: Basic operations
// ---------------------------------------------------------------------------

func Test_TreeMap_EmptyGet(t *testing.T) {
	m := NewTreeMap[int, string]()
	v, ok := m.Get(42)
	assert.False(t, ok)
	assert.Equal(t, "", v)
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Size())
}

func Test_TreeMap_EmptyMinMax(t *testing.T) {
	m := NewTreeMap[int, string]()
	k, v, ok := m.Min()
	assert.False(t, ok)
	assert.Equal(t, 0, k)
	assert.Equal(t, "", v)
	k, v, ok = m.Max()
	assert.False(t, ok)
	assert.Equal(t, 0, k)
	assert.Equal(t, "", v)
}

func Test_TreeMap_SinglePutGet(t *testing.T) {
	m := NewTreeMap[int, string]()
	m.Put(7, "seven")
	assert.Equal(t, 1, m.Size())
	assert.False(t, m.IsEmpty())
	assert.True(t, m.Contains(7))

	v, ok := m.Get(7)
	assert.True(t, ok)
	assert.Equal(t, "seven", v)

	checkRB(t, m)
}

func Test_TreeMap_UpdateExistingKey(t *testing.T) {
	m := NewTreeMap[int, string]()
	m.Put(1, "one")
	m.Put(1, "ONE")
	assert.Equal(t, 1, m.Size())
	v, ok := m.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "ONE", v)
}

func Test_TreeMap_DeleteMissing(t *testing.T) {
	m := NewTreeMap[int, string]()
	m.Put(1, "one")
	m.Delete(99) // no-op, must not panic
	assert.Equal(t, 1, m.Size())
}

func Test_TreeMap_PutDeleteGet(t *testing.T) {
	m := NewTreeMap[int, string]()
	m.Put(5, "five")
	m.Delete(5)
	assert.Equal(t, 0, m.Size())
	assert.True(t, m.IsEmpty())
	_, ok := m.Get(5)
	assert.False(t, ok)
	checkRB(t, m)
}

// ---------------------------------------------------------------------------
// Group 2: Ordering invariants
// ---------------------------------------------------------------------------

func Test_TreeMap_KeysAscending(t *testing.T) {
	keys := []int{5, 3, 8, 1, 9, 2, 7, 4, 6, 10}
	m := NewTreeMap[int, int]()
	for _, k := range keys {
		m.Put(k, k*10)
	}

	got := m.Keys()
	want := make([]int, len(keys))
	copy(want, keys)
	sort.Ints(want)

	assert.Equal(t, want, got)
	checkRB(t, m)
}

func Test_TreeMap_ValuesKeyOrder(t *testing.T) {
	m := NewTreeMap[int, string]()
	m.Put(3, "three")
	m.Put(1, "one")
	m.Put(2, "two")

	assert.Equal(t, []string{"one", "two", "three"}, m.Values())
}

func Test_TreeMap_MinMax(t *testing.T) {
	m := NewTreeMap[int, int]()
	for _, k := range []int{4, 2, 7, 1, 5} {
		m.Put(k, k)
	}

	minK, _, ok := m.Min()
	assert.True(t, ok)
	assert.Equal(t, 1, minK)

	maxK, _, ok := m.Max()
	assert.True(t, ok)
	assert.Equal(t, 7, maxK)
}

func Test_TreeMap_MinAfterDeleteMin(t *testing.T) {
	m := NewTreeMap[int, int]()
	for _, k := range []int{3, 1, 5} {
		m.Put(k, k)
	}
	m.Delete(1)
	minK, _, ok := m.Min()
	assert.True(t, ok)
	assert.Equal(t, 3, minK)
	checkRB(t, m)
}

func Test_TreeMap_MaxAfterDeleteMax(t *testing.T) {
	m := NewTreeMap[int, int]()
	for _, k := range []int{3, 1, 5} {
		m.Put(k, k)
	}
	m.Delete(5)
	maxK, _, ok := m.Max()
	assert.True(t, ok)
	assert.Equal(t, 3, maxK)
	checkRB(t, m)
}

// ---------------------------------------------------------------------------
// Group 3: Floor and Ceiling
// ---------------------------------------------------------------------------

func Test_TreeMap_Floor(t *testing.T) {
	m := NewTreeMap[int, string]()
	for _, k := range []int{1, 3, 5} {
		m.Put(k, fmt.Sprintf("v%d", k))
	}

	tests := []struct {
		query    int
		wantKey  int
		wantFound bool
	}{
		{3, 3, true},  // exact match
		{4, 3, true},  // between keys
		{9, 5, true},  // above all
		{0, 0, false}, // below all
		{1, 1, true},  // equal to minimum
		{5, 5, true},  // equal to maximum
	}

	for _, tc := range tests {
		k, _, ok := m.Floor(tc.query)
		assert.Equal(t, tc.wantFound, ok, "Floor(%d) found", tc.query)
		if tc.wantFound {
			assert.Equal(t, tc.wantKey, k, "Floor(%d) key", tc.query)
		}
	}
}

func Test_TreeMap_Floor_Empty(t *testing.T) {
	m := NewTreeMap[int, int]()
	_, _, ok := m.Floor(5)
	assert.False(t, ok)
}

func Test_TreeMap_Ceiling(t *testing.T) {
	m := NewTreeMap[int, string]()
	for _, k := range []int{1, 3, 5} {
		m.Put(k, fmt.Sprintf("v%d", k))
	}

	tests := []struct {
		query    int
		wantKey  int
		wantFound bool
	}{
		{3, 3, true},  // exact match
		{2, 3, true},  // between keys
		{0, 1, true},  // below all
		{6, 0, false}, // above all
		{1, 1, true},  // equal to minimum
		{5, 5, true},  // equal to maximum
	}

	for _, tc := range tests {
		k, _, ok := m.Ceiling(tc.query)
		assert.Equal(t, tc.wantFound, ok, "Ceiling(%d) found", tc.query)
		if tc.wantFound {
			assert.Equal(t, tc.wantKey, k, "Ceiling(%d) key", tc.query)
		}
	}
}

func Test_TreeMap_Ceiling_Empty(t *testing.T) {
	m := NewTreeMap[int, int]()
	_, _, ok := m.Ceiling(5)
	assert.False(t, ok)
}

// ---------------------------------------------------------------------------
// Group 4: Red-Black structural invariants
// ---------------------------------------------------------------------------

func Test_TreeMap_RBInvariants_AscendingInserts(t *testing.T) {
	m := NewTreeMap[int, int]()
	for i := 0; i < 100; i++ {
		m.Put(i, i)
		checkRB(t, m)
		assert.Equal(t, i+1, m.Size())
	}
}

func Test_TreeMap_RBInvariants_DescendingInserts(t *testing.T) {
	m := NewTreeMap[int, int]()
	for i := 99; i >= 0; i-- {
		m.Put(i, i)
		checkRB(t, m)
	}
	assert.Equal(t, 100, m.Size())
}

func Test_TreeMap_RBInvariants_RandomInserts(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	perm := rng.Perm(100)
	m := NewTreeMap[int, int]()
	for _, k := range perm {
		m.Put(k, k)
		checkRB(t, m)
	}
	assert.Equal(t, 100, m.Size())
}

func Test_TreeMap_RBInvariants_RandomDeletes(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	keys := rng.Perm(50)
	m := NewTreeMap[int, int]()
	for _, k := range keys {
		m.Put(k, k)
	}

	deleteOrder := rng.Perm(50)
	for i, idx := range deleteOrder {
		m.Delete(keys[idx])
		checkRB(t, m)
		assert.Equal(t, 50-i-1, m.Size())
	}
}

// ---------------------------------------------------------------------------
// Group 5: Delete edge cases
// ---------------------------------------------------------------------------

func Test_TreeMap_DeleteRoot_SingleNode(t *testing.T) {
	m := NewTreeMap[int, int]()
	m.Put(1, 1)
	m.Delete(1)
	assert.Equal(t, 0, m.Size())
	assert.Nil(t, m.root)
	checkRB(t, m)
}

func Test_TreeMap_DeleteRoot_WithChildren(t *testing.T) {
	m := NewTreeMap[int, int]()
	m.Put(2, 2)
	m.Put(1, 1)
	m.Put(3, 3)
	m.Delete(2)
	assert.Equal(t, 2, m.Size())
	_, ok := m.Get(2)
	assert.False(t, ok)
	checkRB(t, m)
}

func Test_TreeMap_DeleteBlackLeaf(t *testing.T) {
	// Build a small tree where deleting a black leaf forces deleteFixup.
	m := NewTreeMap[int, int]()
	for _, k := range []int{10, 5, 15, 3, 7} {
		m.Put(k, k)
	}
	m.Delete(3) // likely black leaf
	checkRB(t, m)
	_, ok := m.Get(3)
	assert.False(t, ok)
}

func Test_TreeMap_DeleteNodeWithTwoChildren(t *testing.T) {
	m := NewTreeMap[int, int]()
	for k := 1; k <= 7; k++ {
		m.Put(k, k)
	}
	m.Delete(3) // node with two children — uses in-order successor
	checkRB(t, m)
	_, ok := m.Get(3)
	assert.False(t, ok)
	assert.Equal(t, 6, m.Size())
}

func Test_TreeMap_DeleteAllElements(t *testing.T) {
	m := NewTreeMap[int, int]()
	for k := 0; k < 20; k++ {
		m.Put(k, k)
	}
	for k := 0; k < 20; k++ {
		m.Delete(k)
		checkRB(t, m)
	}
	assert.Equal(t, 0, m.Size())
	assert.Nil(t, m.root)
}

// ---------------------------------------------------------------------------
// Group 6: Range iteration
// ---------------------------------------------------------------------------

func Test_TreeMap_AllEmpty(t *testing.T) {
	m := NewTreeMap[int, int]()
	count := 0
	for range m.All() {
		count++
	}
	assert.Equal(t, 0, count)
}

func Test_TreeMap_AllOrder(t *testing.T) {
	m := NewTreeMap[int, int]()
	for _, k := range []int{5, 3, 8, 1, 9} {
		m.Put(k, k*10)
	}

	var keys []int
	for k := range m.All() {
		keys = append(keys, k)
	}
	assert.Equal(t, []int{1, 3, 5, 8, 9}, keys)
}

func Test_TreeMap_AllKeyValuePairs(t *testing.T) {
	m := NewTreeMap[int, int]()
	for _, k := range []int{2, 1, 3} {
		m.Put(k, k*100)
	}

	pairs := map[int]int{}
	for k, v := range m.All() {
		pairs[k] = v
	}
	assert.Equal(t, map[int]int{1: 100, 2: 200, 3: 300}, pairs)
}

func Test_TreeMap_AllEarlyBreak(t *testing.T) {
	m := NewTreeMap[int, int]()
	for k := 1; k <= 10; k++ {
		m.Put(k, k)
	}

	var keys []int
	for k := range m.All() {
		keys = append(keys, k)
		if k == 3 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3}, keys)
}

func Test_TreeMap_AllAfterDelete(t *testing.T) {
	m := NewTreeMap[int, int]()
	for k := 1; k <= 5; k++ {
		m.Put(k, k)
	}
	m.Delete(2)
	m.Delete(4)

	var keys []int
	for k := range m.All() {
		keys = append(keys, k)
	}
	assert.Equal(t, []int{1, 3, 5}, keys)
}

// ---------------------------------------------------------------------------
// Group 7: Multiple key types (generics)
// ---------------------------------------------------------------------------

func Test_TreeMap_StringKeys(t *testing.T) {
	m := NewTreeMap[string, int]()
	words := []string{"banana", "apple", "cherry", "date", "elderberry"}
	for i, w := range words {
		m.Put(w, i)
	}

	got := m.Keys()
	want := make([]string, len(words))
	copy(want, words)
	sort.Strings(want)
	assert.Equal(t, want, got)
	checkRB(t, m)
}

func Test_TreeMap_Float64Keys(t *testing.T) {
	m := NewTreeMap[float64, string]()
	m.Put(3.14, "pi")
	m.Put(2.71, "e")
	m.Put(1.41, "sqrt2")

	assert.Equal(t, []float64{1.41, 2.71, 3.14}, m.Keys())
	checkRB(t, m)
}

func Test_TreeMap_Uint64Keys(t *testing.T) {
	m := NewTreeMap[uint64, bool]()
	keys := []uint64{100, 50, 200, 75, 150}
	for _, k := range keys {
		m.Put(k, true)
	}

	got := m.Keys()
	assert.Equal(t, []uint64{50, 75, 100, 150, 200}, got)
	checkRB(t, m)
}
