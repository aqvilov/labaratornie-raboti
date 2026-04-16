package main

// первый алгоритм ( наличие элемента в массиве )
func inArray(nums []int, target int) bool {
	for i := 0; i < len(nums); i++ {
		if nums[i] == target {
			return true
		}
	}
	return false
}

// второй максимум
func findSecondMax(arr []int) int {
	if len(arr) < 2 {
		return 0
	}

	max1 := arr[0]
	max2 := arr[0]

	for i := 1; i < len(arr); i++ {
		if arr[i] > max1 {
			max2 = max1
			max1 = arr[i]
		} else if arr[i] > max2 && arr[i] != max1 {
			max2 = arr[i]
		}
	}
	return max2
}

// бинпоиск
func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return -1
}

// таблица умножения
func multi(n int) [][]int {
	table := make([][]int, n)
	for i := 0; i < n; i++ {
		table[i] = make([]int, n)
		for j := 0; j < n; j++ {
			table[i][j] = (i + 1) * (j + 1)
		}
	}
	return table
}
