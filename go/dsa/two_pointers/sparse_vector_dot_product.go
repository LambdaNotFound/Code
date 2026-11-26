package twopointers

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
