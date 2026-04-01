package math

/**
 * 311. Sparse Matrix Multiplication
 *
 * Given two matrices mat1 of size m × k and mat2 of size k × n, return their product matrix result of size m × n.
 *
 * result[i][j] = sum of mat1[i][t] * mat2[t][j] for t in [0, k)
 */
type sparseEntry struct {
	col, val int
}

func compress(mat [][]int) [][]sparseEntry {
	sparse := make([][]sparseEntry, len(mat))
	for i, row := range mat {
		for col, val := range row {
			if val != 0 {
				sparse[i] = append(sparse[i], sparseEntry{col, val})
			}
		}
	}
	return sparse
}

func multiplyCompressed(mat1 [][]int, mat2 [][]int) [][]int {
	m, n := len(mat1), len(mat2[0])
	result := make([][]int, m)
	for i := range result {
		result[i] = make([]int, n)
	}

	// compress both matrices upfront
	sparse1 := compress(mat1)
	sparse2 := compress(mat2) // compress mat2 by rows too

	for i := 0; i < m; i++ {
		for _, e1 := range sparse1[i] {
			// e1.col = t (shared dimension index)
			// multiply with row t of mat2
			for _, e2 := range sparse2[e1.col] {
				result[i][e2.col] += e1.val * e2.val
			}
		}
	}
	return result
}

func multiply(mat1 [][]int, mat2 [][]int) [][]int {
	m, k, n := len(mat1), len(mat1[0]), len(mat2[0])
	result := make([][]int, m)
	for i := range result {
		result[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for t := 0; t < k; t++ {
			if mat1[i][t] == 0 {
				continue // skip entire inner loop
			}
			for j := 0; j < n; j++ {
				result[i][j] += mat1[i][t] * mat2[t][j]
			}
		}
	}
	return result
}
