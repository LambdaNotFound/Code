package topologicalsort

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
 * 207. Course Schedule
 *
 * Return true if you can finish all courses. Otherwise, return false
 */
func canFinish(numCourses int, prerequisites [][]int) bool {
	indegree := make([]int, numCourses)
	adjList := make([][]int, numCourses)
	for _, prerequisite := range prerequisites {
		src, dst := prerequisite[0], prerequisite[1]
		indegree[dst] += 1
		adjList[src] = append(adjList[src], dst)
	}

	queue := []int{}
	for i := range indegree {
		if indegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	res := []int{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		res = append(res, cur)

		for _, v := range adjList[cur] {
			indegree[v] -= 1
			if indegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	return len(res) == numCourses
}

/**
 * 210. Course Schedule II
 *
 * Return the topological ordering
 */
func findOrder(numCourses int, prerequisites [][]int) []int {
	indegree := make([]int, numCourses)
	adjList := make([][]int, numCourses)
	for _, prereq := range prerequisites {
		prereqCourse, course := prereq[1], prereq[0]
		indegree[course]++
		adjList[prereqCourse] = append(adjList[prereqCourse], course)
	}

	queue := []int{}
	for i := range indegree {
		if indegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	order := []int{}
	for len(queue) > 0 {
		course := queue[0]
		queue = queue[1:]
		order = append(order, course)
		for _, nextCourse := range adjList[course] {
			indegree[nextCourse]--
			if indegree[nextCourse] == 0 {
				queue = append(queue, nextCourse)
			}
		}
	}

	if len(order) != numCourses {
		return []int{}
	}

	return order
}

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
