package graph

/**
 * Topological Sort
 *
 * in-degree map + adjacent map + BFS search
 *
 *  indegree := make(map[int][]int)
 *  adjacent := make(map[int][]int)
 *  count := 0
 *
 *  for _,  prerequisite := range  prerequisites {
 *      src :=  prerequisite[1]
 *      dst :=  prerequisite[0]
 *
 *      indegree[dst] = append(indegree[dst], src)
 *      adjacent[src] = append(adjacent[src], dst)
 *  }
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
    canFinishCount := len(queue)
    for len(queue) != 0 {
        // This course is completed
        course := queue[0]
        queue = queue[1:]

        for _, pre := range prerequisites {
            // If course is the prerequite of any other course?
            if course == pre[1] {
                // If yes, then reduce the indegree of that course
                indegree[pre[0]]--
                // Is there a cycle?
                if indegree[pre[0]] < 0 {
                    return false
                }
                // Can this course be completed now?
                if indegree[pre[0]] == 0 {
                    canFinishCount++
                    queue = append(queue, pre[0])
                }
            }
        }
    }

    return canFinishCount == numCourses
}

/**
 * 310. Minimum Height Trees
 *
 * Return a list of all MHTs' root labels.
 *
 * The height of a rooted tree is the number of edges on the
 * longest downward path between the root and a leaf.
 *
 * Adjacency List + Topological Sort
 */
func findMinHeightTrees(n int, edges [][]int) []int {
    if n == 1 {
        return []int{0}
    }

    graph := map[int][]int{}
    for _, edge := range edges {
        src, dst := edge[0], edge[1]
        graph[src] = append(graph[src], dst)
        graph[dst] = append(graph[dst], src)
    }

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
            node := graph[leaf][0] // remove leaf from node
            for i := 0; i < len(graph[node]); i++ {
                if graph[node][i] == leaf {
                    graph[node][i] = graph[node][len(graph[node])-1] // copy/move to truncate
                    break
                }
            }
            graph[node] = graph[node][:len(graph[node])-1]

            if len(graph[node]) == 1 {
                new_leaves = append(new_leaves, node)
            }
        }

        leaves = new_leaves
    }

    return leaves
}
