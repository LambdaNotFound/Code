package greedy

import (
	"container/heap"
	"sort"
)

/**
 * Adjacency List + Topological Sort | DAGs and topo sort
 *
 * in-degree map + adjacent map + BFS search
 *
 *    indegree := make(map[int][]int) => make(map[int]int)
 *    adjacent := make(map[int][]int)
 *    count := 0
 *
 *    for _,  prerequisite := range  prerequisites {
 *        src, dst :=  prerequisite[1], prerequisite[0]
 *
 *        indegree[dst] = append(indegree[dst], src) => indegree[dst] += 1
 *        adjacent[src] = append(adjacent[src], dst)
 *    }
 *
 * Time Complexity: O(n)
 * Space Complexity: O(n)
 *
 */

/**
 * 630. Course Schedule III
 *
 * Greedy + max‑heap
 *
 * Sort courses by their end day. Iterate through them, always taking the
 * current course and pushing its duration into a max‑heap. If the total
 * time exceeds the current course's deadline, drop the course with the
 * longest duration (top of the max‑heap). The number of courses remaining
 * in the heap is the maximum count.
 *
 * This mirrors the classic optimal solution and matches `scheduleCourseHeap`
 * so both implementations return the same result for the tests.
 */
func scheduleCourse(courses [][]int) int {
	if len(courses) == 0 {
		return 0
	}

	day, maxHeap := 0, &MaxHeap{}

	// Sort by deadline (earliest first)
	sort.Slice(courses, func(i, j int) bool {
		return courses[i][1] < courses[j][1]
	})

	for _, course := range courses {
		duration, lastDay := course[0], course[1]

		// Take this course
		heap.Push(maxHeap, duration)
		day += duration

		// If we exceed the deadline, drop the longest course taken so far
		if day > lastDay {
			day -= heap.Pop(maxHeap).(int)
		}
	}

	return maxHeap.Len()
}

type MaxHeap []int

func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Pop() (v interface{}) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return v
}
func (h *MaxHeap) Push(v interface{}) { *h = append(*h, v.(int)) }

func scheduleCourseHeap(courses [][]int) int {
	day, maxHeap := 0, &MaxHeap{}
	sort.Slice(courses, func(i, j int) bool {
		return courses[i][1] < courses[j][1]
	})

	for _, course := range courses {
		heap.Push(maxHeap, course[0])
		day += course[0]

		if day > course[1] {
			day -= heap.Pop(maxHeap).(int)
		}

		if day > course[1] {
			return maxHeap.Len()
		}
	}

	return maxHeap.Len()
}
