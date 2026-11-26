package heap

/**
 * 973. K Closest Points to Origin
 */
import "container/heap"

func kClosest(points [][]int, k int) [][]int {
	maxHeap := &Heap[[2]int]{
		items: make([][2]int, 0, k),
		less:  func(a, b [2]int) bool { return distance(a) > distance(b) },
	}

	for _, point := range points {
		p := [2]int{point[0], point[1]}
		if maxHeap.Len() == k {
			// only replace if new point is closer than current farthest
			if distance(p) < distance(maxHeap.items[0]) {
				heap.Pop(maxHeap)
				heap.Push(maxHeap, p)
			}
		} else {
			heap.Push(maxHeap, p)
		}
	}

	res := make([][]int, 0, k)
	for maxHeap.Len() > 0 {
		p := heap.Pop(maxHeap).([2]int)
		res = append(res, []int{p[0], p[1]})
	}

	return res
}

func distance(point [2]int) int { return point[0]*point[0] + point[1]*point[1] }
