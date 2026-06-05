package heap

/**
 * 973. K Closest Points to Origin
 */
import "container/heap"

func kClosest(points [][]int, k int) [][]int {
	maxHeap := &Heap[[]int]{
		items: make([][]int, 0, k),
		less:  func(a, b []int) bool { return distance(a) > distance(b) },
	}

	for _, point := range points {
		if maxHeap.Len() == k {
			// only replace if new point is closer than current farthest
			if distance(point) < distance(maxHeap.items[0]) {
				heap.Pop(maxHeap)
				heap.Push(maxHeap, point)
			}
		} else {
			heap.Push(maxHeap, point)
		}
	}

	res := make([][]int, 0, k)
	for maxHeap.Len() > 0 {
		res = append(res, heap.Pop(maxHeap).([]int))
	}

	return res
}

func distance(point []int) int { return point[0]*point[0] + point[1]*point[1] }
