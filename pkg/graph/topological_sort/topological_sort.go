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

/**
 * 310. Minimum Height Trees | NOT classic topo sort, Tree center = optimal root
 *
 * Given a tree of n nodes labelled from 0 to n - 1,
 * Return a list of all MHTs' root labels.
 *
 * The height of a rooted tree is the number of edges on the
 * longest downward path between the root and a leaf.
 *
 */
func findMinHeightTrees(n int, edges [][]int) []int {
	if n == 1 {
		return []int{0}
	}

	graph := make([][]int, n)
	degree := make([]int, n) // ← separate degree tracking
	for _, e := range edges {
		graph[e[0]] = append(graph[e[0]], e[1])
		graph[e[1]] = append(graph[e[1]], e[0])
		degree[e[0]]++
		degree[e[1]]++
	}

	queue := []int{}
	for i := 0; i < n; i++ {
		if degree[i] == 1 {
			queue = append(queue, i)
		}
	}

	for len(queue) < n { // len(leaves) == n, every remaining node is a leaf
		n -= len(queue)
		nextQueue := []int{}
		for _, leaf := range queue {
			for _, neighbor := range graph[leaf] {
				degree[neighbor]-- // ← decrement instead of mutate
				if degree[neighbor] == 1 {
					nextQueue = append(nextQueue, neighbor)
				}
			}
		}
		queue = nextQueue
	}

	return queue
}

func findMinHeightTreesTwoPassBFS(n int, edges [][]int) []int {
	if n == 1 {
		return []int{0}
	}

	// build adjacency list
	graph := make([][]int, n)
	for _, e := range edges {
		graph[e[0]] = append(graph[e[0]], e[1])
		graph[e[1]] = append(graph[e[1]], e[0])
	}

	// BFS from a given start, returns (farthest node, parent map)
	bfs := func(start int) (int, []int) {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}

		dist[start] = 0
		queue := []int{start}
		farthest := start

		for len(queue) > 0 {
			node := queue[0]
			queue = queue[1:]
			if dist[node] > dist[farthest] {
				farthest = node
			}
			for _, nei := range graph[node] {
				if dist[nei] == -1 {
					dist[nei] = dist[node] + 1
					parent[nei] = node
					queue = append(queue, nei)
				}
			}
		}
		return farthest, parent
	}

	// pass 1: find one endpoint u of the diameter
	u, _ := bfs(0)

	// pass 2: find other endpoint v + track path
	v, parent := bfs(u)

	// reconstruct diameter path from v back to u
	path := []int{}
	for node := v; node != -1; node = parent[node] {
		path = append(path, node)
	}

	// centers are middle node(s) of diameter path
	m := len(path)
	if m%2 == 1 {
		return []int{path[m/2]}
	}
	return []int{path[m/2-1], path[m/2]}
}
