# Intermediate/Exercise04 - Benchmarking Go Code

## Objective

This exercise aims to provide a practical understanding and application of Go's benchmarking techniques. The goal is to measure and compare the performance of different implementations for the same functionality, identify performance bottlenecks, and make informed decisions regarding code optimization.

## Problem Statement

Calculate the sum of all prime numbers up to a given limit `N`.

## Implementations

Two distinct approaches were implemented to solve the prime summation problem:

1.  **`SumPrimesSequential(n int) int64`**:
    This function calculates the sum of primes in a sequential (non-concurrent) manner. It iterates from 2 up to `N`, checking each number for primality using an optimized `isPrime` helper function (trial division method, skipping even numbers after 2).

2.  **`SumPrimesConcurrent(n int) int64`**:
    This function calculates the sum of primes concurrently. The problem range `[2, N]` is divided into smaller segments. Each segment is assigned to a goroutine (worker). Each worker then calculates the sum of primes within its assigned segment. A `sync.Mutex` is used to safely aggregate these partial sums into a shared `totalSum` variable, preventing race conditions. A `sync.WaitGroup` ensures that the main goroutine waits for all worker goroutines to complete before returning the final sum.

## `isPrime` Helper Function

Both `SumPrimesSequential` and `SumPrimesConcurrent` rely on a common `isPrime(num int) bool` helper function. This function determines if a given number is prime. It includes an optimization to quickly discard even numbers (except 2) and only checks for divisibility by odd numbers up to the square root of the given number.

## Running the Code

To execute the prime summation calculation (using `runner.go`):

1.  Navigate to the `Intermediate/Exercise04_Benchmarking` directory in your terminal.
2.  Run the executable:
    ```bash
    go run ./cmd/primesum/main.go
    ```
    (Or `go run ./cmd/primesum` if running from the `Exercise04_Benchmarking` directory).

## Running the Benchmarks

To measure the performance of both implementations:

1.  Ensure your terminal's working directory is `Intermediate/Exercise04_Benchmarking`.
2.  Execute the benchmark command:
    ```bash
    go test -bench=. -benchmem
    ```
    * `-bench=.`: Runs all benchmark functions in the current directory.
    * `-benchmem`: Reports memory allocation statistics (bytes per operation and allocations per operation).

## Benchmark Results Analysis

The benchmarks were executed on a `13th Gen Intel(R) Core(TM) i5-13600KF` CPU.

### Raw Output:

```text
cpu: 13th Gen Intel(R) Core(TM) i5-13600KF
BenchmarkSumPrimesSequential_Small-20          247124           4782 ns/op            0 B/op           0 allocs/op
BenchmarkSumPrimesSequential_Medium-20            396        2946750 ns/op            0 B/op           0 allocs/op
BenchmarkSumPrimesSequential_Large-20              16       70005062 ns/op            0 B/op           0 allocs/op
BenchmarkSumPrimesSequential_XLarge-20              3      489632833 ns/op            0 B/op           0 allocs/op
BenchmarkSumPrimesConcurrent_Small-20          143547           9064 ns/op         1312 B/op          43 allocs/op
BenchmarkSumPrimesConcurrent_Medium-20           2196         458910 ns/op         1374 B/op          43 allocs/op
BenchmarkSumPrimesConcurrent_Large-20             128        9288807 ns/op         1312 B/op          43 allocs/op
BenchmarkSumPrimesConcurrent_XLarge-20             20       55841275 ns/op         1312 B/op          43 allocs/op
PASS
ok      github.com/Daniel-Q-Reis/GoroutinesFromBeginningToAdvanced/Intermediate/Exercise04_Benchmarking 12.812s
```

### Key Observations:

1.  **Performance (ns/op - Nanoseconds per Operation):**
    * **Small N (1,000):** The `Sequential` implementation (`4782 ns/op`) significantly outperforms the `Concurrent` one (`9064 ns/op`). This indicates that for very small workloads, the overhead of goroutine creation, scheduling, and mutex synchronization outweighs any benefits of parallelism.
    * **Medium N (100,000):** The `Concurrent` implementation (`458910 ns/op`) becomes substantially faster than the `Sequential` one (`2946750 ns/op`). This represents a performance gain of approximately **6.4x**.
    * **Large N (1,000,000):** The `Concurrent` implementation (`9288807 ns/op`) continues to show a strong lead over `Sequential` (`70005062 ns/op`), with a gain of approximately **7.5x**.
    * **X-Large N (4,000,000):** The performance advantage of the `Concurrent` version is even more pronounced (`55841275 ns/op` vs. `489632833 ns/op`), yielding a gain of approximately **8.7x**.

2.  **Memory Consumption (B/op - Bytes per Operation & allocs/op - Allocations per Operation):**
    * **Sequential:** `0 B/op` and `0 allocs/op`. This is highly efficient, indicating minimal dynamic memory allocation and constant memory usage per operation.
    * **Concurrent:** Approximately `1300 B/op` and `43 allocs/op` across all larger `N` values. This constant overhead is attributed to the creation of goroutines, `sync.WaitGroup`, and `sync.Mutex` instances. Crucially, this overhead does not scale with the input size `N`, making the concurrent approach efficient for larger workloads despite this baseline cost.

### Conclusion:

The benchmark results clearly demonstrate that while concurrency introduces a fixed overhead, it provides significant performance benefits for CPU-bound tasks like prime number summation when the workload `N` is sufficiently large. The concurrent solution effectively utilizes multiple CPU cores to distribute the computation, leading to substantial speedups (up to ~8.7x) compared to the sequential approach for larger inputs. For very small inputs, the simplicity and lack of synchronization overhead make the sequential approach preferable.

---