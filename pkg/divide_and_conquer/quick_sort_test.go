package divide_and_conquer

import (
	. "gocode/types"
	"gocode/utils"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_quick_sort(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        expected []int
    }{
        {name: "random_order", input: []int{5, 6, 7, 2, 1, 0}, expected: []int{0, 1, 2, 5, 6, 7}},
        {name: "already_sorted", input: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}},
        {name: "reverse_sorted", input: []int{5, 4, 3, 2, 1}, expected: []int{1, 2, 3, 4, 5}},
        {name: "single_element", input: []int{42}, expected: []int{42}},
        {name: "two_elements", input: []int{3, 1}, expected: []int{1, 3}},
        {name: "with_duplicates", input: []int{3, 1, 2, 1, 3}, expected: []int{1, 1, 2, 3, 3}},
        {name: "all_same", input: []int{7, 7, 7, 7}, expected: []int{7, 7, 7, 7}},
        {name: "negative_numbers", input: []int{-3, -1, -7, 2, 0}, expected: []int{-7, -3, -1, 0, 2}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            arr := make([]int, len(tt.input))
            copy(arr, tt.input)
            quick_sort(arr, 0, len(arr)-1)
            assert.Equal(t, tt.expected, arr)
        })
    }
}

func Test_partition_ascending(t *testing.T) {
    tests := []struct {
        name         string
        input        []int
        low, high    int
        wantIdx      int
        wantArr      []int
    }{
        {
            name: "basic",
            input: []int{7, 3, 4, 6, 5, 5}, low: 0, high: 5,
            wantIdx: 2, wantArr: []int{3, 4, 5, 6, 5, 7},
        },
        {
            name: "pivot_already_smallest",
            input: []int{3, 7, 5, 9, 3}, low: 0, high: 4,
            // pivot = arr[4]=3; nothing < 3, swap arr[0] with arr[4]
            wantIdx: 0, wantArr: []int{3, 7, 5, 9, 3},
        },
        {
            name: "two_elements_swap",
            input: []int{5, 2}, low: 0, high: 1,
            // pivot=2, 5>2 no swap, then arr[0]<->arr[1] → [2,5]
            wantIdx: 0, wantArr: []int{2, 5},
        },
        {
            name: "subrange",
            input: []int{0, 5, 3, 1, 0}, low: 1, high: 4,
            // pivot=arr[4]=0; 5>0,3>0,1>0 → i stays at 1; swap arr[1]<->arr[4]
            wantIdx: 1, wantArr: []int{0, 0, 3, 1, 5},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            arr := make([]int, len(tt.input))
            copy(arr, tt.input)
            got := partition_ascending(arr, tt.low, tt.high)
            assert.Equal(t, tt.wantIdx, got)
            assert.Equal(t, tt.wantArr, arr)
        })
    }
}

func Test_partition_decending(t *testing.T) {
    tests := []struct {
        name      string
        input     []int
        low, high int
        wantIdx   int
        wantArr   []int
    }{
        {
            name: "basic",
            input: []int{7, 3, 4, 6, 5, 5}, low: 0, high: 5,
            wantIdx: 2, wantArr: []int{7, 6, 5, 3, 5, 4},
        },
        {
            name: "two_elements_no_swap",
            input: []int{2, 5}, low: 0, high: 1,
            // pivot=arr[1]=5; arr[0]=2 < 5, not > pivot, i stays 0; swap arr[0]<->arr[1]
            wantIdx: 0, wantArr: []int{5, 2},
        },
        {
            name: "pivot_already_largest",
            input: []int{3, 1, 2, 9}, low: 0, high: 3,
            // pivot=arr[3]=9; 3<9,1<9,2<9 none > pivot; i stays 0; swap arr[0]<->arr[3]
            wantIdx: 0, wantArr: []int{9, 1, 2, 3},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            arr := make([]int, len(tt.input))
            copy(arr, tt.input)
            got := partition_decending(arr, tt.low, tt.high)
            assert.Equal(t, tt.wantIdx, got)
            assert.Equal(t, tt.wantArr, arr)
        })
    }
}

func Test_partition_asc(t *testing.T) {
    tests := []struct {
        name      string
        input     []int
        low, high int
        wantIdx   int
        wantArr   []int
    }{
        {
            name: "basic",
            input: []int{7, 3, 4, 6, 5, 5}, low: 0, high: 5,
            wantIdx: 5, wantArr: []int{5, 3, 4, 6, 5, 7},
        },
        {
            name: "pivot_smallest",
            input: []int{1, 4, 3, 2}, low: 0, high: 3,
            // pivot=arr[0]=1; j from 1..3: 4>1 no,3>1 no,2>1 no; i stays 1; swap arr[0]<->arr[0]
            wantIdx: 0, wantArr: []int{1, 4, 3, 2},
        },
        {
            name: "two_elements",
            input: []int{3, 1}, low: 0, high: 1,
            // pivot=arr[0]=3; j=1: arr[1]=1<=3 swap arr[1]<->arr[1], i=2; swap arr[0]<->arr[1]
            wantIdx: 1, wantArr: []int{1, 3},
        },
        {
            name: "all_smaller",
            input: []int{5, 1, 2, 3, 4}, low: 0, high: 4,
            // pivot=5; all <=5; i ends at 5 (out of range), swap arr[0]<->arr[4]
            wantIdx: 4, wantArr: []int{4, 1, 2, 3, 5},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            arr := make([]int, len(tt.input))
            copy(arr, tt.input)
            got := partition_asc(arr, tt.low, tt.high)
            assert.Equal(t, tt.wantIdx, got)
            assert.Equal(t, tt.wantArr, arr)
        })
    }
}

func Test_partition_dec(t *testing.T) {
    tests := []struct {
        name      string
        input     []int
        low, high int
        wantIdx   int
        wantArr   []int
    }{
        {
            name: "basic",
            input: []int{7, 3, 4, 6, 5, 5}, low: 0, high: 5,
            wantIdx: 0, wantArr: []int{7, 3, 4, 6, 5, 5},
        },
        {
            name: "pivot_smallest",
            input: []int{1, 4, 3, 2}, low: 0, high: 3,
            // pivot=arr[0]=1; i starts at 1; j=1: arr[1]=4>1 swap arr[1]<->arr[1] i=2;
            // j=2: arr[2]=3>1 swap arr[2]<->arr[2] i=3; j=3: arr[3]=2>1 swap arr[3]<->arr[3] i=4;
            // swap arr[0]<->arr[3]: {2,4,3,1}; return 3
            wantIdx: 3, wantArr: []int{2, 4, 3, 1},
        },
        {
            name: "two_elements_descending",
            input: []int{1, 3}, low: 0, high: 1,
            // pivot=1; j=1: 3>1 swap arr[1]<->arr[1],i=2; swap arr[0]<->arr[1]
            wantIdx: 1, wantArr: []int{3, 1},
        },
        {
            name: "already_descending_pivot_largest",
            input: []int{5, 4, 3, 2, 1}, low: 0, high: 4,
            // pivot=5; all < 5; none > pivot; i stays 1; swap arr[0]<->arr[0]
            wantIdx: 0, wantArr: []int{5, 4, 3, 2, 1},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            arr := make([]int, len(tt.input))
            copy(arr, tt.input)
            got := partition_dec(arr, tt.low, tt.high)
            assert.Equal(t, tt.wantIdx, got)
            assert.Equal(t, tt.wantArr, arr)
        })
    }
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
        {
            name:     "single_node",
            list:     []int{1},
            expected: []int{1},
        },
        {
            name:     "already_sorted",
            list:     []int{1, 2, 3, 4},
            expected: []int{1, 2, 3, 4},
        },
        {
            name:     "reverse_sorted",
            list:     []int{4, 3, 2, 1},
            expected: []int{1, 2, 3, 4},
        },
        {
            name:     "with_duplicates",
            list:     []int{3, 1, 2, 1, 3},
            expected: []int{1, 1, 2, 3, 3},
        },
        {
            name:     "all_same",
            list:     []int{5, 5, 5},
            expected: []int{5, 5, 5},
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
        {
            name:     "single_node",
            list:     []int{9},
            expected: []int{9},
        },
        {
            name:     "already_sorted",
            list:     []int{1, 2, 3, 4, 5},
            expected: []int{1, 2, 3, 4, 5},
        },
        {
            name:     "reverse_sorted",
            list:     []int{5, 4, 3, 2, 1},
            expected: []int{1, 2, 3, 4, 5},
        },
        {
            name:     "with_duplicates",
            list:     []int{2, 4, 2, 1, 3},
            expected: []int{1, 2, 2, 3, 4},
        },
        {
            name:     "mixed_negatives",
            list:     []int{10, -3, 7, -1, 2},
            expected: []int{-3, -1, 2, 7, 10},
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
    if head == nil {
        return []int{}
    }
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