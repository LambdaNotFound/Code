package divide_and_conquer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_quick_sort(t *testing.T) {
    array := []int{5, 6, 7, 2, 1, 0}

    expected := []int{0, 1, 2, 5, 6, 7}

    quick_sort(array, 0, len(array)-1)
    assert.Equal(t, expected, array)
}

func Test_partition(t *testing.T) {
    array := []int{7, 3, 4, 6, 5, 5}

    want := 2
    expected := []int{3, 4, 5, 6, 5, 7}

    got := partition(array, 0, len(array)-1)

    assert.Equal(t, want, got)
    assert.Equal(t, expected, array)
}

func Test_partition_asc(t *testing.T) {
    array := []int{7, 3, 4, 6, 5, 5}

    want := 5
    expected := []int{5, 3, 4, 6, 5, 7}

    got := partition_asc(array, 0, len(array)-1)

    assert.Equal(t, want, got)
    assert.Equal(t, expected, array)
}

func Test_partition_decending(t *testing.T) {
    array := []int{7, 3, 4, 6, 5, 5}
    want := 2
    expected := []int{7, 6, 5, 3, 5, 4}

    got := partition_decending(array, 0, len(array)-1)
    
    assert.Equal(t, want, got)
    assert.Equal(t, expected, array)
}

func Test_partition_dec(t *testing.T) {
    array := []int{7, 3, 4, 6, 5, 5}

    want := 0
    expected := []int{7, 3, 4, 6, 5, 5}

    got := partition_dec(array, 0, len(array)-1)

    assert.Equal(t, want, got)
    assert.Equal(t, expected, array)
}
