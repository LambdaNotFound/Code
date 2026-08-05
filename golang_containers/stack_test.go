package containers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Stack(t *testing.T) {
    stack := Stack[int]{}
    isEmpty := stack.IsEmpty()
    assert.Equal(t, true, isEmpty)

    stack.Push(1)
    isEmpty = stack.IsEmpty()
    assert.Equal(t, false, isEmpty)

    stack.Push(2)
    stack.Push(3)
    top := stack.Top()
    assert.Equal(t, 3, top)

    value, ok := stack.Pop()
    assert.Equal(t, 3, value)
    assert.Equal(t, true, ok)

    value = stack.Top()
    assert.Equal(t, 2, value)

    stack.Pop()
    value, ok = stack.Pop()
    assert.Equal(t, 1, value)
    assert.Equal(t, true, ok)

    value, ok = stack.Pop()
    assert.Equal(t, 0, value)
    assert.Equal(t, false, ok)
}

func Test_Stack_Peek(t *testing.T) {
    stack := Stack[int]{}

    // Peek on empty stack returns zero value and false
    val, ok := stack.Peek()
    assert.Equal(t, 0, val)
    assert.Equal(t, false, ok)

    // Peek returns top element without removing it
    stack.Push(1)
    stack.Push(2)
    val, ok = stack.Peek()
    assert.Equal(t, 2, val)
    assert.Equal(t, true, ok)
    assert.Equal(t, false, stack.IsEmpty()) // stack unchanged after Peek
}

func Test_MyQueue(t *testing.T) {
    q := Constructor()
    assert.Equal(t, true, q.Empty())

    // Push elements
    q.Push(1)
    q.Push(2)
    q.Push(3)
    assert.Equal(t, false, q.Empty())

    // Peek returns front without removing
    assert.Equal(t, 1, q.Peek())
    assert.Equal(t, false, q.Empty())

    // Pop returns elements in FIFO order
    assert.Equal(t, 1, q.Pop())
    assert.Equal(t, 2, q.Pop())

    // Push after partial drain (tests interleaved push/pop with lazy stack1->stack2 transfer)
    q.Push(4)
    assert.Equal(t, 3, q.Pop())
    assert.Equal(t, 4, q.Pop())
    assert.Equal(t, true, q.Empty())
}
