# Exercise 02: Simple Worker Pool

## Objective

This exercise demonstrates the implementation of a basic worker pool pattern in Go. A worker pool is a common concurrency pattern used to limit the number of concurrently executing tasks, improving resource management and performance.

## Problem

Create a system where a fixed number of "worker" goroutines process a larger number of "tasks" distributed via a channel. Each task involves a simple calculation (e.g., squaring a number). The main program should generate tasks, distribute them to the pool, and collect the processed results.

## Solution

The solution establishes a producer-consumer model with multiple consumers (workers):
-   A `Task` struct defines the work unit, containing an ID and a number to process.
-   The `main` goroutine launches a dedicated goroutine to act as a **task generator**. This generator creates tasks with random numbers, sends them to a `tasks` channel, and then closes the channel to signal that no more tasks will be sent.
-   A fixed number of `worker` goroutines are started. Each worker reads tasks from the `tasks` channel, performs a simulated computation (squaring the number and sleeping), and sends the result to a `results` channel. Each worker uses `defer wg.Done()` to signal its completion.
-   A separate goroutine in `main` is responsible for **closing the `results` channel**. It waits for all workers to signal completion via a `sync.WaitGroup` (`workerWg`) before safely closing `results`. This prevents panics if `results` is closed prematurely while workers are still trying to send to it.
-   Finally, the `main` goroutine reads all processed results from the `results` channel using a `for-range` loop. It then calculates the square root of the result to display the original number, providing clearer context about the processed data. This loop gracefully terminates when `results` is closed.

This pattern efficiently distributes work among a limited number of concurrent processes, preventing the system from being overwhelmed by too many simultaneous operations.

## How to Run

1. Navigate to this directory:
   ```bash
   cd GoroutinesFromBeginningToAdvanced/Intermediate/Exercise02_WorkerPool
   ```

2. Run the Go program:
   ```bash
   go run main.go
   ```

## Expected Output

The output will show workers processing tasks concurrently, with tasks having random numbers. The results in the main function will display both the squared result and the original number (derived via square root), providing clearer context to the processed data. Due to the nature of concurrent execution, the order of task processing and result reception might vary.

```text
Collecting results...
Worker 1 processing Task ID 1 (Number: 42)
Worker 2 processing Task ID 2 (Number: 15)
Worker 3 processing Task ID 3 (Number: 97)
Worker 4 processing Task ID 4 (Number: 74)
Worker 5 processing Task ID 5 (Number: 86)
Result: 1764 after processing the Number: 42
Worker 1 processing Task ID 6 (Number: 47)
Result: 225 after processing the Number: 15
Worker 2 processing Task ID 7 (Number: 87)
Result: 9409 after processing the Number: 97
Worker 3 processing Task ID 8 (Number: 8)
Result: 5476 after processing the Number: 74
Worker 4 processing Task ID 9 (Number: 85)
Result: 7396 after processing the Number: 86
Worker 5 processing Task ID 10 (Number: 92)
... (output will vary significantly due to random numbers and concurrency)
Results channel closed.
All tasks are done.
```