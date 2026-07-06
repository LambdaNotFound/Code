package utils

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * slices.Clone / maps.Clone (Go 1.21+)
 *
 * Both return a shallow copy: a new top-level slice/map backed by new
 * storage, but element values are copied by plain assignment — nested
 * reference types (slices, maps, pointers) inside the elements are still
 * shared with the source.
 */

func Test_slices_Clone_independence(t *testing.T) {
	original := []int{1, 2, 3}
	cloned := slices.Clone(original)

	cloned[0] = 99
	assert.Equal(t, []int{1, 2, 3}, original) // clone mutation doesn't leak back

	original[1] = 42
	assert.Equal(t, 99, cloned[0])
	assert.Equal(t, 2, cloned[1]) // original mutation doesn't leak into clone
}

func Test_slices_Clone_nil(t *testing.T) {
	var nilSlice []int
	cloned := slices.Clone(nilSlice)

	assert.Nil(t, cloned)
}

func Test_slices_Clone_shallow_nested(t *testing.T) {
	original := [][]string{{"a", "b"}, {"c"}}
	cloned := slices.Clone(original)

	cloned[0] = append([]string{}, "z") // reassigning the element is fine, doesn't affect original
	assert.Equal(t, "a", original[0][0])

	cloned[1][0] = "mutated" // mutating a shared inner slice DOES leak back
	assert.Equal(t, "mutated", original[1][0])
}

func Test_maps_Clone_independence(t *testing.T) {
	original := map[string]int{"a": 1, "b": 2}
	cloned := maps.Clone(original)

	cloned["a"] = 99
	assert.Equal(t, 1, original["a"]) // clone mutation doesn't leak back

	original["b"] = 42
	assert.Equal(t, 2, cloned["b"]) // original mutation doesn't leak into clone
}

func Test_maps_Clone_nil(t *testing.T) {
	var nilMap map[string]int
	cloned := maps.Clone(nilMap)

	assert.Nil(t, cloned)
}

func Test_maps_Clone_shallow_nested(t *testing.T) {
	original := map[string][]int{"a": {1, 2}}
	cloned := maps.Clone(original)

	cloned["a"][0] = 99 // mutating the shared inner slice DOES leak back
	assert.Equal(t, 99, original["a"][0])
}
