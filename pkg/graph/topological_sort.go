package graph

/**
 * Adjacency List + Topological Sort
 *
 * in-degree map + adjacent map + BFS search
 *
 *    indegree := make(map[int][]int)
 *    adjacent := make(map[int][]int)
 *    count := 0
 *
 *    for _,  prerequisite := range  prerequisites {
 *        src, dst :=  prerequisite[1], prerequisite[0]
 *
 *        indegree[dst] = append(indegree[dst], src)
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
    // Store the indegree of each course
    indegree := make([]int, numCourses)
    for _, pre := range prerequisites {
        indegree[pre[0]]++
    }

    // Store the courses that can be completed (indegree == 0) in queue
    var queue []int
    for course, degree := range indegree {
        if degree == 0 {
            queue = append(queue, course)
        }
    }

    // These are the number of courses that can be completed
    completeCount := len(queue)
    for len(queue) != 0 {
        leaf := queue[0]
        queue = queue[1:]

        for _, pre := range prerequisites {
            // If course is the prerequite of any other course?
            if leaf == pre[1] {
                // If yes, then reduce the indegree of that course
                indegree[pre[0]]--
                // Is there a cycle?
                if indegree[pre[0]] < 0 {
                    return false
                }
                // Can this course be completed now?
                if indegree[pre[0]] == 0 {
                    queue = append(queue, pre[0])
                    completeCount++
                }
            }
        }
    }

    return completeCount == numCourses
}

/**
 * 210. Course Schedule II
 *
 * Return the topological ordering
 */
func findOrder(numCourses int, prerequisites [][]int) []int {
    // Store the indegree of each course
    indegree := make([]int, numCourses)
    for _, pre := range prerequisites {
        indegree[pre[0]]++
    }

    // Store the courses that can be completed (indegree == 0) in queue
    var queue []int
    for course, degree := range indegree {
        if degree == 0 {
            queue = append(queue, course)
        }
    }

    // These are the number of courses that can be completed
    completeCount := len(queue)
    res := []int{}
    for len(queue) != 0 {
        leaf := queue[0]
        queue = queue[1:]
        res = append(res, leaf)

        for _, pre := range prerequisites {
            // If course is the prerequite of any other course?
            if leaf == pre[1] {
                // If yes, then reduce the indegree of that course
                indegree[pre[0]]--
                // Is there a cycle?
                if indegree[pre[0]] < 0 {
                    return []int{}
                }
                // Can this course be completed now?
                if indegree[pre[0]] == 0 {
                    queue = append(queue, pre[0])
                    completeCount++
                }
            }
        }
    }

    if completeCount == numCourses {
        return res
    }
    return []int{}
}

/**
 * 630. Course Schedule III
 *
 * You are given an array courses where courses[i] = [durationi, lastDayi]
 * indicate that the ith course should be taken continuously for durationi
 * days and must be finished before or on lastDayi
 */
func scheduleCourse(courses [][]int) int {
    return 0
}

/**
 * 310. Minimum Height Trees
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

    // build the graph: adjacency list + degree count
    graph := map[int][]int{}
    for _, edge := range edges {
        src, dst := edge[0], edge[1]
        graph[src] = append(graph[src], dst)
        graph[dst] = append(graph[dst], src)
    }

    // find all leaves (degree == 1)
    leaves := []int{}
    for k, v := range graph {
        if len(v) == 1 {
            leaves = append(leaves, k)
        }
    }

    for len(leaves) < n {
        n -= len(leaves)

        new_leaves := []int{}
        for _, leaf := range leaves {
            currentLeaf := graph[leaf][0] // remove leaf (degree == 1) from Adjacency List
            for i := 0; i < len(graph[currentLeaf]); i++ {
                if graph[currentLeaf][i] == leaf {
                    graph[currentLeaf] = append(graph[currentLeaf][:i], graph[currentLeaf][i+1:]...) // remove
                    break
                }
            }

            if len(graph[currentLeaf]) == 1 { // add new leaf
                new_leaves = append(new_leaves, currentLeaf)
            }
        }

        leaves = new_leaves
    }

    return leaves
}

/**
 * 329. Longest Increasing Path in a Matrix
 */
func longestIncreasingPath(matrix [][]int) int {
    m := len(matrix)
    n := len(matrix[0])

    indegree := make([][]int, m)
    for i := 0; i < m; i++ {
        indegree[i] = make([]int, n)
    }

    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            for _, dir := range [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
                nextRow := r + dir[0]
                nextCol := c + dir[1]
                if nextRow < 0 || nextCol < 0 || nextRow == m || nextCol == n || 
                    matrix[nextRow][nextCol] <= matrix[r][c] { // increasing
                    continue
                }
                indegree[r][c]++
            }
        }
    }

    queue := [][]int{}
    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if indegree[r][c] == 0 {
                queue = append(queue, []int{r, c})
            }
        }
    }

    steps := 0
    for len(queue) > 0 {
        l := len(queue)
        for i := 0; i < l; i++ {
            r, c := queue[0][0], queue[0][1]
            queue = queue[1:]

            for _, dir := range [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
                nextRow := r + dir[0]
                nextCol := c + dir[1]
                if nextRow < 0 || nextCol < 0 || nextRow == m || nextCol == n || 
                    matrix[nextRow][nextCol] >= matrix[r][c] {
                    continue
                }
                indegree[nextRow][nextCol]--
                if indegree[nextRow][nextCol] == 0 {
                    queue = append(queue, []int{nextRow, nextCol})
                }
            }
        }
        steps++
    }

    return steps
}
