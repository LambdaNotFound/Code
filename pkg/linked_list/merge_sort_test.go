package linked_list

import (
	"gocode/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortList(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"unsorted", []int{4, 2, 1, 3}, []int{1, 2, 3, 4}},
		{"reverse_sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"already_sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"duplicates", []int{3, 1, 2, 1, 3}, []int{1, 1, 2, 3, 3}},
		{"single_node", []int{1}, []int{1}},
		{"two_nodes_sorted", []int{1, 2}, []int{1, 2}},
		{"two_nodes_unsorted", []int{2, 1}, []int{1, 2}},
		{"all_same", []int{2, 2, 2}, []int{2, 2, 2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := sortList(utils.CreateLinkedList(tc.input))
			assert.True(t, utils.VerifyLinkedLists(utils.CreateLinkedList(tc.expected), got))
		})
	}
}
