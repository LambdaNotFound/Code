package solid_coding

/**
 * 54. Spiral Matrix
 */
func spiralOrder(matrix [][]int) []int {
	top, bottom := 0, len(matrix)-1
	left, right := 0, len(matrix[0])-1
	spiral := []int{}

	for top <= bottom && left <= right {
		for i := left; i <= right; i++ {
			spiral = append(spiral, matrix[top][i])
		}
		top++

		for j := top; j <= bottom; j++ {
			spiral = append(spiral, matrix[j][right])
		}
		right--

		if top <= bottom {
			for k := right; k >= left; k-- {
				spiral = append(spiral, matrix[bottom][k])
			}
			bottom--
		}

		if left <= right {
			for l := bottom; l >= top; l-- {
				spiral = append(spiral, matrix[l][left])
			}
			left++
		}
	}

	return spiral
}

func spiralOrderWithMap(matrix [][]int) []int {
	nextDirection := map[string]string{
		"top":    "right",
		"right":  "bottom",
		"bottom": "left",
		"left":   "top",
	}

	m, n := len(matrix), len(matrix[0])
	res := []int{}
	rowTop, rowBottom := 0, m-1
	colLeft, colRight := 0, n-1
	direction := "top"

	for len(res) != m*n {
		if direction == "top" {
			for j := colLeft; j <= colRight; j++ {
				res = append(res, matrix[rowTop][j])
			}
			rowTop++
		} else if direction == "right" {
			for i := rowTop; i <= rowBottom; i++ {
				res = append(res, matrix[i][colRight])
			}
			colRight--
		} else if direction == "bottom" {
			for j := colRight; j >= colLeft; j-- {
				res = append(res, matrix[rowBottom][j])
			}
			rowBottom--
		} else if direction == "left" {
			for i := rowBottom; i >= rowTop; i-- {
				res = append(res, matrix[i][colLeft])
			}
			colLeft++
		}
		direction = nextDirection[direction]
	}

	return res
}
