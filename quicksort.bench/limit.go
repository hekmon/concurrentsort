package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"concurrentsort"
	"sort"
)

// QuickSortMinSizeCompute returns the best usable concurrency slice limit for quicksort.
// Warning: extremly long
func QuickSortMinSizeCompute(startAt, increase, stopAt, sliceSize, nbRun, nbWorkers int) int {
	fmt.Println("* Determining the best limit for concurrency on a slice of size", sliceSize, "with", nbWorkers, "workers")
	type benchdata struct {
		avg    time.Duration
		values []time.Duration
	}
	bench := make(map[int]*benchdata)
	// Placeholders
	var start time.Time
	var total time.Duration
	var currentDuration time.Duration
	var best int
	currentLimit := startAt
	// Start benchmarking
	for currentLimit = startAt; currentLimit <= stopAt; currentLimit += increase {
		fmt.Println()
		fmt.Println("Benchmarking with limit value at", currentLimit)
		// Init
		bench[currentLimit] = new(benchdata)
		bench[currentLimit].values = make([]time.Duration, 0, nbRun)
		// Run benchmark
		for run := 0; run < nbRun; run++ {
			// Init the slice
			fmt.Printf(" %d", run+1)
			toSort := make(concurrentsort.IntSlice, sliceSize)
			rand.Seed(time.Now().UnixNano())
			indexesByWorker := sliceSize / nbWorkers
			var fillGroup sync.WaitGroup
			for worker := 0; worker < nbWorkers; worker++ {
				start := indexesByWorker * worker
				end := indexesByWorker * (worker + 1)
				if worker == nbWorkers-1 && end != len(toSort) {
					end = len(toSort)
				}
				fillGroup.Add(1)
				go func(s, e int) {
					defer fillGroup.Done()
					for i := s; i < e; i++ {
						toSort[i] = rand.Int()
					}
				}(start, end)
			}
			fillGroup.Wait()
			fmt.Print("I")
			// Sort
			start = time.Now()
			concurrentsort.QuickSortCustom(toSort, nbWorkers, currentLimit)
			bench[currentLimit].values = append(bench[currentLimit].values, time.Since(start))
			fmt.Print("S")
		}
		fmt.Println()
		// Compute average
		total = 0
		for _, currentDuration = range bench[currentLimit].values {
			total += currentDuration
		}
		bench[currentLimit].avg = total / time.Duration(nbRun)
		fmt.Println("Average run is", bench[currentLimit].avg, "with a limit of", currentLimit)
	}
	// Search best
	best = startAt
	list := make(sort.IntSlice, 0, len(bench))
	for limit, data := range bench {
		if data.avg < bench[best].avg {
			best = limit
		}
		list = append(list, limit)
	}
	list.Sort()
	// Print results
	fmt.Println()
	fmt.Println("Summary :")
	for _, limit := range list {
		fmt.Printf("%d\t%v", limit, bench[limit].avg)
		if limit == best {
			fmt.Printf(" *")
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("* Best limit is", best, "with", bench[best].avg)
	return best
}
