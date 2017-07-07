package concurrentsort

import (
	"runtime"
	"sync"
)

/*
	Interface and common types
*/

// QuickSortable is an interface which must be satisfied in order to call QuickSort()
type QuickSortable interface {
	Len() int
	LessOrEqual(i, j int) bool
	Swap(i, j int)
	GetSubSliceTo(i int) QuickSortable
	GetSubSliceFrom(j int) QuickSortable
}

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type IntSlice []int

func (p IntSlice) Len() int {
	return len(p)
}
func (p IntSlice) LessOrEqual(i, j int) bool {
	return p[i] <= p[j]
}
func (p IntSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p IntSlice) GetSubSliceTo(i int) QuickSortable {
	return p[:i]
}
func (p IntSlice) GetSubSliceFrom(i int) QuickSortable {
	return p[i:]
}

/*
	Quick Sorting
*/

type quickSortConcurrentManager struct {
	availableWorkers int
	access           sync.Mutex
	rdvpoint         sync.WaitGroup
}

func (qscm *quickSortConcurrentManager) isAWorkerAvailable() bool {
	defer qscm.access.Unlock()
	qscm.access.Lock()
	if qscm.availableWorkers > 0 {
		qscm.availableWorkers--
		qscm.rdvpoint.Add(1)
		return true
	}
	return false
}

func (qscm *quickSortConcurrentManager) workerDone() {
	defer qscm.access.Unlock()
	qscm.access.Lock()
	qscm.availableWorkers++
	qscm.rdvpoint.Done()
}

// QuickSort sorts data using the quicksort algo distributed on multiples workers.
// It is just a high level wrapper which set optimal parameters to QuickSortCustom().
// The number of workers will be equal to runtime.NumCPU() and the concurrent limit adapted to the number of workers.
func QuickSort(data QuickSortable) {
	QuickSortCustom(data, runtime.NumCPU(), int(float64(runtime.NumCPU())*1.5)) // 1.5 magic number: check quicksort.bench package
}

// QuickSortCustom sorts data using the quicksort algo distributed on nbWorkers goroutines
// nbWorkers allows to specify the number of max goroutines which will be used for concurrency.
// concurrentLimit set the minimum slice size limit for concurrency. Check quicksort.bench package for more informations.
func QuickSortCustom(data QuickSortable, nbWorkers int, concurrentLimit int) {
	// Init the concurrent manager
	manager := quickSortConcurrentManager{availableWorkers: nbWorkers - 1} // the current goroutine is the first worker
	// Start worker 1
	quickSort(data, concurrentLimit, &manager)
	manager.rdvpoint.Wait()
}

func quickSort(data QuickSortable, sliceMinSize int, manager *quickSortConcurrentManager) {
	// Start sorting
	if data.Len() > 1 {
		// Select a pivot on the current slice
		pivotIndex := quickSortSelectPivot(data)
		// Partition it and recover new pivot index
		pivotIndex = quickSortPartition(data, pivotIndex)
		// Prepare the subslices
		firstSlice := data.GetSubSliceTo(pivotIndex)
		secondSlice := data.GetSubSliceFrom(pivotIndex + 1)
		// Are some of them eligible to concurrency ?
		firstSideLaunched := false
		secondSideLaunched := false
		if firstSlice.Len() >= sliceMinSize && manager.isAWorkerAvailable() {
			go func() {
				defer manager.workerDone()
				quickSort(firstSlice, sliceMinSize, manager)
			}()
			firstSideLaunched = true
		} else if secondSlice.Len() >= sliceMinSize && manager.isAWorkerAvailable() {
			go func() {
				defer manager.workerDone()
				quickSort(secondSlice, sliceMinSize, manager)
			}()
			secondSideLaunched = true
		}
		// Sort within the same goroutine
		if !firstSideLaunched {
			quickSort(firstSlice, sliceMinSize, manager)
		}
		if !secondSideLaunched {
			quickSort(secondSlice, sliceMinSize, manager)
		}
	}
}

func quickSortPartition(data QuickSortable, pivotIndex int) (newPivotIndex int) {
	// Swap the pivot to the end
	lastIndex := data.Len() - 1
	data.Swap(pivotIndex, lastIndex)
	pivotIndex = lastIndex
	// Launch the quicksort on this part of the slice
	var j int
	for i := 0; i < pivotIndex; i++ {
		if data.LessOrEqual(i, pivotIndex) {
			data.Swap(i, j)
			j++
		}
	}
	// Swap the pivot to the index where all data on the left are less of equals to him
	data.Swap(j, pivotIndex)
	// Return the new pivot index (now at index j)
	return j
}

func quickSortSelectPivot(data QuickSortable) (pivot int) {
	// Classical variant: take last index
	return data.Len() - 1
	// Check with median of medians ?
}
