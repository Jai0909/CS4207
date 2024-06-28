package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
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

func mergeSort(arr []int, ch chan []int) {
	if len(arr) <= 1 {
		ch <- arr
		return
	}

	mid := len(arr) / 2

	leftCh := make(chan []int)
	rightCh := make(chan []int)

	go func() {
		mergeSort(arr[:mid], leftCh)
	}()

	go func() {
		mergeSort(arr[mid:], rightCh)
	}()

	left := <-leftCh
	right := <-rightCh

	close(leftCh)
	close(rightCh)

	ch <- merge(left, right)
}

func measureTime(arr []int, sortFunc func([]int, chan []int), ch chan []int) time.Duration {
	start := time.Now()
	sortFunc(arr, ch)
	return time.Since(start)
}

func generateRandomArray(size, min, max int) []int {
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rand.Intn(max-min+1) + min
	}
	return arr
}

func main() {
	numcores := runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Number of Cores: %v\n", numcores)

	exetimes := make(map[int]time.Duration)

	for i := 1; i <= numcores; i++ {
		arr := generateRandomArray(100000, -99999, 99999)
		runtime.GOMAXPROCS(i)

		ch := make(chan []int)
		go mergeSort(arr, ch)

		mergeSortTime := measureTime(arr, func(arr []int, ch chan []int) { go mergeSort(arr, ch) }, ch)

		exetimes[i] = mergeSortTime
		
		fmt.Printf("Unsorted Array %v\n", arr)
		ms := <-ch
		fmt.Printf("Sorted Array %v\n", ms)
		fmt.Printf("Merge Sort %v Core Execution Time: %v\n", i, mergeSortTime)
	}
	fmt.Println(exetimes)
}
