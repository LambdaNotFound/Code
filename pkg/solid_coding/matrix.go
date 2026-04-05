package solid_coding

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
 * 13. Roman to Integer
 */
func romanToInt(s string) int {
	toInt := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	res := 0
	for i := range s {
		if i+1 < len(s) && toInt[s[i]] < toInt[s[i+1]] {
			res -= toInt[s[i]]
		} else {
			res += toInt[s[i]]
		}
	}

	return res
}
