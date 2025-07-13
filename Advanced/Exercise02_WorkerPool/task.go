package exercise02workerpool

import "time"

// Task represents a unit of work to be Processed
type Task struct {
	ID         int
	Data       int
	Complexity time.Duration
	Result     any
	Err        error
}

// isPrice check if a number is prime (CPU-intensive calculation)
func isPrime(num int) bool {
	if num <= 1 {
		return false
	}

	if num == 2 || num == 3 {
		return true
	}

	for i := 5; i*i <= num; i = i + 6 {
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
	}
	return true
}
