package containers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Min/Max heap ordering
// ---------------------------------------------------------------------------

func Test_MinHeap_Ordering(t *testing.T) {
	h := NewMinHeap[int]()
	for _, v := range []int{5, 1, 8, 3, 2} {
		h.Push(v)
	}
	assert.Equal(t, 5, h.Size())
	prev := -1
	for !h.IsEmpty() {
		v, ok := h.Pop()
		assert.True(t, ok)
		assert.GreaterOrEqual(t, v, prev, "expected ascending order")
		prev = v
	}
	assert.Equal(t, 0, h.Size())
}

func Test_MaxHeap_Ordering(t *testing.T) {
	h := NewMaxHeap[string]()
	h.Push("2025/04/20")
	h.Push("2025/04/22")
	h.Push("2025/04/21")

	v, ok := h.Pop()
	assert.True(t, ok)
	assert.Equal(t, "2025/04/22", v)

	v, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, "2025/04/21", v)

	v, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, "2025/04/20", v)
}

// ---------------------------------------------------------------------------
// Custom comparator — struct elements
// ---------------------------------------------------------------------------

func Test_CustomComparator_Struct(t *testing.T) {
	type item struct {
		priority int
		val      string
	}
	h := NewHeap[item](func(a, b item) bool { return a.priority < b.priority })
	h.Push(item{3, "c"})
	h.Push(item{1, "a"})
	h.Push(item{2, "b"})

	for wantPriority := 1; wantPriority <= 3; wantPriority++ {
		got, ok := h.Pop()
		assert.True(t, ok)
		assert.Equal(t, wantPriority, got.priority)
	}
}

// ---------------------------------------------------------------------------
// Empty heap — never panics
// ---------------------------------------------------------------------------

func Test_Heap_EmptyPop(t *testing.T) {
	h := NewMinHeap[int]()
	val, ok := h.Pop()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

func Test_Heap_EmptyPeek(t *testing.T) {
	h := NewMinHeap[int]()
	val, ok := h.Peek()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

// ---------------------------------------------------------------------------
// Single element
// ---------------------------------------------------------------------------

func Test_Heap_SingleElement(t *testing.T) {
	h := NewMinHeap[int]()
	h.Push(42)
	assert.Equal(t, 1, h.Size())
	assert.False(t, h.IsEmpty())

	top, ok := h.Peek()
	assert.True(t, ok)
	assert.Equal(t, 42, top)
	assert.Equal(t, 1, h.Size()) // Peek must not remove

	val, ok := h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 42, val)
	assert.Equal(t, 0, h.Size())
	assert.True(t, h.IsEmpty())
}

// ---------------------------------------------------------------------------
// Size and IsEmpty invariants
// ---------------------------------------------------------------------------

func Test_Heap_SizeAndIsEmpty(t *testing.T) {
	h := NewMinHeap[int]()
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Size())

	h.Push(10)
	h.Push(20)
	h.Push(30)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 3, h.Size())

	h.Pop()
	assert.Equal(t, 2, h.Size())
}

// ---------------------------------------------------------------------------
// Heapify
// ---------------------------------------------------------------------------

func Test_Heapify_Basic(t *testing.T) {
	h := NewMinHeap[int]()
	h.Heapify([]int{9, 4, 7, 1, 5})
	assert.Equal(t, 5, h.Size())

	expected := []int{1, 4, 5, 7, 9}
	for _, want := range expected {
		got, ok := h.Pop()
		assert.True(t, ok)
		assert.Equal(t, want, got)
	}
}

func Test_Heapify_AlreadySorted(t *testing.T) {
	h := NewMinHeap[int]()
	h.Heapify([]int{1, 2, 3, 4, 5})
	prev := -1
	for !h.IsEmpty() {
		v, _ := h.Pop()
		assert.GreaterOrEqual(t, v, prev)
		prev = v
	}
}

func Test_Heapify_ReverseSorted(t *testing.T) {
	h := NewMinHeap[int]()
	h.Heapify([]int{5, 4, 3, 2, 1})
	prev := -1
	for !h.IsEmpty() {
		v, _ := h.Pop()
		assert.GreaterOrEqual(t, v, prev)
		prev = v
	}
}

func Test_PushAfterHeapify(t *testing.T) {
	h := NewMinHeap[int]()
	h.Heapify([]int{5, 3, 8})
	h.Push(1) // new minimum
	assert.Equal(t, 4, h.Size())

	top, ok := h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 1, top)
}

func Test_Heapify_DoesNotMutateInput(t *testing.T) {
	input := []int{3, 1, 2}
	h := NewMinHeap[int]()
	h.Heapify(input)
	assert.Equal(t, []int{3, 1, 2}, input, "caller's slice must be unchanged")
}

// ---------------------------------------------------------------------------
// Pop until empty
// ---------------------------------------------------------------------------

func Test_Heap_PopUntilEmpty(t *testing.T) {
	h := NewMaxHeap[int]()
	for _, v := range []int{3, 1, 4, 1, 5} {
		h.Push(v)
	}
	prev := int(^uint(0) >> 1) // MaxInt
	for {
		v, ok := h.Pop()
		if !ok {
			break
		}
		assert.LessOrEqual(t, v, prev, "expected descending order")
		prev = v
	}
	assert.Equal(t, 0, h.Size())
}

// ---------------------------------------------------------------------------
// Duplicates
// ---------------------------------------------------------------------------

func Test_Heap_Duplicates(t *testing.T) {
	h := NewMinHeap[int]()
	h.Push(5)
	h.Push(5)
	h.Push(5)
	assert.Equal(t, 3, h.Size())

	for i := 0; i < 3; i++ {
		v, ok := h.Pop()
		assert.True(t, ok)
		assert.Equal(t, 5, v)
	}
	assert.Equal(t, 0, h.Size())
}
