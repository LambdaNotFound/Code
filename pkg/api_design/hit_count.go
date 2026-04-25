package apidesign

/*
 * 362. Design Hit Counter
 *
 * Design a hit counter that records timestamps of incoming hits and efficiently queries
 * the number of hits within the past 5 minutes (300 seconds) using a sliding window.
 *
 * circular buffer
 */

type HitCounter struct {
	queue  []int // just timestamps, no need for a pair type
	counts map[int]int
	total  int
}

func ConstructorHitCounter() HitCounter {
	return HitCounter{
		queue:  make([]int, 0),
		counts: make(map[int]int),
	}
}

func (hc *HitCounter) Hit(timestamp int) {
	hc.counts[timestamp]++
	hc.total++
	if hc.counts[timestamp] == 1 {
		hc.queue = append(hc.queue, timestamp)
	}
}

func (hc *HitCounter) GetHits(timestamp int) int {
	for len(hc.queue) > 0 && hc.queue[0] <= timestamp-300 {
		hc.total -= hc.counts[hc.queue[0]] // look up live count
		delete(hc.counts, hc.queue[0])     // clean up map entry
		hc.queue = hc.queue[1:]
	}
	return hc.total
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
