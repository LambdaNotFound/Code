package heap

import (
	"container/heap"
	"sort"
)

/**
 * 1942. The Number of the Smallest Unoccupied Chair
 */
func smallestChair(times [][]int, targetFriend int) int {
	// Record before sorting — sort destroys the original index mapping.
	targetArrival := times[targetFriend][0]

	sort.Slice(times, func(i, j int) bool {
		return times[i][0] < times[j][0]
	})

	type occupied struct {
		leaveAt int
		seatNum int
	}

	// Who's sitting where, ordered by earliest departure.
	occupiedHeap := &Heap[occupied]{
		less: func(a, b occupied) bool { return a.leaveAt < b.leaveAt },
	}
	// Freed seats waiting to be reused, ordered by smallest number.
	availableHeap := &Heap[int]{
		less: func(a, b int) bool { return a < b },
	}

	nextSeat := 0

	for _, t := range times {
		arrival, leaving := t[0], t[1]

		// Free every seat whose occupant left by this arrival time.
		for occupiedHeap.Len() > 0 && occupiedHeap.Peek().leaveAt <= arrival {
			freed := heap.Pop(occupiedHeap).(occupied)
			heap.Push(availableHeap, freed.seatNum)
		}

		// Assign the smallest available seat, or mint a new one.
		var seatNum int
		if availableHeap.Len() > 0 {
			seatNum = heap.Pop(availableHeap).(int)
		} else {
			seatNum = nextSeat
			nextSeat++
		}

		heap.Push(occupiedHeap, occupied{leaving, seatNum})

		if arrival == targetArrival {
			return seatNum
		}
	}

	return -1 // unreachable if targetFriend is valid
}
