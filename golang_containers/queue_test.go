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

func Test_Queue_Peek(t *testing.T) {
    queue := Queue[int]{}

    // Peek on empty queue returns zero value and false
    val, ok := queue.Peek()
    assert.Equal(t, 0, val)
    assert.Equal(t, false, ok)

    // Peek returns front element without removing it
    queue.Enqueue(10)
    queue.Enqueue(20)
    val, ok = queue.Peek()
    assert.Equal(t, 10, val)
    assert.Equal(t, true, ok)
    assert.Equal(t, 2, queue.Size()) // size unchanged after Peek
}
