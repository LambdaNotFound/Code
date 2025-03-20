package divide_and_conquer

/*
 * Quick Sort,
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

// pivot := arr[high]
func partition(arr []int, low, high int) int {
	var i = low
	for j := i; j < high; j++ {
		if arr[j] < arr[high] { // if arr[j] <= arr[high]
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[i], arr[high] = arr[high], arr[i]
	return i
}

// pivot := arr[low]
func partition_asc(arr []int, low, high int) int {
	i := low + 1
	for j := i; j <= high; j++ {
		if arr[j] <= arr[low] {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[low], arr[i-1] = arr[i-1], arr[low]
	return i - 1
}

// pivot := arr[low]
func partition_dec(arr []int, low, high int) int {
	i := low + 1
	for j := i; j <= high; j++ {
		if arr[j] > arr[low] {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[low], arr[i-1] = arr[i-1], arr[low]
	return i - 1
}
