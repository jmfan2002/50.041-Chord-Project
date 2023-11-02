package util

import "math/big"

// Finds the insertion index of item in arr
func BinarySearch(arr []big.Int, item *big.Int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := low + (high-low)/2
		cmp := arr[mid].Cmp(item)
		if cmp == -1 {
			low = mid + 1
		} else if cmp == 1 {
			high = mid - 1
		} else {
			return mid
		}
	}

	return low
}
