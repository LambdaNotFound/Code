package two_pointers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxArea(t *testing.T) {
    testCases := []struct {
        name     string
        height   []int
        expected int
    }{
        {
            "case 1",
            []int{1, 8, 6, 2, 5, 4, 8, 3, 7},
            49,
        },
        {
            "case 2",
            []int{1, 1},
            1,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := maxArea(tc.height)
            assert.Equal(t, tc.expected, result)
        })
    }
}

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
        {
            "case 2",
            []int{2, 3, 4},
            9,
            nil,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := twoSum(tc.numbers, tc.target)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_threeSum(t *testing.T) {
    testCases := []struct {
        name     string
        numbers  []int
        expected [][]int
    }{
        {
            "case 1",
            []int{-1, 0, 1, 2, -1, -4},
            [][]int{
                {-1, -1, 2}, {-1, 0, 1},
            },
        },
        {
            "case 2",
            []int{0, 1, 1},
            [][]int{},
        },
        {
            "case 3",
            []int{0, 0, 0},
            [][]int{
                {0, 0, 0},
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := threeSum(tc.numbers)
            assert.ElementsMatch(t, tc.expected, result)
        })
    }
}

func Test_sortColors(t *testing.T) {
    testCases := []struct {
        name     string
        numbers  []int
        expected []int
    }{
        {
            "case 1",
            []int{2, 0, 2, 1, 1, 0},
            []int{0, 0, 1, 1, 2, 2},
        },
        {
            "case 2",
            []int{2, 0, 1},
            []int{0, 1, 2},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            sortColors(tc.numbers)
            assert.Equal(t, tc.expected, tc.numbers)
        })
    }
}

func Test_moveZeroes(t *testing.T) {
    testCases := []struct {
        name     string
        numbers  []int
        expected []int
    }{
        {
            "case 1",
            []int{0, 1, 0, 3, 12},
            []int{1, 3, 12, 0, 0},
        },
        {
            "case 2",
            []int{0},
            []int{0},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            moveZeroes(tc.numbers)
            assert.Equal(t, tc.expected, tc.numbers)
        })
    }
}

func Test_removeElement(t *testing.T) {
    tests := []struct {
        name     string
        nums     []int
        val      int
        expected []int
        length   int
    }{
        {name: "basic_case", nums: []int{3, 2, 2, 3}, val: 3, expected: []int{2, 2}, length: 2},
        {name: "no_removal", nums: []int{1, 2, 4}, val: 3, expected: []int{1, 2, 4}, length: 3},
        {name: "remove_all", nums: []int{1, 1, 1}, val: 1, expected: []int{}, length: 0},
        {name: "mixed", nums: []int{0, 1, 2, 2, 3, 0, 4, 2}, val: 2, expected: []int{0, 1, 3, 0, 4}, length: 5},
        {name: "empty", nums: []int{}, val: 1, expected: []int{}, length: 0},
        {name: "single_match", nums: []int{5}, val: 5, expected: []int{}, length: 0},
        {name: "single_no_match", nums: []int{5}, val: 3, expected: []int{5}, length: 1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            numsCopy := append([]int(nil), tt.nums...) // avoid mutating original
            gotLen := removeElement(numsCopy, tt.val)

            if gotLen != tt.length {
                t.Errorf("removeElement(%v, %d) length = %d; want %d",
                    tt.nums, tt.val, gotLen, tt.length)
            }

            gotSlice := numsCopy[:gotLen]
            if !reflect.DeepEqual(gotSlice, tt.expected) {
                t.Errorf("removeElement(%v, %d) slice = %v; want %v",
                    tt.nums, tt.val, gotSlice, tt.expected)
            }
        })
    }
}

func Test_isPalindrome(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected bool
    }{
        {
            "case 1",
            "A man, a plan, a canal: Panama",
            true,
        },
        {
            "case 2",
            "race a car",
            false,
        },
        {
            "case 3",
            " ",
            true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            isPalindrome(tc.input)
            assert.Equal(t, tc.expected, tc.expected)
        })
    }
}
