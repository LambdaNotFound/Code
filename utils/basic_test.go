package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_basics(t *testing.T) {
    array2d := make([][]int, 2)
    var nilSlice []int = nil
    var emptySlice []int = []int{}

    assert.Equal(t, nilSlice, array2d[0])
    assert.ElementsMatch(t, nilSlice, emptySlice)
}
