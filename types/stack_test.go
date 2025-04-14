package types

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
