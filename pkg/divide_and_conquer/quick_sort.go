package divide_and_conquer

import . "gocode/types"

/**
 * Quick Sort
 * T: O(n*log(n)) on average. O(n^2) worst.
 *
 *               partition(start, end)
 *               /          |         \
 * [less than pivot... ], pivot, [greater than pivot... ]
 *
 * recursive structure
 */
func quick_sort(arr []int, start, end int) {
    if start < end {
        pivot_idx := partition(arr, start, end)
        quick_sort(arr, start, pivot_idx-1)
        quick_sort(arr, pivot_idx+1, end)
    }
}

func partition(arr []int, low, high int) int {
    var i = low // pivot is arr[high]
    for j := i; j < high; j++ {
        if arr[j] < arr[high] { // if arr[j] <= arr[high]
            arr[i], arr[j] = arr[j], arr[i]
            i++
        }
    }

    arr[i], arr[high] = arr[high], arr[i]
    return i
}

func partition_asc(arr []int, low, high int) int {
    i := low + 1 // pivot is arr[low]
    for j := i; j <= high; j++ {
        if arr[j] <= arr[low] {
            arr[i], arr[j] = arr[j], arr[i]
            i++
        }
    }

    arr[low], arr[i-1] = arr[i-1], arr[low]
    return i - 1
}

func partition_dec(arr []int, low, high int) int {
    i := low + 1 // pivot is arr[low]
    for j := i; j <= high; j++ {
        if arr[j] > arr[low] {
            arr[i], arr[j] = arr[j], arr[i]
            i++
        }
    }

    arr[low], arr[i-1] = arr[i-1], arr[low]
    return i - 1
}

/**
 * Quick Sort linked list
 *
 * 86. Partition List
 */

 func sortListWithPartition(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return head
    }

    cur := head.Next
    dummySmaller, dummyGreater := ListNode{}, ListNode{}
    curSmaller, curGreater := &dummySmaller, &dummyGreater

    for cur != nil {
        if cur.Val < head.Val {
            curSmaller.Next = cur
            curSmaller = curSmaller.Next
        } else {
            curGreater.Next = cur
            curGreater = curGreater.Next
        }
        cur = cur.Next
    }
    curSmaller.Next = nil
    curGreater.Next = nil

    dummySmaller.Next = sortListWithPartition(dummySmaller.Next)
    dummyGreater.Next = sortListWithPartition(dummyGreater.Next)
    
    cur = dummySmaller.Next
    if cur != nil {
        for cur.Next != nil {
            cur = cur.Next
        }
        cur.Next = head
        head.Next = dummyGreater.Next
        return dummySmaller.Next
    } else {
        head.Next = dummyGreater.Next
        return head
    }
}

func sortList(head *ListNode) *ListNode {
    if head == nil {
        return nil
    }

    tail := head
    for tail.Next != nil {
        tail = tail.Next
    }

    quickSortHelper(head, tail)

    return head
}

func quickSortHelper(head *ListNode, tail *ListNode) {
    if head == nil || head == tail {
        return
    }

    pivot := partitionListSwap(head, tail)

    quickSortHelper(head, pivot)

    quickSortHelper(pivot.Next, tail)
}

func partitionListSwap(head *ListNode, tail *ListNode) *ListNode {
    if head == nil || tail == nil {
        return nil
    }

    pivot := head
    pre, curr := head, head
    for curr != tail.Next {
        if curr.Val < pivot.Val {
            curr.Val, pre.Next.Val = pre.Next.Val, curr.Val
            pre = pre.Next
        }

        curr = curr.Next
    }

    pre.Val, pivot.Val = pivot.Val, pre.Val

    return pre
}
