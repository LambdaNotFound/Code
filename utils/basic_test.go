package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * golang patterns
 */

// named return values
func splitSum(a, b int) (x int, y int) {
    x = a + b
    y = a - b
    return // returns x and y
}

func Test_slice_nil(t *testing.T) {
    var nilSlice []int = nil
    var emptySlice []int = []int{}

    array2d := make([][]int, 2)

    assert.Equal(t, nilSlice, array2d[0])
    assert.NotEqual(t, emptySlice, array2d[0])
    assert.ElementsMatch(t, nilSlice, emptySlice)

    var sliceA = []int{1, 2, 3}
    var sliceB = sliceA[:0]
    assert.Equal(t, emptySlice, sliceB)
}

func Test_slice_range(t *testing.T) {
    slice := []int{1, 2, 3}
    count := 0

    for _, val := range slice { // i < len(slice) leads to infinite loop
        count += 1

        slice = append(slice, val)
    }

    assert.Equal(t, 3, count)

    sliceIndexed := slice[len(slice)+1:]
    assert.Equal(t, 0, len(sliceIndexed))
    assert.Equal(t, []int{}, sliceIndexed)

    sliceIndexed = slice[:0]
    assert.Equal(t, 0, len(sliceIndexed))
    assert.Equal(t, []int{}, sliceIndexed)
}

func Test_slice_remove_item(t *testing.T) {
    slice := []int{1, 2, 3}
    target := 2
    expected := []int{1, 3}

    for i, _ := range slice { // remove item from a slice
        if slice[i] == target {
            slice[i] = slice[len(slice)-1]
            break
        }
    }
    slice = slice[:len(slice)-1]

    assert.Equal(t, slice, expected)

    slice = []int{1, 3, 5, 8}
    expected = []int{1, 3, 5}
    slice = RemoveIf(slice, func(x int) bool { return x%2 == 0 })
    assert.Equal(t, slice, expected)
}

func Test_slice_spread(t *testing.T) {
    slice := []int{3, 2, 1, 0}
    expected := []int{3, 2, 1, 0}

    sliceNew := append([]int{}, slice...)
    assert.Equal(t, expected, sliceNew)
}

func Test_map_key_value_lookup(t *testing.T) {
    byteToIntMap := make(map[byte]int)

    byteToIntMap['a'] = 1

    assert.Equal(t, 1, byteToIntMap['a'])
    assert.Equal(t, 0, byteToIntMap['b'])
}

func Test_map_key_delete(t *testing.T) {
    byteToIntMap := make(map[byte]int)

    byteToIntMap['a'] = 1
    byteToIntMap['b'] = 2
    byteToIntMap['c'] = 3

    assert.Equal(t, 3, len(byteToIntMap))

    delete(byteToIntMap, 'a')
    assert.Equal(t, 2, len(byteToIntMap))
    assert.Equal(t, 0, byteToIntMap['a'])
}
