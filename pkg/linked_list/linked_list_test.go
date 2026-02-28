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
		expected *ListNode
	}{
		{
			"case 1",
			[]int{1, 2, 3, 4, 5},
			utils.CreateLinkedList([]int{5, 4, 3, 2, 1}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			head := utils.CreateLinkedList(tc.input)
			result := reverseList(head)
			equal := utils.VerifyLinkedLists(tc.expected, result)
			assert.Equal(t, true, equal)

			head = utils.CreateLinkedList(tc.input)
			result = reverseList_recursive(head)
			equal = utils.VerifyLinkedLists(tc.expected, result)
			assert.Equal(t, true, equal)

			head = utils.CreateLinkedList(tc.input)
			result = reverseList_recursive2(head)
			equal = utils.VerifyLinkedLists(tc.expected, result)
			assert.Equal(t, true, equal)
		})
	}
}

func Test_mergeTwoLists(t *testing.T) {
	list1 := utils.CreateLinkedList([]int{1, 3, 5, 7})
	list2 := utils.CreateLinkedList([]int{2, 4, 6, 8})

	want := utils.CreateLinkedList([]int{1, 2, 3, 4, 5, 6, 7, 8})
	got := mergeTwoLists(list1, list2)
	equal := utils.VerifyLinkedLists(want, got)

	assert.Equal(t, true, equal)
}

func Test_swapPairs(t *testing.T) {
	list1 := utils.CreateLinkedList([]int{1, 2, 3, 4})

	want := utils.CreateLinkedList([]int{2, 1, 4, 3})
	got := swapPairs(list1)
	equal := utils.VerifyLinkedLists(want, got)

	assert.Equal(t, true, equal)
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
