package binarysearch

import (
	"math"
	"sort"
)

/**
 * sort.Search(n, f) requires the predicate to go [false ... false ... true ... true]
 * sort.Search(n, f) returns the smallest index i in [0, n] at which f(i) is true.
 *
 * All false → returns n (no index in [0, n) satisfied it)
 * All true → returns 0 (the first index is already true)
 */

/*
 * You're building a simplified border security solution that places detection towers at various locations along a border.
 *
 * Different sensors are available for the towers that can detect up to a fixed range.
 *
 * Given the positions of the towers and know border crossing locations, return the minimum range that a sensor would need to be able to detect to cover all border crossing locations.
 *
 * As an example, imagine that your inputs are the following:
 *
 * 1 2 3 4 6 10 12 - Border Crossings
 * 1     4 6 - Tower
 *
 * Output: 6
 */

func minSensorRange(crossings, towers []int) int {
	sort.Ints(towers)
	numTowers, maxMinDist := len(towers), math.MinInt32
	for _, crossing := range crossings {
		rightIdx := sort.Search(numTowers, func(i int) bool { return towers[i] > crossing })

		nearestTowerDist := math.MaxInt32
		if rightIdx < numTowers {
			nearestTowerDist = min(nearestTowerDist, towers[rightIdx]-crossing)
		}
		if rightIdx > 0 {
			nearestTowerDist = min(nearestTowerDist, crossing-towers[rightIdx-1])
		}
		maxMinDist = max(maxMinDist, nearestTowerDist)
	}

	return maxMinDist
}
