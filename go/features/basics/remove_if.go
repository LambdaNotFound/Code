package basics

// RemoveIf filters s in place, keeping the elements pred rejects. It reuses
// s's backing array, so the caller must not keep using s afterwards.
func RemoveIf[T any](s []T, pred func(T) bool) []T {
	result := s[:0] // reuse underlying array
	for _, v := range s {
		if !pred(v) {
			result = append(result, v)
		}
	}
	return result
}
