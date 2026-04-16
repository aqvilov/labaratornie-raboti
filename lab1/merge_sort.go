package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func mergeSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func generateRandomArray(size int) []int { // рандомные числа в массив
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rand.Intn(10000)
	}
	return arr
}

func testSort(size int) {
	arr := generateRandomArray(size)
	// замеры
	start := time.Now()
	mergeSort(arr)
	elapsed := time.Since(start)

	fmt.Printf("Размер: %d, Время: %v\n", size, elapsed)
}

func memoryTest(size int) {
	arr := generateRandomArray(size)

	var m1, m2 runtime.MemStats
	runtime.GC() // очищаем все что больше не используется
	runtime.ReadMemStats(&m1)

	mergeSort(arr)

	runtime.ReadMemStats(&m2)

	allocKB := (m2.TotalAlloc - m1.TotalAlloc) / 1024

	fmt.Printf("Память: %d KB\n", allocKB)
}

func main() {
	sizes := []int{100, 1000, 5000, 10000}

	for _, size := range sizes {
		testSort(size)
		memoryTest(size)
	}
}
