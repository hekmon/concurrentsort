package concurrentsort

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func qsInit(toSort *IntSlice, sliceSize int) {
	fmt.Println("Initializing an empty slice with", sliceSize, "slots")
	*toSort = make(IntSlice, sliceSize)
	fmt.Println("Filling it up with random numbers")
	rand.Seed(time.Now().UnixNano())
	nbCPUs := runtime.NumCPU()
	indexesByCPU := sliceSize / nbCPUs
	var fillGroup sync.WaitGroup
	for cpu := 0; cpu < nbCPUs; cpu++ {
		start := indexesByCPU * cpu
		end := indexesByCPU * (cpu + 1)
		if cpu == nbCPUs-1 && end != len(*toSort) {
			end = len(*toSort)
		}
		fillGroup.Add(1)
		go func(s, e int) {
			defer fillGroup.Done()
			for i := s; i < e; i++ {
				(*toSort)[i] = rand.Int()
			}
		}(start, end)
	}
	fillGroup.Wait()
	fmt.Println("Init done")
}

func qsLaunch(sliceSize, nbWorkers int) error {
	// Create the slice placeholder
	var toSort IntSlice
	// Init the slice
	qsInit(&toSort, sliceSize)
	// Sort
	fmt.Println("Start sorting with", nbWorkers, "workers")
	start := time.Now()
	QuickSort(toSort, nbWorkers, nil)
	fmt.Println("Sorted in", time.Since(start))
	// Check
	fmt.Println("Checking slice order...")
	var previous int
	for i := 1; i < len(toSort); i++ {
		previous = i - 1
		if toSort[previous] > toSort[i] {
			return fmt.Errorf("value at index %d is greater (%d) than value at index %d (%d)", previous, toSort[previous], i, toSort[i])
		}
	}
	return nil
}

func TestQuickSort100000(t *testing.T) {
	var err error
	size := 100000
	if err = qsLaunch(size, 1); err != nil {
		t.Error(err.Error())
	}
	fmt.Println()
	if err = qsLaunch(size, runtime.NumCPU()); err != nil {
		t.Error(err.Error())
	}
	fmt.Println()
}

func TestQuickSort1000000(t *testing.T) {
	var err error
	size := 1000000
	if err = qsLaunch(size, 1); err != nil {
		t.Error(err.Error())
	}
	fmt.Println()
	if err = qsLaunch(size, runtime.NumCPU()); err != nil {
		t.Error(err.Error())
	}
	fmt.Println()
}

func TestQuickSort10000000(t *testing.T) {
	var err error
	size := 10000000
	if err = qsLaunch(size, 1); err != nil {
		t.Error(err.Error())
	}
	fmt.Println()
	if err = qsLaunch(size, runtime.NumCPU()); err != nil {
		t.Error(err.Error())
	}
	fmt.Println()
}
