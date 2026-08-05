package utils

// DeepCopyMatrix returns a new 2D slice with the same values as m — each row
// gets its own backing array, so mutating the copy never affects m.
func DeepCopyMatrix[T any](m [][]T) [][]T {
	cp := make([][]T, len(m))
	for i, row := range m {
		cp[i] = make([]T, len(row))
		copy(cp[i], row)
	}
	return cp
}
