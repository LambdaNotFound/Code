package two_pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_twoSum(t *testing.T) {
    testCases := []struct {
        name     string
        numbers  []int
        target   int
        expected []int
    }{
        {
            "case 1",
            []int{2, 7, 11, 15},
            9,
            []int{1, 2},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := twoSum(tc.numbers, tc.target)
            assert.Equal(t, tc.expected, result)
        })
    }
}
