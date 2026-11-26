package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateLinkedList(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int // nil means expect a nil head
	}{
		{"empty_input", []int{}, nil},
		{"single_element", []int{42}, []int{42}},
		{"multiple_elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			head := CreateLinkedList(tc.input)
			if tc.expected == nil {
				assert.Nil(t, head)
				return
			}
			cur := head
			for _, v := range tc.expected {
				assert.NotNil(t, cur)
				assert.Equal(t, v, cur.Val)
				cur = cur.Next
			}
			assert.Nil(t, cur) // no extra tail nodes
		})
	}
}

func Test_VerifyLinkedLists(t *testing.T) {
	tests := []struct {
		name     string
		l1, l2   []int
		expected bool
	}{
		{"both_nil", nil, nil, true},
		{"equal_single", []int{1}, []int{1}, true},
		{"equal_multi", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"value_mismatch", []int{1, 2, 3}, []int{1, 2, 4}, false},
		{"l1_shorter", []int{1, 2}, []int{1, 2, 3}, false},
		{"l2_shorter", []int{1, 2, 3}, []int{1, 2}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := VerifyLinkedLists(
				CreateLinkedList(tc.l1),
				CreateLinkedList(tc.l2),
			)
			assert.Equal(t, tc.expected, result)
		})
	}
}
