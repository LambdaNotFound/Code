package linked_list

import (
	. "gocode/types"
	"gocode/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_reorderList(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"odd_length", []int{1, 2, 3, 4, 5}, []int{1, 5, 2, 4, 3}},
		{"even_length", []int{1, 2, 3, 4}, []int{1, 4, 2, 3}},
		{"two_nodes", []int{1, 2}, []int{1, 2}},
		{"single_node", []int{1}, []int{1}},
		{"three_nodes", []int{1, 2, 3}, []int{1, 3, 2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			head := utils.CreateLinkedList(tc.input)
			reorderList(head) // in-place
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), head))
		})
	}
}

func Test_isPalindrome(t *testing.T) {
	testCases := []struct {
		name     string
		build    func() *ListNode
		expected bool
	}{
		{
			name: "palindrome_odd",
			build: func() *ListNode {
				return utils.CreateLinkedList([]int{1, 2, 1})
			},
			expected: true,
		},
		{
			name: "palindrome_even",
			build: func() *ListNode {
				return utils.CreateLinkedList([]int{1, 2, 2, 1})
			},
			expected: true,
		},
		{
			name: "not_palindrome",
			build: func() *ListNode {
				return utils.CreateLinkedList([]int{1, 2, 3})
			},
			expected: false,
		},
		{
			name: "single_node",
			build: func() *ListNode {
				return &ListNode{Val: 5}
			},
			expected: true,
		},
		{
			name: "two_nodes_palindrome",
			build: func() *ListNode {
				return utils.CreateLinkedList([]int{1, 1})
			},
			expected: true,
		},
		{
			name: "two_nodes_not_palindrome",
			build: func() *ListNode {
				return utils.CreateLinkedList([]int{1, 2})
			},
			expected: false,
		},
		{
			name: "longer_palindrome",
			build: func() *ListNode {
				return utils.CreateLinkedList([]int{1, 2, 3, 2, 1})
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, isPalindrome(tc.build()))
		})
	}
}
