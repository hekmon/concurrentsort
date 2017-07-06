package main

import (
	"concurrentsort"
)

func main() {
	best := QuickSortMinSizeCompute(0, 2, 20, 1<<24, 1<<7, 4)
	concurrentsort.QuickSortMinSizeForConcurrency = best
}
