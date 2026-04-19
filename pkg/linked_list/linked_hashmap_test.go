package linked_list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * 1429. First Unique Number
 *
 * Test strategy:
 *   - Verify the three LeetCode canonical examples.
 *   - Cover every code branch in Add (freq==1, freq==2, freq>2).
 *   - Verify insertion-order semantics: the first unique is always the one
 *     that was inserted earliest, even if later values were also added first.
 *   - Verify ShowFirstUnique is read-only (idempotent).
 *   - Verify correct behaviour on edge inputs: empty init, all-duplicate init,
 *     single element, high frequency, new unique after list empties.
 *
 * Table-driven pattern:
 *   Each test case is a named sequence of typed operations applied to a single
 *   FirstUnique instance. Two operation kinds are supported:
 *
 *     add(v)   — calls fu.Add(v)
 *     show(w)  — calls fu.ShowFirstUnique() and asserts the result equals w
 *
 *   repeat(n, op) expands one operation into n consecutive copies; used to
 *   exercise the freq>2 no-op branch without cluttering the table.
 */

// opKind distinguishes the two operations a test step can perform.
type opKind int

const (
	opAdd  opKind = iota // call Add(val)
	opShow               // call ShowFirstUnique(), expect want
)

// op is a single test step.
type op struct {
	kind opKind
	val  int // opAdd: value to add | opShow: expected return value
}

// add returns an Add step.
func add(v int) op { return op{kind: opAdd, val: v} }

// show returns a ShowFirstUnique step that expects want.
func show(want int) op { return op{kind: opShow, val: want} }

// repeat returns n copies of o, used for high-frequency stress steps.
func repeat(n int, o op) []op {
	ops := make([]op, n)
	for i := range ops {
		ops[i] = o
	}
	return ops
}

func Test_FirstUnique(t *testing.T) {
	tests := []struct {
		name string
		init []int
		ops  []op
	}{
		// -----------------------------------------------------------------------
		// LeetCode canonical examples
		// -----------------------------------------------------------------------
		{
			name: "leetcode_example1",
			// Initial queue: [2, 3, 5] — all unique, first unique = 2.
			init: []int{2, 3, 5},
			ops: []op{
				show(2),
				add(5), show(2),  // 5 → duplicate; uniques: [2, 3]
				add(2), show(3),  // 2 → duplicate; uniques: [3]
				add(3), show(-1), // 3 → duplicate; uniques: []
			},
		},
		{
			name: "leetcode_example2",
			// Initial queue: all duplicates — no unique from the start.
			init: []int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
			ops: []op{
				show(-1),
				add(7), show(-1), // freq > 2: no-op, still no unique
				add(3), show(3),  // 3 is new; first unique = 3
			},
		},
		{
			name: "leetcode_example3",
			// New unique added after all initials become duplicates.
			init: []int{2, 3, 5},
			ops: []op{
				show(2),
				add(5), add(2), add(3),
				show(-1),
				add(4), show(4), // entirely new value → only unique
			},
		},

		// -----------------------------------------------------------------------
		// Edge inputs
		// -----------------------------------------------------------------------
		{
			name: "empty_init_returns_minus_one",
			init: []int{},
			ops:  []op{show(-1)},
		},
		{
			name: "single_unique_element",
			init: []int{42},
			ops:  []op{show(42)},
		},
		{
			name: "single_element_becomes_duplicate_on_add",
			init: []int{42},
			ops: []op{
				show(42),
				add(42), show(-1),
			},
		},
		{
			name: "all_duplicates_in_init_returns_minus_one",
			init: []int{1, 1, 2, 2, 3, 3},
			ops:  []op{show(-1)},
		},

		// -----------------------------------------------------------------------
		// Insertion-order semantics
		// -----------------------------------------------------------------------
		{
			name: "insertion_order_preserved_across_adds",
			// Unique elements must surface in their original insertion order,
			// regardless of when duplicates are introduced.
			init: []int{1, 2, 3},
			ops: []op{
				show(1),
				add(1), show(2),  // 1 → duplicate; front advances to 2
				add(2), show(3),  // 2 → duplicate; front advances to 3
				add(3), show(-1), // 3 → duplicate; list empty
			},
		},
		{
			name: "later_unique_does_not_jump_queue",
			// Adding a new unique value must append it after existing uniques,
			// never promote it to the front.
			init: []int{5, 3, 8},
			ops: []op{
				show(5),
				add(10), show(5),  // 10 is new and unique — but 5 is still first
				add(5), show(3),   // 5 → duplicate; first unique = 3
				add(3), show(8),   // 3 → duplicate; first unique = 8
				add(8), show(10),  // 8 → duplicate; first unique = 10
			},
		},

		// -----------------------------------------------------------------------
		// Add branch coverage: freq > 2 (no-op path)
		// -----------------------------------------------------------------------
		{
			name: "high_frequency_add_is_noop_no_crash",
			// Adding the same value many times must not corrupt the list or panic.
			// The freq>2 branch in Add must be a true no-op.
			init: []int{1},
			ops:  append(repeat(100, add(1)), show(-1)),
		},
		{
			name: "high_frequency_does_not_affect_other_uniques",
			// Spamming add(1) 50 times: only the freq==2 hit matters; the rest
			// are no-ops. The other uniques (2, 3) must be untouched.
			init: []int{1, 2, 3},
			ops:  append(repeat(50, add(1)), show(2)),
		},

		// -----------------------------------------------------------------------
		// ShowFirstUnique is read-only
		// -----------------------------------------------------------------------
		{
			name: "show_first_unique_is_idempotent",
			// Repeated calls must return the same value without mutating state.
			init: []int{3, 1, 2},
			ops:  repeat(10, show(3)),
		},

		// -----------------------------------------------------------------------
		// New unique added after list is fully emptied
		// -----------------------------------------------------------------------
		{
			name: "new_unique_recovers_after_empty_list",
			init: []int{1, 1},
			ops: []op{
				show(-1),
				add(5), show(5), // first entirely-new unique
			},
		},
		{
			name: "multiple_new_uniques_after_empty_list",
			init: []int{1, 1, 2, 2},
			ops: []op{
				show(-1),
				add(3), show(3),
				add(4), show(3),  // 4 is new but 3 is still first
				add(3), show(4),  // 3 → duplicate; 4 becomes first
			},
		},

		// -----------------------------------------------------------------------
		// Interleaved Add and ShowFirstUnique
		// -----------------------------------------------------------------------
		{
			name: "interleaved_add_and_show",
			init: []int{10, 20, 30},
			ops: []op{
				show(10),
				add(30), show(10), // 30 → dup
				add(10), show(20), // 10 → dup
				add(40), show(20), // 40 is new; 20 still leads
				add(20), show(40), // 20 → dup; 40 now leads
				add(40), show(-1), // 40 → dup; list empty
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fu := ConstructorLinkedHashmap(tc.init)
			for stepIdx, o := range tc.ops {
				switch o.kind {
				case opAdd:
					fu.Add(o.val)
				case opShow:
					assert.Equal(t, o.val, fu.ShowFirstUnique(),
						"step %d: ShowFirstUnique()", stepIdx)
				}
			}
		})
	}
}
