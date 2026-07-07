package math

/**
 * 311. Sparse Matrix Multiplication
 *
 * Given two matrices mat1 of size m × k and mat2 of size k × n, return their product matrix result of size m × n.
 *
 * result[i][j] = sum of mat1[i][t] * mat2[t][j] for t in [0, k)
 *
 * Input:  mat1 (m×k), mat2 (k×n)
 * Output: result matrix (m×n)
 *
 * result[i][j] = Σ mat1[i][t] * mat2[t][j]
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

/**
 * 1570. Dot Product of Two Sparse Vectors
 *
 * Input:  vec1 (1×n), vec2 (1×n)
 * Output: single scalar
 *
 * result = Σ vec1[i] * vec2[i]
 */
type SparseVector struct {
	entries []sparseVectorEntry // only non-zero (idx, val) pairs
}

type sparseVectorEntry struct {
	idx, val int
}

func Constructor(nums []int) SparseVector {
	sv := SparseVector{}
	for i, val := range nums {
		if val != 0 {
			sv.entries = append(sv.entries, sparseVectorEntry{i, val})
		}
	}
	return sv
}

// dot product via two-pointer merge — O(L1 + L2)
func (v *SparseVector) dotProduct(vec SparseVector) int {
	result := 0
	i, j := 0, 0
	for i < len(v.entries) && j < len(vec.entries) {
		ei, ej := v.entries[i], vec.entries[j]
		if ei.idx == ej.idx {
			result += ei.val * ej.val
			i++
			j++
		} else if ei.idx < ej.idx {
			i++
		} else {
			j++
		}
	}
	return result
}

// hashmap based solution
type SparseVectorMap struct {
	nonZero map[int]int // idx → val
}

func ConstructorSparseVector(nums []int) SparseVectorMap {
	sv := SparseVectorMap{nonZero: map[int]int{}}
	for i, val := range nums {
		if val != 0 {
			sv.nonZero[i] = val
		}
	}
	return sv
}

// iterate shorter vector, look up in longer
func (v *SparseVectorMap) dotProductMap(vec SparseVectorMap) int {
	// always iterate the smaller map
	a, b := v.nonZero, vec.nonZero
	if len(a) > len(b) {
		a, b = b, a
	}
	result := 0
	for idx, val := range a {
		result += val * b[idx] // b[idx]=0 if missing
	}
	return result
}
