package concurrentsort

import (
	"sync"
)

// QuickSortMinSizeForConcurrency will prevent starting another goroutine if the slice is too small
var QuickSortMinSizeForConcurrency = 1000 // good value for IntSlice
// Need proper benchmarking to compute the right value vs the cost of a goroutine for each type

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

// QuickSort sorts data using the quicksort algo distributed on nbWorkers goroutines
func QuickSort(data QuickSortable, nbWorkers int) {
	manager := quickSortConcurrentManager{availableWorkers: nbWorkers - 1}
	quickSort(data, &manager)
	manager.rdvpoint.Wait()
}

func quickSort(data QuickSortable, manager *quickSortConcurrentManager) {
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
		// Try to side launch sub slices quicksort
		if firstSlice.Len() >= QuickSortMinSizeForConcurrency && manager.isAWorkerAvailable() {
			go func() {
				defer manager.workerDone()
				quickSort(firstSlice, manager)
			}()
			firstSideLaunched = true
		} else if secondSlice.Len() >= QuickSortMinSizeForConcurrency && manager.isAWorkerAvailable() {
			// If first did not make the side launch, maybe the second will ?
			go func() {
				defer manager.workerDone()
				quickSort(secondSlice, manager)
			}()
			secondSideLaunched = true
		}
		// Sort within the same goroutine
		if !firstSideLaunched {
			quickSort(firstSlice, manager)
		}
		if !secondSideLaunched {
			quickSort(secondSlice, manager)
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
