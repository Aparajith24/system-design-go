package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func countPrimesSequential(start, end int) int {
	count := 0
	for n := start; n < end; n++ {
		if isPrime(n) {
			count++
		}
	}
	return count
}

func countPrimesParallel(start, end, workers int) int {
	var total int64
	var wg sync.WaitGroup

	chunk := (end - start) / workers
	for w := 0; w < workers; w++ {
		chunkStart := start + w*chunk
		chunkEnd := chunkStart + chunk
		if w == workers-1 {
			chunkEnd = end // last worker absorbs any remainder
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			local := countPrimesSequential(s, e)
			atomic.AddInt64(&total, int64(local))
		}(chunkStart, chunkEnd)
	}

	wg.Wait()
	return int(total)
}

func main() {
	const rangeEnd = 3_000_000
	numCPU := runtime.NumCPU()

	fmt.Printf("CPUs available (GOMAXPROCS): %d\n", numCPU)
	fmt.Printf("counting primes in [0, %d)\n\n", rangeEnd)

	start := time.Now()
	seqCount := countPrimesSequential(0, rangeEnd)
	seqElapsed := time.Since(start)
	fmt.Printf("sequential: found %d primes in %s\n", seqCount, seqElapsed)

	start = time.Now()
	parCount := countPrimesParallel(0, rangeEnd, numCPU)
	parElapsed := time.Since(start)
	fmt.Printf("parallel (%d goroutines): found %d primes in %s\n", numCPU, parCount, parElapsed)

	speedup := float64(seqElapsed) / float64(parElapsed)
	fmt.Printf("\nspeedup: %.2fx\n", speedup)
}
