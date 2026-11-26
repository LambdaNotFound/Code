package solid_coding

import "slices"

/**
 * 54. Spiral Matrix
 */
func spiralOrder(matrix [][]int) []int {
	m, n := len(matrix), len(matrix[0])
	next := map[string]string{"top": "right", "right": "bottom", "bottom": "left", "left": "top"}
	res := []int{}

	dir := "top"
	topRow, rightCol, bottomRow, leftCol := 0, n-1, m-1, 0
	for len(res) < m*n {
		switch {
		case dir == "top":
			for c := leftCol; c <= rightCol; c++ {
				res = append(res, matrix[topRow][c])
			}
			topRow++
		case dir == "right":
			for r := topRow; r <= bottomRow; r++ {
				res = append(res, matrix[r][rightCol])
			}
			rightCol--
		case dir == "bottom":
			for c := rightCol; c >= leftCol; c-- {
				res = append(res, matrix[bottomRow][c])
			}
			bottomRow--
		case dir == "left":
			for r := bottomRow; r >= topRow; r-- {
				res = append(res, matrix[r][leftCol])
			}
			leftCol++
		}
		dir = next[dir]
	}

	return res
}

/**
 * 73. Set Matrix Zeroes
 */
func setZeroes(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])
	rows, cols := make([]bool, m), make([]bool, n)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] == 0 {
				rows[i] = true
				cols[j] = true
			}
		}
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if rows[i] || cols[j] {
				matrix[i][j] = 0
			}
		}
	}
}

func setZeroesOptimal(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])
	firstRowHasZero, firstColHasZero := false, false
	for i := 0; i < m; i++ {
		if matrix[i][0] == 0 {
			firstColHasZero = true
		}
	}
	for j := 0; j < n; j++ {
		if matrix[0][j] == 0 {
			firstRowHasZero = true
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if matrix[i][j] == 0 {
				matrix[i][0] = 0
				matrix[0][j] = 0
			}
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if matrix[i][0] == 0 || matrix[0][j] == 0 {
				matrix[i][j] = 0
			}
		}
	}

	if firstRowHasZero {
		for j := 0; j < n; j++ {
			matrix[0][j] = 0
		}
	}
	if firstColHasZero {
		for i := 0; i < m; i++ {
			matrix[i][0] = 0
		}
	}
}

/**
 * 48. Rotate Image
 *
 *      reverse + diagonal swap
 *  1, 2, 3    7, 8, 9    7, 4, 1
 *  4, 5, 6    4, 5, 6    8, 5, 2
 *  7, 8, 9    1, 2, 3    9, 6, 3
 *
 */
func rotateImage(matrix [][]int) {
	n := len(matrix)

	// slices.Reverse(matrix)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}
}

/**
 * 189. Rotate Array
 *
 * [1,2,3,4,5,6,7], k = 3 => [5,6,7,1,2,3,4]
 *
 * reverse subarray + full reverse
 * [0:n-k] [n-k,n]
 * 4,3,2,1, 7,6,5
 */
func rotateArray(nums []int, k int) {
	n := len(nums)
	k = k % n
	slices.Reverse(nums[:n-k])
	slices.Reverse(nums[n-k : n])
	slices.Reverse(nums)
}
