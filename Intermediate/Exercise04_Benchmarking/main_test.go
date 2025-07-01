package exercise04benchmarking_test

import (
	"testing" // The testing package is required for benchmarks

	exercise04benchmarking "github.com/Daniel-Q-Reis/GoroutinesFromBeginningToAdvanced/Intermediate/Exercise04_Benchmarking"
)

// Define various values of N to test different scales for benchmarks.
// These variables will be used across different benchmark functions.
var (
	N_SMALL  = 1000
	N_MEDIUM = 100000
	N_LARGE  = 1000000
	N_XLARGE = 4000000 // Your original large value
)

// --- Benchmarks for SumPrimesSequential ---

// BenchmarkSumPrimesSequential_Small benchmarks the sequential prime sum for N_SMALL.
// b *testing.B is the benchmark context.
func BenchmarkSumPrimesSequential_Small(b *testing.B) {
	// The loop runs b.N times. b.N is automatically adjusted by the Go testing framework
	// to make sure the benchmark runs for a sufficient amount of time to get reliable results.
	for i := 0; i < b.N; i++ {
		exercise04benchmarking.SumPrimesSequential(N_SMALL) // Call the function to be benchmarked
	}
}

// BenchmarkSumPrimesSequential_Medium benchmarks the sequential prime sum for N_MEDIUM.
func BenchmarkSumPrimesSequential_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exercise04benchmarking.SumPrimesSequential(N_MEDIUM)
	}
}

// BenchmarkSumPrimesSequential_Large benchmarks the sequential prime sum for N_LARGE.
func BenchmarkSumPrimesSequential_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exercise04benchmarking.SumPrimesSequential(N_LARGE)
	}
}

// BenchmarkSumPrimesSequential_XLarge benchmarks the sequential prime sum for N_XLARGE.
func BenchmarkSumPrimesSequential_XLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exercise04benchmarking.SumPrimesSequential(N_XLARGE)
	}
}

// --- Benchmarks for SumPrimesConcurrent ---

// BenchmarkSumPrimesConcurrent_Small benchmarks the concurrent prime sum for N_SMALL.
func BenchmarkSumPrimesConcurrent_Small(b *testing.B) {
	// Optional: You could explicitly set GOMAXPROCS here for the benchmark,
	// but runtime.NumCPU() is generally good as a default in the function itself.
	// For testing scalability, you might use:
	// oldMaxProcs := runtime.GOMAXPROCS(runtime.NumCPU())
	// defer runtime.GOMAXPROCS(oldMaxProcs) // Restore original GOMAXPROCS after benchmark
	for i := 0; i < b.N; i++ {
		exercise04benchmarking.SumPrimesConcurrent(N_SMALL)
	}
}

// BenchmarkSumPrimesConcurrent_Medium benchmarks the concurrent prime sum for N_MEDIUM.
func BenchmarkSumPrimesConcurrent_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exercise04benchmarking.SumPrimesConcurrent(N_MEDIUM)
	}
}

// BenchmarkSumPrimesConcurrent_Large benchmarks the concurrent prime sum for N_LARGE.
func BenchmarkSumPrimesConcurrent_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exercise04benchmarking.SumPrimesConcurrent(N_LARGE)
	}
}

// BenchmarkSumPrimesConcurrent_XLarge benchmarks the concurrent prime sum for N_XLARGE.
func BenchmarkSumPrimesConcurrent_XLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exercise04benchmarking.SumPrimesConcurrent(N_XLARGE)
	}
}

// Example of using b.RunParallel() for benchmarks that benefit from parallelism within the test.
// This is more complex for this specific aggregation problem, but it's a good pattern to know.
// b.RunParallel distributes b.N iterations among multiple goroutines.
// to run this benchmark in terminal execute -> go test -bench=. -benchmem      //for more info check README
