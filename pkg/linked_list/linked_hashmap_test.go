package linked_list

import (
	"sync"
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
		{
			name: "leetcode_example1",
			// Initial queue: [2, 3, 5] — all unique, first unique = 2.
			init: []int{2, 3, 5},
			ops: []op{
				show(2),
				add(5), show(2), // 5 → duplicate; uniques: [2, 3]
				add(2), show(3), // 2 → duplicate; uniques: [3]
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
				add(3), show(3), // 3 is new; first unique = 3
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
		{
			name: "insertion_order_preserved_across_adds",
			// Unique elements must surface in their original insertion order,
			// regardless of when duplicates are introduced.
			init: []int{1, 2, 3},
			ops: []op{
				show(1),
				add(1), show(2), // 1 → duplicate; front advances to 2
				add(2), show(3), // 2 → duplicate; front advances to 3
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
				add(10), show(5), // 10 is new and unique — but 5 is still first
				add(5), show(3), // 5 → duplicate; first unique = 3
				add(3), show(8), // 3 → duplicate; first unique = 8
				add(8), show(10), // 8 → duplicate; first unique = 10
			},
		},
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
		{
			name: "show_first_unique_is_idempotent",
			// Repeated calls must return the same value without mutating state.
			init: []int{3, 1, 2},
			ops:  repeat(10, show(3)),
		},
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
				add(4), show(3), // 4 is new but 3 is still first
				add(3), show(4), // 3 → duplicate; 4 becomes first
			},
		},
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

/**
 * Concurrency tests — always run with the race detector:
 *
 *   go test -race ./pkg/linked_list/... -run Test_FirstUnique_Concurrent -v
 *
 * Each sub-test launches multiple goroutines that call Add and/or
 * ShowFirstUnique concurrently. The race detector reports any unsynchronised
 * access even if the test's functional assertions happen to pass.
 *
 * Four scenarios are covered:
 *
 *   1. concurrent_adds_distinct_values
 *      Pure write contention: N goroutines each add a unique value.
 *      The RWMutex write lock must serialise all freq++ / list-append ops.
 *
 *   2. concurrent_reads_only
 *      Pure read contention: N goroutines call ShowFirstUnique in parallel.
 *      The RLock must allow all readers in simultaneously and return a
 *      consistent value every time.
 *
 *   3. concurrent_mixed_reads_and_writes
 *      Simultaneous readers and writers: half the goroutines add new values,
 *      the other half call ShowFirstUnique. Exercises the RLock/Lock interplay
 *      under real mixed load; no assertion on specific values, only no race.
 *
 *   4. concurrent_same_value_becomes_duplicate
 *      Multiple goroutines race to add the same value. Exactly one of them
 *      will hit freq==2 and remove the node from the list; the rest hit the
 *      freq>2 no-op branch. After the WaitGroup drains the result must be
 *      deterministic: that value is gone, the next unique leads.
 */
func Test_FirstUnique_Concurrent(t *testing.T) {
	const goroutines = 50 // enough to create scheduling pressure without being slow

	// ---------------------------------------------------------------------------
	// 1. Pure write contention — N goroutines each add a distinct value.
	//    After all complete, every added value has freq==1 and lives in the list.
	//    ShowFirstUnique must return a valid value (not -1) and must not race.
	// ---------------------------------------------------------------------------
	t.Run("concurrent_adds_distinct_values", func(t *testing.T) {
		fu := ConstructorLinkedHashmap([]int{})

		var wg sync.WaitGroup
		for i := 1; i <= goroutines; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				fu.Add(val)
			}(i)
		}
		wg.Wait()

		// All goroutines added a unique value; the list must be non-empty.
		got := fu.ShowFirstUnique()
		assert.NotEqual(t, -1, got, "at least one unique must remain after concurrent distinct adds")
	})

	// ---------------------------------------------------------------------------
	// 2. Pure read contention — N goroutines call ShowFirstUnique simultaneously.
	//    The RWMutex must admit all readers concurrently; every call must see
	//    the same value (no torn reads, no stale state).
	// ---------------------------------------------------------------------------
	t.Run("concurrent_reads_only", func(t *testing.T) {
		fu := ConstructorLinkedHashmap([]int{7, 3, 5}) // 7 is first unique

		results := make([]int, goroutines)
		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = fu.ShowFirstUnique()
			}(i)
		}
		wg.Wait()

		for i, v := range results {
			assert.Equal(t, 7, v, "goroutine %d: ShowFirstUnique() must return 7 (no writes occurred)", i)
		}
	})

	// ---------------------------------------------------------------------------
	// 3. Mixed reads and writes — writer goroutines add new unique values while
	//    reader goroutines call ShowFirstUnique. Validates that RLock/Lock do not
	//    interleave in a way that corrupts the list or exposes a partial state.
	//    Functional correctness is hard to pin down under arbitrary scheduling,
	//    so the assertion is minimal: ShowFirstUnique must never panic and must
	//    return either -1 or a value that was actually added (not garbage).
	// ---------------------------------------------------------------------------
	t.Run("concurrent_mixed_reads_and_writes", func(t *testing.T) {
		fu := ConstructorLinkedHashmap([]int{})

		// Writers add values 1..goroutines, each exactly once.
		addedVals := make(map[int]bool, goroutines)
		for i := 1; i <= goroutines; i++ {
			addedVals[i] = true
		}

		var wg sync.WaitGroup

		// Half the goroutines write.
		for i := 1; i <= goroutines; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				fu.Add(val)
			}(i)
		}

		// Half the goroutines read concurrently with those writes.
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				v := fu.ShowFirstUnique()
				// v is either -1 (nothing added yet) or a value we actually added.
				assert.True(t, v == -1 || addedVals[v],
					"ShowFirstUnique() returned %d which was never added", v)
			}()
		}

		wg.Wait()
	})

	// ---------------------------------------------------------------------------
	// 4. Same value added by N goroutines — deterministic final state.
	//    Initial state: {1, 2} — both unique; 1 is first.
	//    N goroutines all call Add(1). The first to execute hits freq==2 and
	//    removes 1 from the list; all subsequent hits are freq>2 no-ops.
	//    After the WaitGroup drains, 1 must be gone and 2 must lead.
	// ---------------------------------------------------------------------------
	t.Run("concurrent_same_value_becomes_duplicate", func(t *testing.T) {
		fu := ConstructorLinkedHashmap([]int{1, 2}) // freq[1]=1, freq[2]=1

		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fu.Add(1) // exactly one goroutine hits freq==2; rest are no-ops
			}()
		}
		wg.Wait()

		// 1 has been removed (freq >= 2); 2 is the sole unique.
		assert.Equal(t, 2, fu.ShowFirstUnique(),
			"after concurrent duplicate adds of 1, the first unique must be 2")
	})
}
