package divide_and_conquer

import . "gocode/types"

/**
 * Quick Sort
 * T: O(n*log(n)) on average. O(n^2) worst.
 *
 * [less than pivot... ] pivot [greater than pivot... ]
 */
func quick_sort(arr []int, start, end int) {
    if start < end {
        pivot_idx := partition(arr, start, end)
        quick_sort(arr, start, pivot_idx-1)
        quick_sort(arr, pivot_idx+1, end)
    }
}

func partition(arr []int, low, high int) int {
    var i = low // pivot := arr[high]
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
    i := low + 1 // pivot := arr[low]
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
    i := low + 1 // pivot := arr[low]
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

    pivot := partitionList(head, tail)

    quickSortHelper(head, pivot)

    quickSortHelper(pivot.Next, tail)
}

func partitionList(head *ListNode, tail *ListNode) *ListNode {
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
