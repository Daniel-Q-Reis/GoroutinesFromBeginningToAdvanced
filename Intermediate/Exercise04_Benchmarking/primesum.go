package exercise04benchmarking

import (
	// For formatted I/O (printing to console)
	"math"    // For math.Sqrt function
	"runtime" // For runtime.NumCPU to get the number of logical CPUs
	"sync"    // For synchronization primitives like sync.Mutex and sync.WaitGroup
	// For measuring elapsed time (used in main for demonstration)
)

// isPrime checks if a given number is prime.
// This function is a core component and its efficiency directly impacts overall performance.
func isPrime(num int) bool {
	// Numbers less than or equal to 1 are not prime.
	if num <= 1 {
		return false
	}
	// 2 is the only even prime number.
	if num == 2 {
		return true
	}
	// All other even numbers (greater than 2) are not prime.
	if num%2 == 0 {
		return false
	}

	// For odd numbers, we only need to check for divisibility by odd numbers
	// up to its square root.
	// Converting num to float64 for math.Sqrt, then casting back to int.
	sqrt := int(math.Sqrt(float64(num)))

	// Start checking from 3 and increment by 2 (skipping even numbers).
	for i := 3; i <= sqrt; i += 2 {
		if num%i == 0 {
			return false // Found a divisor, so it's not prime
		}
	}
	return true // No divisors found, so it's prime
}

// sumPrimesInRange calculates the sum of prime numbers within a specified range [start, end].
// It updates a shared totalSum using a mutex for thread-safety.
func sumPrimesInRange(start, end int, totalSum *int64, mu *sync.Mutex) {
	localSum := int64(0) // Accumulate sum locally first to minimize mutex contention

	// Ensure the starting number is at least 2, as primes start from 2.
	if start < 2 {
		start = 2
	}

	// Iterate through the given range and sum up prime numbers.
	for i := start; i <= end; i++ {
		if isPrime(i) {
			localSum += int64(i) // Add to local sum if prime
		}
	}

	// Acquire the mutex to safely update the shared totalSum.
	mu.Lock()
	*totalSum += localSum // Add the local sum to the global sum
	mu.Unlock()           // Release the mutex
}

// SumPrimesSequential calculates the sum of all prime numbers up to 'n' sequentially.
func SumPrimesSequential(n int) int64 {
	sum := int64(0) // Initialize sum (using int64 to accommodate large sums)
	// Iterate from 2 up to n, checking each number for primality.
	for i := 2; i <= n; i++ {
		if isPrime(i) {
			sum += int64(i) // Add to sum if prime
		}
	}
	return sum
}

// SumPrimesConcurrent calculates the sum of all prime numbers up to 'n' concurrently
// using a work-pooling pattern with a mutex for aggregation.
func SumPrimesConcurrent(n int) int64 {
	var totalSum int64    // Shared variable to store the final sum
	var mu sync.Mutex     // Mutex to protect access to totalSum
	var wg sync.WaitGroup // WaitGroup to wait for all worker goroutines to complete

	// Determine the number of workers based on the number of logical CPUs available.
	numWorkers := runtime.NumCPU()
	if numWorkers == 0 { // Fallback in case NumCPU returns 0 for some reason
		numWorkers = 1
	}

	// Calculate the approximate size of each segment (range) for workers.
	// Using int64 for division to handle potentially large 'n'.
	segmentSize := n / numWorkers

	// Launch goroutines (workers) to process segments of the number range.
	for i := 0; i < numWorkers; i++ {
		start := i * segmentSize
		// End of the current segment. Subtracting 1 to ensure disjoint ranges.
		end := start + segmentSize - 1

		// For the last worker, adjust its 'end' to cover up to 'n'
		// in case 'n' is not perfectly divisible by 'numWorkers'.
		if i == numWorkers-1 {
			end = n
		}

		// Ensure the starting point of any segment is at least 2,
		// as primes start from 2. The isPrime function already handles values <= 1.
		if start < 2 {
			start = 2
		}

		wg.Add(1) // Increment the WaitGroup counter for each new goroutine
		// Launch a goroutine, passing the segment's start and end.
		// These parameters are passed by value to avoid closure issues.
		go func(s, e int) {
			defer wg.Done()                        // Decrement WaitGroup counter when the goroutine finishes
			sumPrimesInRange(s, e, &totalSum, &mu) // Call the helper function for this segment
		}(start, end)
	}

	wg.Wait()       // Wait for all worker goroutines to complete their execution
	return totalSum // Return the final accumulated sum
}
