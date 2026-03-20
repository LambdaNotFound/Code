package heap

/**
 * 973. K Closest Points to Origin
 */
import "container/heap"

func kClosest(points [][]int, k int) [][]int {
	maxHeap := &PointHeap{
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

type PointHeap struct {
	items [][]int
	less  func([]int, []int) bool
}

func (h *PointHeap) Len() int           { return len(h.items) }
func (h *PointHeap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *PointHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *PointHeap) Push(item interface{}) {
	h.items = append(h.items, item.([]int))
}

func (h *PointHeap) Pop() interface{} {
	item := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return item
}
