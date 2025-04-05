package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_slice_basics(t *testing.T) {
    var nilSlice []int = nil
    var emptySlice []int = []int{}

    array2d := make([][]int, 2)

    assert.Equal(t, nilSlice, array2d[0])
    assert.NotEqual(t, emptySlice, array2d[0])
    assert.ElementsMatch(t, nilSlice, emptySlice)
}
