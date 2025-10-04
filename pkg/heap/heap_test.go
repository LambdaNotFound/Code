package heap

import (
	. "gocode/types"
	"gocode/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeKLists(t *testing.T) {
    testCases := []struct {
        name     string
        input    []*ListNode
        expected *ListNode
    }{
        {
            "case 1",
            []*ListNode{
                utils.CreateLinkedList([]int{1, 4, 5}),
                utils.CreateLinkedList([]int{1, 3, 4}),
                utils.CreateLinkedList([]int{2, 6}),
            },
            utils.CreateLinkedList([]int{1, 1, 2, 3, 4, 4, 5, 6}),
        },
        {
            "case 2",
            []*ListNode{},
            utils.CreateLinkedList([]int{}),
        },
        {
            "case 3",
            []*ListNode{
                utils.CreateLinkedList([]int{}),
            },
            utils.CreateLinkedList([]int{}),
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := mergeKLists(tc.input)
            equal := utils.VerifyLinkedLists(tc.expected, result)
            assert.Equal(t, true, equal)
        })
    }
}

func Test_kClosest(t *testing.T) {
    tests := []struct {
        name     string
        points   [][]int
        k        int
        expected [][]int
    }{
        {
            name:     "basic_case",
            points:   [][]int{{1, 3}, {-2, 2}},
            k:        1,
            expected: [][]int{{-2, 2}},
        },
        {
            name:     "two_closest",
            points:   [][]int{{3, 3}, {5, -1}, {-2, 4}},
            k:        2,
            expected: [][]int{{3, 3}, {-2, 4}},
        },
        {
            name:     "all_points",
            points:   [][]int{{2, 2}, {1, 1}, {0, 0}},
            k:        3,
            expected: [][]int{{2, 2}, {1, 1}, {0, 0}}, // any order
        },
        {
            name:     "single_point",
            points:   [][]int{{10, 10}},
            k:        1,
            expected: [][]int{{10, 10}},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := kClosest(tt.points, tt.k)

            // Since order of closest points can vary, we check via sets
            assert.Equal(t, len(tt.expected), len(result))

            // convert slices to maps for comparison
            toMap := func(points [][]int) map[[2]int]bool {
                m := make(map[[2]int]bool)
                for _, p := range points {
                    m[[2]int{p[0], p[1]}] = true
                }
                return m
            }

            gotMap := toMap(result)
            wantMap := toMap(tt.expected)
            assert.Equal(t, wantMap, gotMap)
        })
    }
}