package topologicalsort

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
 *
 * Time: O(V + E)
 *     Building indegree and adjList: O(E)
 *     Seeding the queue: O(V)
 *     BFS: each course dequeued once O(V), each edge traversed once O(E)
 *
 * Space: O(V + E)
 *     adjList: O(V + E) — V sublists with E total entries
 *     indegree, queue, order: O(V) each
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
