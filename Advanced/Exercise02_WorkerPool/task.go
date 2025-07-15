package exercise02workerpool

import "time"

// Task represents a unit of work to be processed within the worker pool.
// It holds all necessary information for a worker to perform its computation.
type Task struct {
	ID         int           // Unique identifier for the task.
	Data       int           // The input data for the CPU-intensive calculation (e.g., number to check for primality).
	Complexity time.Duration // Simulated duration for processing this specific task.
	Result     any           // Stores the outcome of the task's processing (e.g., boolean for isPrime). 'any' type allows flexibility for different task results.
	Err        error         // Stores any error that occurred during task processing. Nil if successful.
}

// isPrime checks if a given number is prime.
// This function serves as the CPU-intensive calculation that workers will perform.
// It uses an optimized algorithm for primality testing (checking divisibility only by 6k ± 1).
func isPrime(num int) bool {
	// Numbers less than or equal to 1 are not prime by definition.
	if num <= 1 {
		return false
	}

	// 2 and 3 are prime numbers.
	if num == 2 || num == 3 {
		return true
	}

	// If the number is divisible by 2 or 3, it's not prime (unless it's 2 or 3 itself,
	// which are handled above). This also handles all multiples of 2 and 3.
	if num%2 == 0 || num%3 == 0 {
		return false
	}

	// All primes greater than 3 can be expressed in the form 6k ± 1.
	// We only need to check divisors up to the square root of 'num'.
	// We increment 'i' by 6 because we've already checked divisibility by 2 and 3.
	// So, we check 'i' (6k-1) and 'i+2' (6k+1).
	for i := 5; i*i <= num; i = i + 6 {
		// If 'num' is divisible by 'i' or 'i+2', it's not prime.
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
	}
	// If no divisors are found, the number is prime.
	return true
}
