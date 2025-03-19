package divide_and_conquer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_partition_asc(t *testing.T) {
	array := []int{7, 3, 4, 6, 5, 5}

	want := 5
	array_result := []int{5, 3, 4, 6, 5, 7}

	got := partition_asc(array, 0, len(array)-1)

	assert.Equal(t, want, got)
	assert.Equal(t, array_result, array)
}

func Test_partition_dec(t *testing.T) {
	array := []int{7, 3, 4, 6, 5, 5}

	want := 0
	array_result := []int{7, 3, 4, 6, 5, 5}

	got := partition_dec(array, 0, len(array)-1)

	assert.Equal(t, want, got)
	assert.Equal(t, array_result, array)
}
