package solid_coding

import "sort"

/**
 * 54. Spiral Matrix
 */
func spiralOrder(matrix [][]int) []int {
    nextDirection := map[string]string{
        "top":    "right",
        "right":  "bottom",
        "bottom": "left",
        "left":   "top",
    }

    m, n := len(matrix), len(matrix[0])
    res := []int{}
    rtop, rbottom := 0, m-1
    cleft, cright := 0, n-1
    direction := "top"

    for len(res) != m*n {
        if direction == "top" {
            for j := cleft; j <= cright; j++ {
                res = append(res, matrix[rtop][j])
            }
            rtop++
        } else if direction == "right" {
            for i := rtop; i <= rbottom; i++ {
                res = append(res, matrix[i][cright])
            }
            cright--
        } else if direction == "bottom" {
            for j := cright; j >= cleft; j-- {
                res = append(res, matrix[rbottom][j])
            }
            rbottom--
        } else if direction == "left" {
            for i := rbottom; i >= rtop; i-- {
                res = append(res, matrix[i][cleft])
            }
            cleft++
        }
        direction = nextDirection[direction]
    }

    return res
}

/*
 * 56. Merge Intervals
 *
 * [[1,4],[2,3]] => [[1,4]]
 */
func merge(intervals [][]int) [][]int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })

    result := [][]int{}
    result = append(result, intervals[0])
    for _, interval := range intervals {
        current := result[len(result)-1]
        if current[1] < interval[0] {
            result = append(result, interval)
        } else {
            if current[1] < interval[1] {
                current[1] = interval[1]
            }
        }
    }

    return result
}
