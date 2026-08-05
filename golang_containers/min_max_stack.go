package containers

import "cmp"

/**
 * MinMaxStack[T] — Stack with O(1) Minimum and Maximum (LeetCode 155 extended)
 *
 * A LIFO stack that answers GetMin() and GetMax() in O(1) at all times,
 * alongside the standard Push, Pop, and Peek operations.
 *
 * Algorithm — three parallel slices:
 *   items  : the main stack (LIFO order)
 *   minAux : minAux[i] = min(items[0], items[1], ..., items[i])
 *   maxAux : maxAux[i] = max(items[0], items[1], ..., items[i])
 *
 * Structural invariant:
 *   len(items) == len(minAux) == len(maxAux) at all times.
 *   Maintained by appending to all three on Push and truncating all three
 *   by one index on Pop.
 *
 * Why push to aux on every Push (not only on a new min/max):
 *   Consider pushing [3, 1, 1]. If minAux only records 1 once, popping
 *   the first 1 would incorrectly restore the min to 3 while the second 1
 *   is still on the stack. Pushing every time costs O(n) total aux space
 *   (the same order as the main stack) but eliminates this edge case entirely.
 *
 * Why cmp.Ordered (not a custom comparator):
 *   MinMaxStack tracks the mathematical min and max of the values themselves.
 *   There is only one sensible definition of min/max for ordered types, unlike
 *   Heap where "priority" is caller-defined (e.g. sort Tasks by deadline).
 *   If ordering by a field is needed, extract that field into a MinMaxStack
 *   of the field's type (e.g. MinMaxStack[time.Time]).
 *
 * Complexity:
 *   Push, Pop, Peek, GetMin, GetMax, Size, IsEmpty   O(1)
 *   Space                                             O(3n)
 */

// MinMaxStack is a generic LIFO stack with O(1) min and max queries.
// T must be an ordered type (any numeric type or string).
// The zero value is ready to use; NewMinMaxStack is provided for consistency.
type MinMaxStack[T cmp.Ordered] struct {
	items  []T // main stack
	minAux []T // minAux[i] = running minimum from bottom up to position i
	maxAux []T // maxAux[i] = running maximum from bottom up to position i
}

// NewMinMaxStack returns an empty, ready-to-use MinMaxStack.
func NewMinMaxStack[T cmp.Ordered]() *MinMaxStack[T] {
	return &MinMaxStack[T]{}
}

// Push adds x to the top of the stack and updates the min/max records.
// Time: O(1).
func (s *MinMaxStack[T]) Push(x T) {
	s.items = append(s.items, x)
	if len(s.minAux) == 0 {
		s.minAux = append(s.minAux, x)
		s.maxAux = append(s.maxAux, x)
	} else {
		s.minAux = append(s.minAux, min(x, s.minAux[len(s.minAux)-1]))
		s.maxAux = append(s.maxAux, max(x, s.maxAux[len(s.maxAux)-1]))
	}
}

// Pop removes and returns the top element.
// Returns (zero, false) if the stack is empty — never panics.
// Time: O(1).
func (s *MinMaxStack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	n := len(s.items) - 1
	top := s.items[n]
	s.items = s.items[:n]
	s.minAux = s.minAux[:n]
	s.maxAux = s.maxAux[:n]
	return top, true
}

// Peek returns the top element without removing it.
// Returns (zero, false) if the stack is empty — never panics.
// Time: O(1).
func (s *MinMaxStack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// GetMin returns the minimum value currently in the stack.
// Returns (zero, false) if the stack is empty.
// Time: O(1).
func (s *MinMaxStack[T]) GetMin() (T, bool) {
	if len(s.minAux) == 0 {
		var zero T
		return zero, false
	}
	return s.minAux[len(s.minAux)-1], true
}

// GetMax returns the maximum value currently in the stack.
// Returns (zero, false) if the stack is empty.
// Time: O(1).
func (s *MinMaxStack[T]) GetMax() (T, bool) {
	if len(s.maxAux) == 0 {
		var zero T
		return zero, false
	}
	return s.maxAux[len(s.maxAux)-1], true
}

// Size returns the number of elements currently in the stack.
// Time: O(1).
func (s *MinMaxStack[T]) Size() int {
	return len(s.items)
}

// IsEmpty reports whether the stack has no elements.
// Time: O(1).
func (s *MinMaxStack[T]) IsEmpty() bool {
	return len(s.items) == 0
}
