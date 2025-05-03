package backtracking

import (
	"reflect"
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

            result = letterCombinationsBacktrack(tc.digits)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_combinationSum(t *testing.T) {
    testCases := []struct {
        name       string
        candidates []int
        target     int
        expected   [][]int
    }{
        {
            "case 1",
            []int{2, 3, 6, 7},
            7,
            [][]int{
                {2, 2, 3},
                {7},
            },
        },
        {
            "case 2",
            []int{2, 3, 5},
            7,
            [][]int{
                {2, 2, 2, 2},
                {2, 3, 3},
                {3, 5},
            },
        },
        {
            "case 3",
            []int{2},
            1,
            [][]int{},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := combinationSum(tc.candidates, tc.target)
            reflect.DeepEqual(tc.expected, result)
        })
    }
}
