package apidesign

/*
 * 362. Design Hit Counter
 *
 * Design a hit counter that records timestamps of incoming hits and efficiently queries
 * the number of hits within the past 5 minutes (300 seconds) using a sliding window.
 *
 * 0. Hashmap + queue
 * 1. circular buffer
 *
 * follow up: accessed by multi-thread
 *
 */

type HitCounter struct {
	queue  []int       // just timestamps, no need for a pair type
	counts map[int]int // count each timestamp bucket
	total  int
}

func NewHitCounter() HitCounter {
	return HitCounter{
		queue:  make([]int, 0),
		counts: make(map[int]int),
	}
}

// Time: O(1)
// Space: O(W) W = number of distinct timestamps
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
 * Circular Buffer (alternative implementation)
 *
 * Time: O(1)
 * Space: O(1)
 */
type CircularBufferHitCounter struct {
	timestamps [300]int
	counts     [300]int
}

func NewCircularBufferHitCounter() CircularBufferHitCounter {
	return CircularBufferHitCounter{}
}

func (hc *CircularBufferHitCounter) Hit(timestamp int) {
	idx := timestamp % 300
	if hc.timestamps[idx] == timestamp {
		hc.counts[idx]++
	} else {
		hc.timestamps[idx] = timestamp
		hc.counts[idx] = 1
	}
}

func (hc *CircularBufferHitCounter) GetHits(timestamp int) int {
	total := 0
	for i := 0; i < 300; i++ {
		if hc.timestamps[i] > timestamp-300 {
			total += hc.counts[i]
		}
	}
	return total
}
