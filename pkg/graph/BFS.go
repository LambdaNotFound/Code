package graph

import (
	. "gocode/containers"
	. "gocode/types"
)

/**
 * 994. Rotting Oranges
 *
 * 0 representing an empty cell,
 * 1 representing a fresh orange, or
 * 2 representing a rotten orange.
 *
 */
func orangesRotting(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    queue := Queue[[2]int]{} // queue of [row, col]
    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if grid[r][c] == 2 {
                queue.Enqueue([2]int{r, c})
            }
        }
    }

    directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

    minutes := 0
    for !queue.IsEmpty() {
        for i := 0; i < queue.Size(); i++ {
            orange, _ := queue.Peek()
            queue.Dequeue()
            for _, dir := range directions {
                r, c := orange[0]+dir[0], orange[1]+dir[1]
                if r < 0 || r >= m || c < 0 || c >= n {
                    continue
                }
                if grid[r][c] != 1 {
                    continue
                }
                grid[r][c] = 2
                queue.Enqueue([2]int{r, c})
            }
        }
        if !queue.IsEmpty() {
            minutes++
        }
    }

    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if grid[r][c] == 1 {
                return -1
            }
        }
    }
    return minutes
}

func orangesRotting_slice(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    queue := [][2]int{} // queue of [row, col]
    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if grid[r][c] == 2 {
                queue = append(queue, [2]int{r, c})
            }
        }
    }

    directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

    minutes := 0
    for len(queue) > 0 {
        for _, orange := range queue {
            queue = queue[1:]
            for _, dir := range directions {
                r, c := orange[0]+dir[0], orange[1]+dir[1]
                if r < 0 || r >= m || c < 0 || c >= n {
                    continue
                }
                if grid[r][c] != 1 {
                    continue
                }
                grid[r][c] = 2
                queue = append(queue, [2]int{r, c})
            }
        }
        if len(queue) > 0 {
            minutes++
        }
    }

    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if grid[r][c] == 1 {
                return -1
            }
        }
    }
    return minutes
}

/**
 * 542. 01 Matrix
 * a multi-source BFS approach
 */
func updateMatrix(mat [][]int) [][]int {
    if mat == nil || len(mat) == 0 || len(mat[0]) == 0 {
        return [][]int{}
    }

    m, n := len(mat), len(mat[0])
    queue := make([][]int, 0)
    MAX_VALUE := m * n

    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            if mat[i][j] == 0 {
                queue = append(queue, []int{i, j})
            } else {
                mat[i][j] = MAX_VALUE
            }
        }
    }

    directions := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

    for len(queue) > 0 {
        cell := queue[0]
        queue = queue[1:]
        for _, dir := range directions {
            r, c := cell[0]+dir[0], cell[1]+dir[1] // check if there's a shorter path
            if r >= 0 && r < m && c >= 0 && c < n && mat[r][c] > mat[cell[0]][cell[1]]+1 {
                queue = append(queue, []int{r, c})
                mat[r][c] = mat[cell[0]][cell[1]] + 1
            }
        }
    }

    return mat
}

/**
 * 199. Binary Tree Right Side View
 */
func rightSideView(root *TreeNode) []int {
    res := []int{}
    if root == nil {
        return res
    }

    queue := []*TreeNode{root}
    for len(queue) > 0 {
        levelSize := len(queue)
        for i := 0; i < levelSize; i++ {
            cur := queue[0]
            queue = queue[1:]

            if cur.Left != nil {
                queue = append(queue, cur.Left)
            }
            if cur.Right != nil {
                queue = append(queue, cur.Right)
            }

            if i == levelSize-1 {
                res = append(res, cur.Val)
            }
        }
    }
    return res
}

/**
 * 102. Binary Tree Level Order Traversal
 */
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return [][]int{}
    }

    var res [][]int
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        qlen := len(queue)
        var level []int
        for i := 0; i < qlen; i++ {
            node := queue[0]
            queue = queue[1:]

            level = append(level, node.Val)
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        res = append(res, level)
    }
    return res
}
