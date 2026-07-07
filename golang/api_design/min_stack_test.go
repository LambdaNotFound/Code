package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MinStack(t *testing.T) {
    stack := ConstructorMinStack()

    stack.Push(-2)
    stack.Push(0)
    stack.Push(-3)
    val := stack.GetMin()
    assert.Equal(t, -3, val)
    stack.Pop()
    val = stack.Top()
    assert.Equal(t, 0, val)
    val = stack.GetMin()
    assert.Equal(t, -2, val)
}
