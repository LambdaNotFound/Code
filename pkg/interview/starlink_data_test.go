package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lastNonOverlapping(t *testing.T) {
	windows := []Window{
		{start: 0, end: 3},
		{start: 2, end: 5},
		{start: 5, end: 8},
		{start: 6, end: 9},
	}

	t.Run("no_prior_non_overlapping", func(t *testing.T) {
		assert.Equal(t, -1, lastNonOverlapping(windows, 1))
	})

	t.Run("finds_non_overlapping_predecessor", func(t *testing.T) {
		// windows[2].start=5, windows[0].end=3<=5 and windows[1].end=5<=5 → last is index 1
		assert.Equal(t, 1, lastNonOverlapping(windows, 2))
	})

	t.Run("first_window_has_no_predecessor", func(t *testing.T) {
		assert.Equal(t, -1, lastNonOverlapping(windows, 0))
	})
}

func Test_maxData(t *testing.T) {
	t.Run("empty_windows", func(t *testing.T) {
		total, chosen := maxData([]Window{})
		assert.Equal(t, 0, total)
		assert.Nil(t, chosen)
	})

	t.Run("single_window", func(t *testing.T) {
		windows := []Window{{start: 0, end: 5, data: 10, satID: 1}}
		total, chosen := maxData(windows)
		assert.Equal(t, 10, total)
		assert.Len(t, chosen, 1)
	})

	t.Run("non_overlapping_all_selected", func(t *testing.T) {
		windows := []Window{
			{start: 0, end: 3, data: 10, satID: 1},
			{start: 3, end: 6, data: 20, satID: 2},
			{start: 6, end: 9, data: 15, satID: 3},
		}
		total, chosen := maxData(windows)
		assert.Equal(t, 45, total)
		assert.Len(t, chosen, 3)
	})

	t.Run("overlapping_best_combination", func(t *testing.T) {
		windows := []Window{
			{start: 0, end: 3, data: 10, satID: 1},
			{start: 1, end: 4, data: 15, satID: 2},
			{start: 5, end: 8, data: 20, satID: 1},
		}
		// windows[1] (15) and windows[2] (20) don't overlap → total 35
		total, _ := maxData(windows)
		assert.Equal(t, 35, total)
	})

	t.Run("driver_example", func(t *testing.T) {
		maxDataTest() // covers the demo driver
	})
}
