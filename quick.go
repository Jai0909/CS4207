package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func partition(arr []int, low, high int) int {
	pivotIndex := rand.Intn(high-low+1) + low
	pivot := arr[pivotIndex]
	arr[pivotIndex], arr[high] = arr[high], arr[pivotIndex]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return i
}

func quicksort(arr []int, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(arr) <= 1 {
		return
	}

	pivotIndex := partition(arr, 0, len(arr)-1)

	leftWg := &sync.WaitGroup{}
	rightWg := &sync.WaitGroup{}

	leftWg.Add(1)
	go func() {
		quicksort(arr[:pivotIndex], leftWg)
	}()

	rightWg.Add(1)
	go func() {
		quicksort(arr[pivotIndex+1:], rightWg)
	}()

	leftWg.Wait()
	rightWg.Wait()
}

func measureTime(arr []int, sortFunc func([]int, *sync.WaitGroup), wg *sync.WaitGroup) time.Duration {
	start := time.Now()
	sortFunc(arr, wg)
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

		var wg sync.WaitGroup
		wg.Add(1)
		quicksortTime := measureTime(arr, quicksort, &wg)

		exetimes[i] = quicksortTime

		fmt.Printf("Unsorted Array %v\n", arr)
		wg.Wait()
		fmt.Printf("Sorted Array %v\n", arr)
		fmt.Printf("Quicksort %v Core Execution Time: %v\n", i, quicksortTime)
	}
	fmt.Println(exetimes)
}
