# ConcurrentSort
Sorting algorithms parallelized in Golang.

Simply doing them as exercise.

[![GoDoc](https://godoc.org/github.com/Hekmon/concurrentsort?status.svg)](https://godoc.org/github.com/Hekmon/concurrentsort)

## QuickSort

```
Initializing an empty slice with 100000 slots
Filling it up with random numbers
Init done
Start sorting with 1 workers
Sorted in 90.433691ms
Checking slice order...

Initializing an empty slice with 100000 slots
Filling it up with random numbers
Init done
Start sorting with 8 workers
Sorted in 23.142308ms
Checking slice order...

Initializing an empty slice with 1000000 slots
Filling it up with random numbers
Init done
Start sorting with 1 workers
Sorted in 779.117719ms
Checking slice order...

Initializing an empty slice with 1000000 slots
Filling it up with random numbers
Init done
Start sorting with 8 workers
Sorted in 143.434931ms
Checking slice order...

Initializing an empty slice with 10000000 slots
Filling it up with random numbers
Init done
Start sorting with 1 workers
Sorted in 7.277401429s
Checking slice order...

Initializing an empty slice with 10000000 slots
Filling it up with random numbers
Init done
Start sorting with 8 workers
Sorted in 1.407103771s
Checking slice order...

PASS
ok  	_/tmp/cs/env/src/concurrentsort	17.841s
```
