package containers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Queue(t *testing.T) {
    queue := Queue[int]{}
    isEmpty := queue.IsEmpty()
    assert.Equal(t, true, isEmpty)

    queue.Enqueue(1)
    isEmpty = queue.IsEmpty()
    assert.Equal(t, false, isEmpty)
    size := queue.Size()
    assert.Equal(t, 1, size)

    queue.Enqueue(2)
    queue.Enqueue(3)
    value, ok := queue.Dequeue()
    assert.Equal(t, 1, value)
    assert.Equal(t, true, ok)
    size = queue.Size()
    assert.Equal(t, 2, size)

    queue.Dequeue()
    front := queue.Front()
    assert.Equal(t, 3, front)
    value, ok = queue.Dequeue()
    assert.Equal(t, 3, value)
    assert.Equal(t, true, ok)

    
    value, ok = queue.Dequeue()
    assert.Equal(t, 0, value)
    assert.Equal(t, false, ok)
}
