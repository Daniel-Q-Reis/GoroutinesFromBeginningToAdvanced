package main

import (
	"fmt"
	"time"

	exercise04benchmarking "github.com/Daniel-Q-Reis/GoroutinesFromBeginningToAdvanced/Intermediate/Exercise04_Benchmarking"
)

func main() {
	n := 4000000

	fmt.Printf("Calculating primes up to %d...\n", n)

	// --- Sequential Calculation ---
	startTimeSeq := time.Now()
	sumSeq := exercise04benchmarking.SumPrimesSequential(n)
	elapsedTimeSeq := time.Since(startTimeSeq)
	fmt.Printf("Sequential sum of primes up to %d: %d (Elapsed: %s)\n", n, sumSeq, elapsedTimeSeq)

	// --- Concurrent Calculation ---
	startTimeConc := time.Now()
	sumConc := exercise04benchmarking.SumPrimesConcurrent(n)
	elapsedTimeConc := time.Since(startTimeConc)
	fmt.Printf("Concurrent sum of primes up to %d: %d (Elapsed: %s)\n", n, sumConc, elapsedTimeConc)

	// --- Result Verification ---
	if sumSeq == sumConc {
		fmt.Println("Results match! Both sequential and concurrent methods produced the same sum.")
	} else {
		fmt.Println("WARNING: Results DO NOT match between sequential and concurrent methods!")
		fmt.Printf("Sequential: %d, Concurrent: %d\n", sumSeq, sumConc)
	}
}
