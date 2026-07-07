package linked_list

import (
	. "gocode/types"
	"gocode/utils"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_reverseList(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"multiple_nodes", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"single_node", []int{1}, []int{1}},
		{"two_nodes", []int{1, 2}, []int{2, 1}},
		{"nil", []int{}, []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			want := utils.CreateLinkedList(tc.expected)
			impls := []struct {
				name string
				fn   func(*ListNode) *ListNode
			}{
				{"iterative", reverseList_iterative},
				{"recursive", reverseList_recursive},
				{"recursive_alt", reverseList},
			}
			for _, impl := range impls {
				head := utils.CreateLinkedList(tc.input)
				result := impl.fn(head)
				assert.True(t, utils.VerifyLinkedLists(want, result), "%s failed", impl.name)
			}
		})
	}
}

func Test_mergeTwoLists(t *testing.T) {
	testCases := []struct {
		name     string
		l1, l2   []int
		expected []int
	}{
		{"interleaved", []int{1, 3, 5, 7}, []int{2, 4, 6, 8}, []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"l1_nil", []int{}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"l2_nil", []int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{"both_nil", []int{}, []int{}, []int{}},
		{"l1_shorter", []int{1, 2}, []int{3, 4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{"duplicates", []int{1, 1, 2}, []int{1, 2, 3}, []int{1, 1, 1, 2, 2, 3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeTwoLists(utils.CreateLinkedList(tc.l1), utils.CreateLinkedList(tc.l2))
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))
		})
	}
}

func Test_swapPairs(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"even_length", []int{1, 2, 3, 4}, []int{2, 1, 4, 3}},
		{"odd_length", []int{1, 2, 3}, []int{2, 1, 3}},
		{"single_node", []int{1}, []int{1}},
		{"nil", []int{}, []int{}},
		{"two_nodes", []int{1, 2}, []int{2, 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := swapPairs(utils.CreateLinkedList(tc.input))
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))
		})
	}
}

func Test_hasCycle(t *testing.T) {
	tests := []struct {
		name     string
		build    func() *ListNode
		expected bool
	}{
		{
			name: "no_cycle",
			build: func() *ListNode {
				n1 := &ListNode{Val: 1}
				n2 := &ListNode{Val: 2}
				n3 := &ListNode{Val: 3}
				n1.Next = n2
				n2.Next = n3
				return n1
			},
			expected: false,
		},
		{
			name: "simple_cycle",
			build: func() *ListNode {
				n1 := &ListNode{Val: 1}
				n2 := &ListNode{Val: 2}
				n1.Next = n2
				n2.Next = n1 // cycle back to n1
				return n1
			},
			expected: true,
		},
		{
			name: "longer_cycle",
			build: func() *ListNode {
				n1 := &ListNode{Val: 1}
				n2 := &ListNode{Val: 2}
				n3 := &ListNode{Val: 3}
				n4 := &ListNode{Val: 4}
				n1.Next = n2
				n2.Next = n3
				n3.Next = n4
				n4.Next = n2 // cycle back to n2
				return n1
			},
			expected: true,
		},
		{
			name: "single_node_no_cycle",
			build: func() *ListNode {
				return &ListNode{Val: 1}
			},
			expected: false,
		},
		{
			name: "single_node_cycle",
			build: func() *ListNode {
				n1 := &ListNode{Val: 1}
				n1.Next = n1 // self-cycle
				return n1
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := tt.build()
			result := hasCycle(head)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_middleNode(t *testing.T) {
	tests := []struct {
		name     string
		build    func() *ListNode
		expected []int
	}{
		{
			name: "odd_length",
			build: func() *ListNode {
				n1 := &ListNode{Val: 1}
				n2 := &ListNode{Val: 2}
				n3 := &ListNode{Val: 3}
				n4 := &ListNode{Val: 4}
				n5 := &ListNode{Val: 5}
				n1.Next, n2.Next, n3.Next, n4.Next = n2, n3, n4, n5
				return n1
			},
			expected: []int{3, 4, 5}, // middle is 3
		},
		{
			name: "even_length",
			build: func() *ListNode {
				n1 := &ListNode{Val: 1}
				n2 := &ListNode{Val: 2}
				n3 := &ListNode{Val: 3}
				n4 := &ListNode{Val: 4}
				n5 := &ListNode{Val: 5}
				n6 := &ListNode{Val: 6}
				n1.Next, n2.Next, n3.Next, n4.Next, n5.Next = n2, n3, n4, n5, n6
				return n1
			},
			expected: []int{4, 5, 6}, // middle is 4 (second middle)
		},
		{
			name: "single_node",
			build: func() *ListNode {
				return &ListNode{Val: 1}
			},
			expected: []int{1},
		},
		{
			name: "two_nodes",
			build: func() *ListNode {
				n1 := &ListNode{Val: 1}
				n2 := &ListNode{Val: 2}
				n1.Next = n2
				return n1
			},
			expected: []int{2}, // second node is the middle
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := tt.build()
			mid := middleNode(head)

			// Collect values from mid to end
			vals := []int{}
			for cur := mid; cur != nil; cur = cur.Next {
				vals = append(vals, cur.Val)
			}

			assert.Equal(t, tt.expected, vals)
		})
	}
}

func Test_copyRandomList(t *testing.T) {
	tests := []struct {
		name  string
		build func() *Node
	}{
		{
			name:  "nil_input",
			build: func() *Node { return nil },
		},
		{
			name: "single_node",
			build: func() *Node {
				return &Node{Val: 1}
			},
		},
		{
			name: "two_nodes_no_random",
			build: func() *Node {
				n1 := &Node{Val: 1}
				n2 := &Node{Val: 2}
				n1.Next = n2
				return n1
			},
		},
		{
			name: "random_points_forward",
			build: func() *Node {
				n1 := &Node{Val: 1}
				n2 := &Node{Val: 2}
				n1.Next = n2
				n1.Random = n2
				return n1
			},
		},
		{
			name: "random_points_backward",
			build: func() *Node {
				n1 := &Node{Val: 1}
				n2 := &Node{Val: 2}
				n3 := &Node{Val: 3}
				n1.Next, n2.Next = n2, n3
				n3.Random = n1
				return n1
			},
		},
		{
			name: "random_self_loop",
			build: func() *Node {
				n1 := &Node{Val: 1}
				n1.Random = n1
				return n1
			},
		},
		{
			name: "full_random_structure",
			build: func() *Node {
				n1 := &Node{Val: 7}
				n2 := &Node{Val: 13}
				n3 := &Node{Val: 11}
				n4 := &Node{Val: 10}
				n5 := &Node{Val: 1}
				n1.Next, n2.Next, n3.Next, n4.Next = n2, n3, n4, n5
				n2.Random = n1
				n3.Random = n5
				n4.Random = n3
				n5.Random = n1
				return n1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := tt.build()
			got := copyRandomList(head)

			assert.True(t, verifyCopyRandomList(head, got))
		})
	}
}

// verifyCopyRandomList checks that got is a valid deep copy of head:
// same structure (Next chain), same values, same Random relationships,
// and no shared node pointers between original and copy.
//
// TODO: simplify this method.
func verifyCopyRandomList(head, got *Node) bool {
	if head == nil && got == nil {
		return true
	}
	if head == nil || got == nil {
		return false
	}

	// Build index map for original list: node -> index
	origIdx := make(map[*Node]int)
	for i, cur := 0, head; cur != nil; cur, i = cur.Next, i+1 {
		origIdx[cur] = i
	}

	// Build ordered slice of original nodes by index
	n := len(origIdx)
	origNodes := make([]*Node, n)
	for cur := head; cur != nil; cur = cur.Next {
		origNodes[origIdx[cur]] = cur
	}

	// Traverse copy and verify structure + deep-copy property
	copyCur := got
	for i := 0; copyCur != nil; copyCur, i = copyCur.Next, i+1 {
		if i >= n {
			return false
		}
		orig := origNodes[i]

		// Deep copy: copy node must not be in original set
		if _, inOrig := origIdx[copyCur]; inOrig {
			return false
		}
		// Same value
		if copyCur.Val != orig.Val {
			return false
		}
		// Random: if orig has random at index j, copy must have random pointing to copy node at j
		if orig.Random != nil {
			if copyCur.Random == nil {
				return false
			}
			j := origIdx[orig.Random]
			// copyCur.Random must be the j-th copy node; verify by walking
			copyJ := got
			for k := 0; k < j && copyJ != nil; k++ {
				copyJ = copyJ.Next
			}
			if copyJ != copyCur.Random {
				return false
			}
		} else {
			if copyCur.Random != nil {
				return false
			}
		}
	}
	if copyCur != nil {
		return false
	}

	return true
}

func Test_removeNthFromEnd(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		n        int
		expected []int
	}{
		{"remove_last", []int{1, 2, 3, 4, 5}, 1, []int{1, 2, 3, 4}},
		{"remove_first", []int{1, 2, 3, 4, 5}, 5, []int{2, 3, 4, 5}},
		{"remove_middle", []int{1, 2, 3, 4, 5}, 3, []int{1, 2, 4, 5}},
		{"single_node", []int{1}, 1, []int{}},
		{"two_nodes_remove_first", []int{1, 2}, 2, []int{2}},
		{"two_nodes_remove_last", []int{1, 2}, 1, []int{1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := removeNthFromEnd(utils.CreateLinkedList(tc.input), tc.n)
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))
		})
	}
}

func Test_rotateRight(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		k        int
		expected []int
	}{
		{"k=2", []int{1, 2, 3, 4, 5}, 2, []int{4, 5, 1, 2, 3}},
		{"k=0", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"k=length", []int{1, 2, 3}, 3, []int{1, 2, 3}},
		{"k_gt_length", []int{1, 2, 3}, 5, []int{2, 3, 1}}, // 5 % 3 = 2, last 2 move to front
		{"single_node", []int{1}, 10, []int{1}},
		{"two_nodes", []int{1, 2}, 1, []int{2, 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := rotateRight(utils.CreateLinkedList(tc.input), tc.k)
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))
		})
	}
}

func Test_addTwoNumbers(t *testing.T) {
	testCases := []struct {
		name     string
		l1, l2   []int // digits in reverse order (LSB first)
		expected []int
	}{
		{"no_carry", []int{1, 2, 3}, []int{4, 5, 6}, []int{5, 7, 9}},
		{"with_carry", []int{9, 9, 9}, []int{1}, []int{0, 0, 0, 1}},
		{"carry_at_end", []int{5}, []int{5}, []int{0, 1}},
		{"different_lengths", []int{2, 4, 3}, []int{5, 6, 4}, []int{7, 0, 8}}, // 342 + 465 = 807
		{"zeros", []int{0}, []int{0}, []int{0}},
		{"l2_longer", []int{1}, []int{9, 9}, []int{0, 0, 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := addTwoNumbers(utils.CreateLinkedList(tc.l1), utils.CreateLinkedList(tc.l2))
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))
		})
	}
}

func Test_oddEvenList(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"five_nodes", []int{1, 2, 3, 4, 5}, []int{1, 3, 5, 2, 4}},
		{"four_nodes", []int{2, 1, 3, 5}, []int{2, 3, 1, 5}},
		{"single_node", []int{1}, []int{1}},
		{"two_nodes", []int{1, 2}, []int{1, 2}},
		{"three_nodes", []int{1, 2, 3}, []int{1, 3, 2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := oddEvenList(utils.CreateLinkedList(tc.input))
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))

			got = oddEvenListCalude(utils.CreateLinkedList(tc.input))
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))
		})
	}
}

func Test_reverseKGroup(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		k        int
		expected []int
	}{
		{"k=2", []int{1, 2, 3, 4, 5}, 2, []int{2, 1, 4, 3, 5}},
		{"k=3", []int{1, 2, 3, 4, 5}, 3, []int{3, 2, 1, 4, 5}},
		{"k=1_no_change", []int{1, 2, 3}, 1, []int{1, 2, 3}},
		{"k=length", []int{1, 2, 3, 4}, 4, []int{4, 3, 2, 1}},
		{"single_node", []int{1}, 1, []int{1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := reverseKGroup(utils.CreateLinkedList(tc.input), tc.k)
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))
		})
	}
}
