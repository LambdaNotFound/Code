package math

/**
 * 36. Valid Sudoku
 *
 * Each row must contain the digits 1-9 without repetition.
 * Each column must contain the digits 1-9 without repetition.
 * Each of the nine 3 x 3 sub-boxes of the grid must contain the digits 1-9 without repetition.
 */

/*
┌─────────┬─────────┬─────────┐
│ sq(0,0) │ sq(0,1) │ sq(0,2) │
│  [0]    │  [1]    │  [2]    │
├─────────┼─────────┼─────────┤
│ sq(1,0) │ sq(1,1) │ sq(1,2) │
│  [3]    │  [4]    │  [5]    │
├─────────┼─────────┼─────────┤
│ sq(2,0) │ sq(2,1) │ sq(2,2) │
│  [6]    │  [7]    │  [8]    │
└─────────┴─────────┴─────────┘
*/
func isValidSudoku(board [][]byte) bool {
	var rows, columns, squares [9][9]bool
	for i, row := range board {
		for j, v := range row {
			if v != '.' {
				k := int(v) - 49
				if rows[i][k] || columns[j][k] || squares[i/3*3+j/3][k] {
					return false
				}
				rows[i][k], columns[j][k], squares[i/3*3+j/3][k] = true, true, true
			}
		}
	}
	return true
}
