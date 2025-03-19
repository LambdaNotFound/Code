package divide_and_conquer

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
