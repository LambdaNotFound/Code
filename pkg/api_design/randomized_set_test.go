package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RandomizedSet(t *testing.T) {
	t.Run("insert_returns_true_for_new_value", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		assert.True(t, rs.Insert(1))
	})

	t.Run("insert_returns_false_for_duplicate", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		rs.Insert(1)
		assert.False(t, rs.Insert(1))
	})

	t.Run("remove_returns_true_for_existing_value", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		rs.Insert(1)
		assert.True(t, rs.Remove(1))
	})

	t.Run("remove_returns_false_for_missing_value", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		assert.False(t, rs.Remove(99))
	})

	t.Run("insert_after_remove_succeeds", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		rs.Insert(1)
		rs.Remove(1)
		assert.True(t, rs.Insert(1))
	})

	t.Run("get_random_returns_only_element", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		rs.Insert(42)
		// with a single element, GetRandom must always return it
		for i := 0; i < 10; i++ {
			assert.Equal(t, 42, rs.GetRandom())
		}
	})

	t.Run("get_random_returns_element_from_set", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		rs.Insert(1)
		rs.Insert(2)
		rs.Insert(3)
		allowed := map[int]bool{1: true, 2: true, 3: true}
		for i := 0; i < 20; i++ {
			assert.True(t, allowed[rs.GetRandom()])
		}
	})

	t.Run("remove_last_element_then_reinsert", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		rs.Insert(1)
		rs.Insert(2)
		rs.Remove(2) // removes last element (swap with itself)
		assert.True(t, rs.Insert(2))
	})

	t.Run("leetcode_example", func(t *testing.T) {
		rs := ConstructorRandomizedSet()
		assert.True(t, rs.Insert(1))
		assert.False(t, rs.Remove(2))
		assert.True(t, rs.Insert(2))
		assert.Contains(t, []int{1, 2}, rs.GetRandom())
		assert.True(t, rs.Remove(1))
		assert.False(t, rs.Insert(2))
		assert.Equal(t, 2, rs.GetRandom())
	})
}
