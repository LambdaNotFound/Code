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
 *
 * 1. DFS + memo
 * 2. BFS + topological sort
 */
func longestIncreasingPath(matrix [][]int) int {
    m := len(matrix)
    n := len(matrix[0])

    inDegree := make([][]int, m)
    for i := 0; i < m; i++ {
        inDegree[i] = make([]int, n)
    }

    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            for _, dir := range [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
                nr := r + dir[0]
                nc := c + dir[1]
                if nr < 0 || nc < 0 || nr == m || nc == n || matrix[nr][nc] <= matrix[r][c] {
                    continue
                }
                inDegree[r][c]++
            }
        }
    }

    queue := [][]int{}
    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if inDegree[r][c] == 0 {
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
                nr := r + dir[0]
                nc := c + dir[1]
                if nr < 0 || nc < 0 || nr == m || nc == n || matrix[nr][nc] >= matrix[r][c] {
                    continue
                }
                inDegree[nr][nc]--
                if inDegree[nr][nc] == 0 {
                    queue = append(queue, []int{nr, nc})
                }
            }
        }
        steps++
    }

    return steps
}

func longestIncreasingPathMemoization(matrix [][]int) int {
    if len(matrix) == 0 || len(matrix[0]) == 0 {
        return 0
    }

    m, n := len(matrix), len(matrix[0])
    memo := make([][]int, m)
    for i := range memo {
        memo[i] = make([]int, n)
    }

    dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

    var dfs func(int, int) int
    dfs = func(x, y int) int {
        if memo[x][y] != 0 {
            return memo[x][y]
        }
        maxLen := 1
        for _, dir := range dirs {
            nx, ny := x+dir[0], y+dir[1]
            if nx >= 0 && nx < m && ny >= 0 && ny < n && matrix[nx][ny] > matrix[x][y] {
                length := 1 + dfs(nx, ny)
                if length > maxLen {
                    maxLen = length
                }
            }
        }
        memo[x][y] = maxLen
        return maxLen
    }

    res := 0
    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            res = max(res, dfs(i, j))
        }
    }
    return res
}
