package backtracking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateParenthesisCalude(t *testing.T) {
	testCases := []struct {
		name     string
		n        int
		expected []string
	}{
		{
			name: "n=3",
			n:    3,
			expected: []string{
				"((()))", "(()())", "(())()", "()(())", "()()()",
			},
		},
		{
			name:     "n=1",
			n:        1,
			expected: []string{"()"},
		},
		{
			name:     "n=2",
			n:        2,
			expected: []string{"(())", "()()"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := generateParenthesisCalude(tc.n)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}
