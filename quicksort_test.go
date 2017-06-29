package concurrentsort

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

const size = 10000000

var toSort IntSlice

func init() {
	fmt.Println("Initializing an empty slice with", size, "slots")
	toSort = make(IntSlice, size)
	fmt.Println("Filling it up with random numbers")
	rand.Seed(time.Now().UnixNano())
	nbCPUs := runtime.NumCPU()
	indexesByCPU := size / nbCPUs
	var fillGroup sync.WaitGroup
	for cpu := 0; cpu < nbCPUs; cpu++ {
		start := indexesByCPU * cpu
		end := indexesByCPU * (cpu + 1)
		if cpu == nbCPUs-1 && end != len(toSort) {
			end = len(toSort)
		}
		fillGroup.Add(1)
		go func(s, e int) {
			defer fillGroup.Done()
			for i := s; i < e; i++ {
				toSort[i] = int(rand.Int())
				// toSort[i] = int(rand.Int31n(20))
			}
		}(start, end)
	}
	fillGroup.Wait()
	fmt.Println("Setup done")
}

func TestQuickSort(t *testing.T) {
	// Sort
	nbWorkers := runtime.NumCPU()
	fmt.Println("Start sorting with", nbWorkers, "workers")
	start := time.Now()
	QuickSort(toSort, nbWorkers)
	fmt.Println("Sorted in", time.Since(start))
	// Check
	fmt.Println("Checking slice order...")
	var previous int
	for i := 1; i < len(toSort); i++ {
		previous = i - 1
		if toSort[previous] > toSort[i] {
			t.Errorf("Error: index %d is greater (%d) than index %d (%d)\n", previous, toSort[previous], i, toSort[i])
		}
	}
}
