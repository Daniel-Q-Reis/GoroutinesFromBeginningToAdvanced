# Exercise 03: Goroutine Cancellation with Context

## Objective

This exercise demonstrates the critical use of `context.Context` for managing the lifecycle of goroutines, specifically for performing graceful cancellation. This is essential for preventing goroutine leaks and ensuring proper resource cleanup in long-running concurrent applications.

## Problem

In concurrent Go applications, goroutines often run indefinitely or for long periods. It's crucial to have a mechanism to signal these goroutines to stop their work and shut down gracefully when they are no longer needed (e.g., due to a user request, a timeout, or a system shutdown). Without such a mechanism, goroutines might continue consuming resources unnecessarily.

## Solution

The solution utilizes `context.WithCancel` to create a `Context` that can be cancelled.
-   The `main` function initializes a parent `Context` and a `cancel` function using `context.WithCancel`.
-   Multiple `cancellableWorker` goroutines are launched. Each worker receives this `Context` as an argument.
-   Inside each `cancellableWorker`, a `select` statement continuously monitors two channels:
    -   `<-ctx.Done()`: This case is triggered when the `cancel()` function is called. Upon receiving this signal, the worker prints a message and `returns`, effectively shutting down the goroutine gracefully.
    -   `<-taskChannel`: This case receives simulated tasks. The worker processes the task and performs a simulated sleep.
    -   A `default` case with a short `time.Sleep` is included to prevent the `select` from blocking indefinitely when neither a task nor a cancellation signal is immediately available. This allows the worker to repeatedly check the `ctx.Done()` channel.
-   A separate goroutine in `main` is responsible for simulating the event that triggers cancellation. After a predefined duration, it calls `cancel()`, sending the cancellation signal to all associated workers.
-   `sync.WaitGroup` ensures that the `main` goroutine waits for all workers to acknowledge the cancellation and shut down before the program exits, confirming proper cleanup.

This pattern is a cornerstone of robust Go concurrent programming, enabling controlled termination of goroutines.

## How to Run

1. Navigate to this directory:
   ```bash
   cd GoroutinesFromBeginningToAdvanced/Intermediate/Exercise03_GoroutineCancellation
   ```

2. Run the Go program:
   ```bash
   go run main.go
   ```

## Expected Output

The output will show workers processing tasks for a period. After approximately 1.1 seconds, the cancellation signal will be sent. You will then observe the workers acknowledging the cancellation and gracefully shutting down.

```text
Worker 1: Processing task 0
Worker 2: Processing task 1
Worker 3: Processing task 2
Worker 1: Processing task 3
Worker 2: Processing task 4
Worker 3: Processing task 5
Worker 3: Processing task 6
Main: Context cancelled. Signaling workers to stop.
Worker 1: Context cancelled. Shutting down. 
Worker 2: Context cancelled. Shutting down. 
Worker 3: Context cancelled. Shutting down. 
The program successfully stopped workers using ''cancel() context.Context''
```
