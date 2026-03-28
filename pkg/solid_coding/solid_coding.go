package solid_coding

/**
 * 54. Spiral Matrix
 */
func spiralOrder(matrix [][]int) []int {
	top, bottom := 0, len(matrix)-1
	left, right := 0, len(matrix[0])-1
	res := make([]int, 0, len(matrix)*len(matrix[0])) // pre-allocate exact capacity

	for left <= right && top <= bottom {
		// Traverse top row left → right
		for col := left; col <= right; col++ {
			res = append(res, matrix[top][col])
		}
		top++

		// Traverse right column top → bottom
		for row := top; row <= bottom; row++ {
			res = append(res, matrix[row][right])
		}
		right--

		// Traverse bottom row right → left (only if row still valid)
		if top <= bottom {
			for col := right; col >= left; col-- {
				res = append(res, matrix[bottom][col])
			}
			bottom--
		}

		// Traverse left column bottom → top (only if col still valid)
		if left <= right {
			for row := bottom; row >= top; row-- {
				res = append(res, matrix[row][left])
			}
			left++
		}
	}

	return res
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
	numRows, numCols := len(matrix), len(matrix[0])

	// Extra variable since matrix[0][0] can't represent both
	// row 0 and col 0 simultaneously
	firstColHasZero := false

	// Pass 1: use first row/col as markers
	for row := 0; row < numRows; row++ {
		// Check col 0 separately since it shares matrix[0][0]
		if matrix[row][0] == 0 {
			firstColHasZero = true
		}
		// Start col from 1 — col 0 is handled by firstColHasZero
		for col := 1; col < numCols; col++ {
			if matrix[row][col] == 0 {
				matrix[row][0] = 0 // mark row
				matrix[0][col] = 0 // mark col
			}
		}
	}

	// Pass 2: zero out cells based on markers in first row/col
	// Start from row=1, col=1 — don't touch the markers yet!
	for row := 1; row < numRows; row++ {
		for col := 1; col < numCols; col++ {
			if matrix[row][0] == 0 || matrix[0][col] == 0 {
				matrix[row][col] = 0
			}
		}
	}

	// Pass 3: handle first row using matrix[0][0] as marker
	if matrix[0][0] == 0 {
		for col := 0; col < numCols; col++ {
			matrix[0][col] = 0
		}
	}

	// Pass 4: handle first col using firstColHasZero
	if firstColHasZero {
		for row := 0; row < numRows; row++ {
			matrix[row][0] = 0
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
