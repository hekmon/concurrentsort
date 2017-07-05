# ConcurrentSort
Sorting algorithms parallelized in Golang.

Simply doing them as exercise.

## QuickSort

```
Initializing an empty slice with 100000 slots
Filling it up with random numbers
Init done
Start sorting with 1 workers
Sorted in 70.856406ms
Checking slice order...

Initializing an empty slice with 100000 slots
Filling it up with random numbers
Init done
Start sorting with 8 workers
Sorted in 16.557851ms
Checking slice order...

Initializing an empty slice with 1000000 slots
Filling it up with random numbers
Init done
Start sorting with 1 workers
Sorted in 660.800135ms
Checking slice order...

Initializing an empty slice with 1000000 slots
Filling it up with random numbers
Init done
Start sorting with 8 workers
Sorted in 156.058668ms
Checking slice order...

Initializing an empty slice with 10000000 slots
Filling it up with random numbers
Init done
Start sorting with 1 workers
Sorted in 6.435971655s
Checking slice order...

Initializing an empty slice with 10000000 slots
Filling it up with random numbers
Init done
Start sorting with 8 workers
Sorted in 1.456817945s
Checking slice order...

PASS
ok  	_/tmp/cs/env/src/concurrentsort	16.854s
```
