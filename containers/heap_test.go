package containers

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MinHeap(t *testing.T) {
    minHeap := &Heap[int]{
        less: func(a, b int) bool { return a < b },
    }
    heap.Init(minHeap)

    minHeap.PushItem(5)
    minHeap.PushItem(3)
    minHeap.PushItem(8)
    value := minHeap.Peek()
    assert.Equal(t, 3, value)

    value = minHeap.PopItem()
    assert.Equal(t, 3, value)

    value = minHeap.Peek()
    assert.Equal(t, 5, value)
}

func Test_MaxHeap(t *testing.T) {
    minHeap := &Heap[string]{
        less: func(a, b string) bool { return a > b },
    }
    heap.Init(minHeap)

    minHeap.PushItem("2025/04/20")
    minHeap.PushItem("2025/04/21")
    minHeap.PushItem("2025/04/22")
    value := minHeap.Peek()
    assert.Equal(t, "2025/04/22", value)

    value = minHeap.PopItem()
    assert.Equal(t, "2025/04/22", value)

    value = minHeap.Peek()
    assert.Equal(t, "2025/04/21", value)
}

func Test_Heap_SingleElement(t *testing.T) {
    h := &Heap[int]{less: func(a, b int) bool { return a < b }}
    heap.Init(h)

    h.PushItem(42)
    assert.Equal(t, 1, h.Len())
    assert.Equal(t, 42, h.Peek())

    val := h.PopItem()
    assert.Equal(t, 42, val)
    assert.Equal(t, 0, h.Len())
}

func Test_Heap_Duplicates(t *testing.T) {
    h := &Heap[int]{less: func(a, b int) bool { return a < b }}
    heap.Init(h)

    h.PushItem(5)
    h.PushItem(5)
    h.PushItem(5)
    assert.Equal(t, 3, h.Len())
    assert.Equal(t, 5, h.PopItem())
    assert.Equal(t, 5, h.PopItem())
    assert.Equal(t, 5, h.PopItem())
    assert.Equal(t, 0, h.Len())
}
