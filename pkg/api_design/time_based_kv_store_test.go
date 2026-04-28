package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TimeMap(t *testing.T) {
    timeMap := ConstructorTimeMap()

    timeMap.Set("foo", "bar", 1)
    val := timeMap.Get("foo", 1)
    assert.Equal(t, "bar", val)

    val = timeMap.Get("foo", 3)
    assert.Equal(t, "bar", val)

    timeMap.Set("foo", "bar2", 4)
    val = timeMap.Get("foo", 4)
    assert.Equal(t, "bar2", val)

    val = timeMap.Get("foo", 5)
    assert.Equal(t, "bar2", val)
}

// Test_GetByUpperBound covers the sort.Search-based upper-bound variant.
//
// Branch coverage:
//   - missing key → ""
//   - timestamp before all entries (index==0) → ""
//   - timestamp exactly equal to an entry → that entry's value
//   - timestamp between two entries → value of the lower entry (floor)
//   - timestamp after all entries → last value
func Test_GetByUpperBound(t *testing.T) {
	t.Run("missing_key_returns_empty", func(t *testing.T) {
		tm := ConstructorTimeMap()
		assert.Equal(t, "", tm.GetByUpperBound("missing", 5))
	})

	t.Run("timestamp_before_all_entries_returns_empty", func(t *testing.T) {
		tm := ConstructorTimeMap()
		tm.Set("k", "v1", 10)
		tm.Set("k", "v2", 20)
		// query timestamp 3 < 10 → index==0 → ""
		assert.Equal(t, "", tm.GetByUpperBound("k", 3))
	})

	t.Run("timestamp_exact_match", func(t *testing.T) {
		tm := ConstructorTimeMap()
		tm.Set("k", "bar", 1)
		tm.Set("k", "bar2", 4)
		assert.Equal(t, "bar", tm.GetByUpperBound("k", 1))
		assert.Equal(t, "bar2", tm.GetByUpperBound("k", 4))
	})

	t.Run("timestamp_between_entries_returns_floor", func(t *testing.T) {
		tm := ConstructorTimeMap()
		tm.Set("k", "high", 10)
		tm.Set("k", "low", 20)
		// query 15: upper-bound lands at index 1 (first timestamp > 15 is 20) → arr[0]="high"
		assert.Equal(t, "high", tm.GetByUpperBound("k", 15))
	})

	t.Run("timestamp_after_all_entries_returns_last", func(t *testing.T) {
		tm := ConstructorTimeMap()
		tm.Set("k", "first", 5)
		tm.Set("k", "last", 10)
		assert.Equal(t, "last", tm.GetByUpperBound("k", 99))
	})

	t.Run("same_behavior_as_get_for_valid_queries", func(t *testing.T) {
		// GetByUpperBound and Get must agree on the floor value for the same inputs.
		tm := ConstructorTimeMap()
		tm.Set("foo", "bar", 1)
		tm.Set("foo", "bar2", 4)
		for _, ts := range []int{1, 3, 4, 5} {
			assert.Equal(t, tm.Get("foo", ts), tm.GetByUpperBound("foo", ts),
				"mismatch at timestamp %d", ts)
		}
	})
}

func Test_TimeMap_2(t *testing.T) {
    timeMap := ConstructorTimeMap()

    timeMap.Set("love", "high", 10)
    timeMap.Set("love", "low", 20)
    val := timeMap.Get("love", 5)
    assert.Equal(t, "", val)

    val = timeMap.Get("love", 10)
    assert.Equal(t, "high", val)

    val = timeMap.Get("love", 15)
    assert.Equal(t, "high", val)

    val = timeMap.Get("love", 20)
    assert.Equal(t, "low", val)

    val = timeMap.Get("love", 25)
    assert.Equal(t, "low", val)
}
