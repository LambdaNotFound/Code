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

/**
 * 475. Heaters
 *
 * Time: O((H + N) log N)
 *     Sorting heaters: O(N log N) where N = len(heaters)
 *     For each of H = len(houses) houses, sort.Search does a binary search over N heaters: O(H log N)
 *
 * Space: O(log N)
 *     Sorting is in-place but uses O(log N) stack space for Go's sort.Ints (introsort)
 */
func findRadius(houses []int, heaters []int) int {
	sort.Ints(heaters)

	maxMinDist, n := math.MinInt32, len(heaters)
	for _, house := range houses {
		rightIdx := sort.Search(n, func(i int) bool {
			return house < heaters[i]
		})
		leftIdx := rightIdx - 1

		minDist := math.MaxInt32
		if rightIdx < n {
			minDist = min(minDist, heaters[rightIdx]-house)
		}
		if leftIdx >= 0 {
			minDist = min(minDist, house-heaters[leftIdx])
		}
		maxMinDist = max(maxMinDist, minDist)
	}

	return maxMinDist
}
