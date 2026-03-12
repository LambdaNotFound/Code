package backtracking

import "strings"

/**
 * 37. Sudoku Solver
 */
func solveSudoku(board [][]byte) {
	var isValid func(int, int, byte) bool
	isValid = func(row, col int, num byte) bool {
		for i := 0; i < 9; i++ { // current (row, col)
			if board[row][i] == num || board[i][col] == num {
				return false
			}
		}

		startRow, startCol := (row/3)*3, (col/3)*3 // sub 3 x 3 board
		for i := startRow; i < startRow+3; i++ {
			for j := startCol; j < startCol+3; j++ {
				if board[i][j] == num {
					return false
				}
			}
		}
		return true
	}

	var dfs func() bool
	dfs = func() bool {
		for row := 0; row < 9; row++ {
			for col := 0; col < 9; col++ {
				if board[row][col] == '.' {
					for num := byte('1'); num <= '9'; num++ { // enumerate
						if isValid(row, col, num) {
							board[row][col] = num
							if dfs() {
								return true
							}
							board[row][col] = '.'
						}
					}
					return false
				}
			}
		}
		return true
	}
	dfs()
}

/**
 * 51. N-Queens
 * 52. N-Queens II
 */
func solveNQueens(n int) [][]string {
	var result [][]string

	board := make([]string, n)
	for i := range board {
		board[i] = strings.Repeat(".", n)
	}

	var isValid func(int, int) bool
	isValid = func(row, col int) bool {
		// check all rows above
		for i := 0; i < row; i++ {
			if board[i][col] == 'Q' {
				return false
			}
		}
		// upper-left diagonal
		for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
			if board[i][j] == 'Q' {
				return false
			}
		}
		// upper-right diagonal
		for i, j := row-1, col+1; i >= 0 && j < n; i, j = i-1, j+1 {
			if board[i][j] == 'Q' {
				return false
			}
		}
		return true
	}

	var dfs func(int)
	dfs = func(row int) {
		if row == n {
			temp := make([]string, n)
			copy(temp, board) // deep copy
			result = append(result, temp)
			return
		}

		for col := 0; col < n; col++ {
			if isValid(row, col) {
				newRow := []byte(board[row]) // []byte <> string
				newRow[col] = 'Q'
				board[row] = string(newRow)

				dfs(row + 1)

				newRow[col] = '.'
				board[row] = string(newRow)
			}
		}
	}

	dfs(0)
	return result
}
