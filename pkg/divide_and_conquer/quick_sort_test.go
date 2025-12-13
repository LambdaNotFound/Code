package divide_and_conquer

import (
	. "gocode/types"
	"gocode/utils"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_quick_sort(t *testing.T) {
    array := []int{5, 6, 7, 2, 1, 0}

    expected := []int{0, 1, 2, 5, 6, 7}

    quick_sort(array, 0, len(array)-1)
    assert.Equal(t, expected, array)
}

func Test_partition_ascending(t *testing.T) {
    array := []int{7, 3, 4, 6, 5, 5}
    // 3, 7, 4, 6, 5, 5
    //                ^pivot
    // 3, 4, 7, 6, 5, 5
    // 3, 4, 5, 6, 5, 7

    want := 2
    expected := []int{3, 4, 5, 6, 5, 7}
    got := partition_ascending(array, 0, len(array)-1)

    assert.Equal(t, want, got)
    assert.Equal(t, expected, array)
}

func Test_partition_decending(t *testing.T) {
    array := []int{7, 3, 4, 6, 5, 5}
    // 7, 6, 4, 3, 5, 5
    //                ^pivot
    // 7, 6, 4, 3, 5, 5
    // 7, 6, 5, 3, 5, 4

    want := 2
    expected := []int{7, 6, 5, 3, 5, 4}
    got := partition_decending(array, 0, len(array)-1)

    assert.Equal(t, want, got)
    assert.Equal(t, expected, array)
}

func Test_partition_asc(t *testing.T) {
    array := []int{7, 3, 4, 6, 5, 5}
    // 7, 3, 4, 6, 5, 5
    // ^pivot
    // 5, 3, 4, 6, 5, 7

    want := 5
    expected := []int{5, 3, 4, 6, 5, 7}

    got := partition_asc(array, 0, len(array)-1)

    assert.Equal(t, want, got)
    assert.Equal(t, expected, array)
}

func Test_partition_dec(t *testing.T) {
    array := []int{7, 3, 4, 6, 5, 5}
    // 7, 6, 4, 3, 5, 5
    // ^pivot
    // 7, 6, 4, 3, 5, 5

    want := 0
    expected := []int{7, 3, 4, 6, 5, 5}

    got := partition_dec(array, 0, len(array)-1)

    assert.Equal(t, want, got)
    assert.Equal(t, expected, array)
}

func Test_sortListQuickSort(t *testing.T) {
    testCases := []struct {
        name     string
        list     []int
        expected []int
    }{
        {
            "case 1",
            []int{4, 2, 1, 3},
            []int{1, 2, 3, 4},
        },
        {
            "case 2",
            []int{-1, 5, 3, 4, 0},
            []int{-1, 0, 3, 4, 5},
        },
        {
            "case 3",
            []int{},
            []int{},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            list := utils.CreateLinkedList(tc.list)
            result := sortListQuickSort(list)
            expected := utils.CreateLinkedList(tc.expected)
            isEqual := utils.VerifyLinkedLists(result, expected)
            assert.Equal(t, true, isEqual)
        })
    }
}

func Test_sortListWithCopy(t *testing.T) {
    testCases := []struct {
        name     string
        list     []int
        expected []int
    }{
        {
            "case 1",
            []int{4, 2, 1, 3},
            []int{1, 2, 3, 4},
        },
        {
            "case 2",
            []int{-1, 5, 3, 4, 0},
            []int{-1, 0, 3, 4, 5},
        },
        {
            "case 3",
            []int{},
            []int{},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            list := utils.CreateLinkedList(tc.list)
            result := sortListWithCopy(list)
            expected := utils.CreateLinkedList(tc.expected)
            isEqual := utils.VerifyLinkedLists(result, expected)
            assert.Equal(t, true, isEqual)
        })
    }
}

func buildList(nums []int) *ListNode {
    dummy := &ListNode{}
    curr := dummy
    for _, n := range nums {
        curr.Next = &ListNode{Val: n}
        curr = curr.Next
    }
    return dummy.Next
}

func listToSlice(head *ListNode) []int {
    var res []int
    for head != nil {
        res = append(res, head.Val)
        head = head.Next
    }
    return res
}

func TestSortList(t *testing.T) {
    testCases := []struct {
        name     string
        input    []int
        expected []int
    }{
        {
            name:     "empty list",
            input:    []int{},
            expected: []int{},
        },
        {
            name:     "single node",
            input:    []int{1},
            expected: []int{1},
        },
        {
            name:     "already sorted",
            input:    []int{1, 2, 3, 4},
            expected: []int{1, 2, 3, 4},
        },
        {
            name:     "reverse sorted",
            input:    []int{4, 3, 2, 1},
            expected: []int{1, 2, 3, 4},
        },
        {
            name:     "with duplicates",
            input:    []int{4, 2, 1, 3, 2},
            expected: []int{1, 2, 2, 3, 4},
        },
        {
            name:     "negative numbers",
            input:    []int{-1, 5, 3, 4, 0},
            expected: []int{-1, 0, 3, 4, 5},
        },
        {
            name:     "mixed positives and negatives",
            input:    []int{10, -3, 7, 2, -1},
            expected: []int{-3, -1, 2, 7, 10},
        },
        {
            name:     "all equal",
            input:    []int{5, 5, 5, 5},
            expected: []int{5, 5, 5, 5},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            head := buildList(tc.input)
            sorted := sortListQuickSort(head)
            got := listToSlice(sorted)

            assert.Equal(t, tc.expected, got)
        })
    }
}