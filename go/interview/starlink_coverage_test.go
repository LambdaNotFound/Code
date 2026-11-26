package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stationLeftRight(t *testing.T) {
	s := Station{pos: 100, reach: 30}
	assert.Equal(t, 70.0, s.left())
	assert.Equal(t, 130.0, s.right())
}

func Test_minStationsCover(t *testing.T) {
	t.Run("example_from_driver", func(t *testing.T) {
		stations := []Station{
			{id: 1, pos: 100, reach: 120},
			{id: 2, pos: 300, reach: 150},
			{id: 3, pos: 250, reach: 80},
			{id: 4, pos: 550, reach: 180},
			{id: 5, pos: 800, reach: 220},
			{id: 6, pos: 700, reach: 50},
		}
		count, chosen := minStationsCover(0, 1000, stations)
		assert.Equal(t, 4, count)
		assert.Len(t, chosen, 4)
	})

	t.Run("single_station_covers_pass", func(t *testing.T) {
		stations := []Station{{id: 1, pos: 50, reach: 60}}
		count, chosen := minStationsCover(0, 100, stations)
		assert.Equal(t, 1, count)
		assert.Len(t, chosen, 1)
	})

	t.Run("gap_makes_coverage_impossible", func(t *testing.T) {
		stations := []Station{
			{id: 1, pos: 10, reach: 5},  // covers [5, 15]
			{id: 2, pos: 50, reach: 10}, // covers [40, 60] — gap at [15,40]
		}
		count, chosen := minStationsCover(0, 60, stations)
		assert.Equal(t, -1, count)
		assert.Nil(t, chosen)
	})

	t.Run("exact_boundary_coverage", func(t *testing.T) {
		stations := []Station{
			{id: 1, pos: 5, reach: 5},  // covers [0, 10]
			{id: 2, pos: 15, reach: 5}, // covers [10, 20]
		}
		count, _ := minStationsCover(0, 20, stations)
		assert.Equal(t, 2, count)
	})
}

func Test_minStationsCoverTest(t *testing.T) {
	// covers the demo driver function
	minStationsCoverTest()
}
