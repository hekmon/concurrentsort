package main

import (
	"runtime"

	"concurrentsort"
)

func main() {
	best := QuickSortMinSizeCompute(20, 5, 80, 1<<24, 1<<7, runtime.NumCPU())
	concurrentsort.QuickSortMinSizeForConcurrency = best
}
