package apidesign

/*
 * 362. Design Hit Counter
 *
 * Design a hit counter that records timestamps of incoming hits and efficiently queries
 * the number of hits within the past 5 minutes (300 seconds) using a sliding window.
 */

type HitCounter struct {
	queue []int
}

func ConstructorHitCounter() HitCounter {
	return HitCounter{queue: make([]int, 0)}
}

func (hc *HitCounter) Hit(timestamp int) {
	hc.queue = append(hc.queue, timestamp)
}

func (hc *HitCounter) GetHits(timestamp int) int {
	// evict hits outside 300 second window from front
	for len(hc.queue) > 0 && hc.queue[0] <= timestamp-300 {
		hc.queue = hc.queue[1:]
	}
	return len(hc.queue)
}

/*
 * Circular Buffer (alternative implementation — commented out to avoid redeclaration)
 *
 * type HitCounter struct {
 * 	timestamps [300]int
 * 	counts     [300]int
 * }
 *
 * func (hc *HitCounter) Hit(timestamp int) {
 * 	idx := timestamp % 300
 * 	if hc.timestamps[idx] == timestamp {
 * 		hc.counts[idx]++
 * 	} else {
 * 		hc.timestamps[idx] = timestamp
 * 		hc.counts[idx] = 1
 * 	}
 * }
 *
 * func (hc *HitCounter) GetHits(timestamp int) int {
 * 	total := 0
 * 	for i := 0; i < 300; i++ {
 * 		if hc.timestamps[i] > timestamp-300 {
 * 			total += hc.counts[i]
 * 		}
 * 	}
 * 	return total
 * }
 */

/*
 * Binary Search (alternative implementation — commented out to avoid redeclaration)
 *
 * type HitCounter struct {
 * 	timestamps []int
 * }
 *
 * func (hc *HitCounter) Hit(timestamp int) {
 * 	hc.timestamps = append(hc.timestamps, timestamp)
 * }
 *
 * func (hc *HitCounter) GetHits(timestamp int) int {
 * 	cutoff := timestamp - 300
 * 	left, right := 0, len(hc.timestamps)
 * 	for left < right {
 * 		mid := left + (right-left)/2
 * 		if hc.timestamps[mid] <= cutoff {
 * 			left = mid + 1
 * 		} else {
 * 			right = mid
 * 		}
 * 	}
 * 	return len(hc.timestamps) - left
 * }
 */
