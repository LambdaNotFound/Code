package backtracking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_letterCombinations(t *testing.T) {
    testCases := []struct {
        name     string
        digits   string
        expected []string
    }{
        {
            "case 1",
            "23",
            []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"},
        },
        {
            "case 2",
            "2",
            []string{"a", "b", "c"},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := letterCombinations(tc.digits)
            assert.Equal(t, tc.expected, result)
        })
    }
}
