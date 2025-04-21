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
