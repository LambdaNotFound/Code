package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_slice_nil(t *testing.T) {
    var nilSlice []int = nil
    var emptySlice []int = []int{}

    array2d := make([][]int, 2)

    assert.Equal(t, nilSlice, array2d[0])
    assert.NotEqual(t, emptySlice, array2d[0])
    assert.ElementsMatch(t, nilSlice, emptySlice)
}

func Test_slice_basic(t *testing.T) {
    slice := []int{ 1, 2, 3 }
    count := 0

    for _, val := range slice { // i < len(slice), infinite loop
        count += 1

        slice = append(slice, val)
    }

    assert.Equal(t, 3, count)
}
