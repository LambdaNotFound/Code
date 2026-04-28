package heap

import (
	"container/heap"
	. "gocode/types"
	"gocode/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeKLists(t *testing.T) {
    testCases := []struct {
        name     string
        input    []*ListNode
        expected *ListNode
    }{
        {
            "case 1",
            []*ListNode{
                utils.CreateLinkedList([]int{1, 4, 5}),
                utils.CreateLinkedList([]int{1, 3, 4}),
                utils.CreateLinkedList([]int{2, 6}),
            },
            utils.CreateLinkedList([]int{1, 1, 2, 3, 4, 4, 5, 6}),
        },
        {
            "case 2",
            []*ListNode{},
            utils.CreateLinkedList([]int{}),
        },
        {
            "case 3",
            []*ListNode{
                utils.CreateLinkedList([]int{}),
            },
            utils.CreateLinkedList([]int{}),
        },
        {
            "single_list",
            []*ListNode{
                utils.CreateLinkedList([]int{1, 2, 3}),
            },
            utils.CreateLinkedList([]int{1, 2, 3}),
        },
        {
            "single_element_lists",
            []*ListNode{
                utils.CreateLinkedList([]int{3}),
                utils.CreateLinkedList([]int{1}),
                utils.CreateLinkedList([]int{2}),
            },
            utils.CreateLinkedList([]int{1, 2, 3}),
        },
        {
            "already_sorted",
            []*ListNode{
                utils.CreateLinkedList([]int{1, 2}),
                utils.CreateLinkedList([]int{3, 4}),
            },
            utils.CreateLinkedList([]int{1, 2, 3, 4}),
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := mergeKLists(tc.input)
            equal := utils.VerifyLinkedLists(tc.expected, result)
            assert.Equal(t, true, equal)
        })
    }
}

func Test_kClosest(t *testing.T) {
    tests := []struct {
        name     string
        points   [][]int
        k        int
        expected [][]int
    }{
        {
            name:     "basic_case",
            points:   [][]int{{1, 3}, {-2, 2}},
            k:        1,
            expected: [][]int{{-2, 2}},
        },
        {
            name:     "two_closest",
            points:   [][]int{{3, 3}, {5, -1}, {-2, 4}},
            k:        2,
            expected: [][]int{{3, 3}, {-2, 4}},
        },
        {
            name:     "all_points",
            points:   [][]int{{2, 2}, {1, 1}, {0, 0}},
            k:        3,
            expected: [][]int{{2, 2}, {1, 1}, {0, 0}}, // any order
        },
        {
            name:     "single_point",
            points:   [][]int{{10, 10}},
            k:        1,
            expected: [][]int{{10, 10}},
        },
        {
            name:     "origin_point",
            points:   [][]int{{0, 0}, {1, 1}},
            k:        1,
            expected: [][]int{{0, 0}},
        },
        {
            name:     "negatives_are_closer",
            points:   [][]int{{-1, -1}, {5, 5}, {2, 2}},
            k:        1,
            expected: [][]int{{-1, -1}},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := kClosest(tt.points, tt.k)

            // Since order of closest points can vary, we check via sets
            assert.Equal(t, len(tt.expected), len(result))

            // convert slices to maps for comparison
            toMap := func(points [][]int) map[[2]int]bool {
                m := make(map[[2]int]bool)
                for _, p := range points {
                    m[[2]int{p[0], p[1]}] = true
                }
                return m
            }

            gotMap := toMap(result)
            wantMap := toMap(tt.expected)
            assert.Equal(t, wantMap, gotMap)
        })
    }
}

func Test_MedianFinder(t *testing.T) {
    tests := []struct {
        name     string
        nums     []int
        expected []float64 // expected median after each AddNum
    }{
        {
            name:     "single",
            nums:     []int{1},
            expected: []float64{1},
        },
        {
            name:     "two",
            nums:     []int{1, 2},
            expected: []float64{1, 1.5},
        },
        {
            name:     "three",
            nums:     []int{1, 2, 3},
            expected: []float64{1, 1.5, 2},
        },
        {
            name:     "four",
            nums:     []int{1, 2, 3, 4},
            expected: []float64{1, 1.5, 2, 2.5},
        },
        {
            name:     "reverse_order",
            nums:     []int{4, 3, 2, 1},
            expected: []float64{4, 3.5, 3, 2.5},
        },
        {
            name:     "duplicates",
            nums:     []int{2, 2, 2, 2},
            expected: []float64{2, 2, 2, 2},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mf := Constructor()
            for i, num := range tt.nums {
                mf.AddNum(num)
                assert.InDelta(t, tt.expected[i], mf.FindMedian(), 1e-9,
                    "median after AddNum(%d) at index %d", num, i)
            }
        })
    }
}

// ---------------------------------------------------------------------------
// 480. Sliding Window Median
//
// Both implementations (lazy-deletion two-heap and insertion-sort naive) are
// driven with the same table so they stay in sync.
//
// Test strategy:
//   - LeetCode canonical example (odd k=3, mix of positives and negatives).
//   - k=1: every element is trivially its own median.
//   - k=2 (even): result is always the average of two adjacent elements.
//   - k=len(nums): single window — one median for the entire array.
//   - All-same values: median equals that value regardless of k.
//   - Two elements, even k: straightforward average.
//   - Negatives only: ensures heap ordering works below zero.
//   - All-duplicate odd k: median stays constant throughout.
//   - Large identical values (even k): validates float arithmetic on big ints.
// ---------------------------------------------------------------------------
func Test_medianSlidingWindow(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []float64
	}{
		{
			// [1,3,-1]→1, [3,-1,-3]→-1, [-1,-3,5]→-1, [-3,5,3]→3, [5,3,6]→5, [3,6,7]→6
			name:     "leetcode_example1_odd_k",
			nums:     []int{1, 3, -1, -3, 5, 3, 6, 7},
			k:        3,
			expected: []float64{1, -1, -1, 3, 5, 6},
		},
		{
			// k=1: each element is its own single-element window.
			name:     "k1_trivial_each_element_is_median",
			nums:     []int{5, 3, 1, 2, 4},
			k:        1,
			expected: []float64{5, 3, 1, 2, 4},
		},
		{
			// k=2 (even): median = (a+b)/2.0 for each adjacent pair.
			name:     "k2_even_window_averages",
			nums:     []int{1, 2, 3, 4},
			k:        2,
			expected: []float64{1.5, 2.5, 3.5},
		},
		{
			// k == len(nums): exactly one window covering the whole array.
			// sorted [1,2,3,4] → median = (2+3)/2 = 2.5
			name:     "k_equals_array_length",
			nums:     []int{3, 1, 2, 4},
			k:        4,
			expected: []float64{2.5},
		},
		{
			// All values identical — median is always that value.
			name:     "all_same_values_even_k",
			nums:     []int{5, 5, 5, 5},
			k:        2,
			expected: []float64{5, 5, 5},
		},
		{
			// Two elements, even k=2: one window, average of the two.
			name:     "two_elements_k2",
			nums:     []int{2, 1},
			k:        2,
			expected: []float64{1.5},
		},
		{
			// Negative numbers only, odd k=3.
			// [-7,-1,-5]→-5, [-1,-5,-2]→-2
			name:     "negatives_only_odd_k",
			nums:     []int{-7, -1, -5, -2},
			k:        3,
			expected: []float64{-5, -2},
		},
		{
			// All duplicates, odd k=3: median stays at 1 for every window.
			name:     "all_duplicates_odd_k",
			nums:     []int{1, 1, 1, 1, 1},
			k:        3,
			expected: []float64{1, 1, 1},
		},
		{
			// Two identical large integers, even k=2.
			// (2147483647 + 2147483647) / 2 = 2147483647 — no integer overflow in float.
			name:     "large_identical_values_even_k",
			nums:     []int{2147483647, 2147483647},
			k:        2,
			expected: []float64{2147483647},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := medianSlidingWindow(tt.nums, tt.k)
			assert.Equal(t, len(tt.expected), len(got), "result length mismatch")
			for i := range tt.expected {
				assert.InDelta(t, tt.expected[i], got[i], 1e-9,
					"window %d: medianSlidingWindow", i)
			}

			gotNaive := medianSlidingWindowInsertion(tt.nums, tt.k)
			assert.Equal(t, len(tt.expected), len(gotNaive), "naive result length mismatch")
			for i := range tt.expected {
				assert.InDelta(t, tt.expected[i], gotNaive[i], 1e-9,
					"window %d: medianSlidingWindowInsertion", i)
			}
		})
	}
}

func Test_ListNodeMinHeap(t *testing.T) {
    t.Run("push_pop_ordered", func(t *testing.T) {
        h := &ListNodeMinHeap{}
        heap.Init(h)

        heap.Push(h, &ListNode{Val: 5})
        heap.Push(h, &ListNode{Val: 1})
        heap.Push(h, &ListNode{Val: 3})

        assert.Equal(t, 3, h.Len())

        v1 := heap.Pop(h).(*ListNode)
        assert.Equal(t, 1, v1.Val)
        v2 := heap.Pop(h).(*ListNode)
        assert.Equal(t, 3, v2.Val)
        v3 := heap.Pop(h).(*ListNode)
        assert.Equal(t, 5, v3.Val)
        assert.Equal(t, 0, h.Len())
    })

    t.Run("swap_and_less_via_heap_ordering", func(t *testing.T) {
        h := &ListNodeMinHeap{
            {Val: 10},
            {Val: 2},
            {Val: 7},
        }
        heap.Init(h) // establishes heap property via Less/Swap
        assert.Equal(t, 2, (*h)[0].Val)
    })
}

func Test_Heap(t *testing.T) {
    t.Run("min_heap_push_peek_len_pop", func(t *testing.T) {
        h := NewHeap(func(a, b int) bool { return a < b })
        assert.Equal(t, 0, h.Len())

        heap.Push(h, 5)
        heap.Push(h, 2)
        heap.Push(h, 8)
        heap.Push(h, 1)

        assert.Equal(t, 4, h.Len())
        assert.Equal(t, 1, h.Peek()) // min at top

        v := heap.Pop(h).(int)
        assert.Equal(t, 1, v)
        assert.Equal(t, 3, h.Len())
        assert.Equal(t, 2, h.Peek())

        v = heap.Pop(h).(int)
        assert.Equal(t, 2, v)
        v = heap.Pop(h).(int)
        assert.Equal(t, 5, v)
        v = heap.Pop(h).(int)
        assert.Equal(t, 8, v)
        assert.Equal(t, 0, h.Len())
    })

    t.Run("max_heap_push_peek_pop", func(t *testing.T) {
        h := NewHeap(func(a, b int) bool { return a > b })
        heap.Push(h, 3)
        heap.Push(h, 10)
        heap.Push(h, 7)

        assert.Equal(t, 10, h.Peek()) // max at top

        v := heap.Pop(h).(int)
        assert.Equal(t, 10, v)
        assert.Equal(t, 7, h.Peek())
    })

    t.Run("single_element_pop", func(t *testing.T) {
        h := NewHeap(func(a, b int) bool { return a < b })
        heap.Push(h, 42)
        assert.Equal(t, 1, h.Len())
        v := heap.Pop(h).(int)
        assert.Equal(t, 42, v)
        assert.Equal(t, 0, h.Len())
    })

    t.Run("duplicates", func(t *testing.T) {
        h := NewHeap(func(a, b int) bool { return a < b })
        heap.Push(h, 3)
        heap.Push(h, 3)
        heap.Push(h, 3)
        assert.Equal(t, 3, h.Len())
        assert.Equal(t, 3, h.Peek())
        heap.Pop(h)
        assert.Equal(t, 2, h.Len())
    })
}