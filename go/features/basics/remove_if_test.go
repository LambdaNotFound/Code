package basics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RemoveIf(t *testing.T) {
	t.Run("remove_none", func(t *testing.T) {
		result := RemoveIf([]int{1, 2, 3}, func(x int) bool { return false })
		assert.Equal(t, []int{1, 2, 3}, result)
	})
	t.Run("remove_all", func(t *testing.T) {
		result := RemoveIf([]int{1, 2, 3}, func(x int) bool { return true })
		assert.Empty(t, result)
	})
	t.Run("remove_first", func(t *testing.T) {
		result := RemoveIf([]int{1, 2, 3}, func(x int) bool { return x == 1 })
		assert.Equal(t, []int{2, 3}, result)
	})
	t.Run("remove_last", func(t *testing.T) {
		result := RemoveIf([]int{1, 2, 3}, func(x int) bool { return x == 3 })
		assert.Equal(t, []int{1, 2}, result)
	})
	t.Run("remove_evens", func(t *testing.T) {
		result := RemoveIf([]int{1, 2, 3, 4, 5}, func(x int) bool { return x%2 == 0 })
		assert.Equal(t, []int{1, 3, 5}, result)
	})
	t.Run("empty_input", func(t *testing.T) {
		result := RemoveIf([]int{}, func(x int) bool { return true })
		assert.Empty(t, result)
	})
	t.Run("string_type_param", func(t *testing.T) {
		result := RemoveIf([]string{"a", "bb", "ccc", "d"}, func(x string) bool { return len(x) > 1 })
		assert.Equal(t, []string{"a", "d"}, result)
	})
}
