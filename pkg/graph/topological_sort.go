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
