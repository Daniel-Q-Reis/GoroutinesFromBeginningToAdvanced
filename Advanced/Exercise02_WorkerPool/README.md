# Advanced/Exercise02 - Distributed Worker Pool with Dynamic Task Distribution

## Objective

This exercise demonstrates the implementation of a sophisticated concurrent system in Go: a distributed worker pool. It showcases how to dynamically generate and distribute tasks among a pool of workers, handle varying task complexities, and ensure a robust, graceful shutdown.

## Problem Statement

The challenge is to simulate a system where a "Producer" continuously generates "Data Processing Tasks." These tasks, each with varying simulated "Complexity" (CPU-intensive calculation + `time.Sleep`), need to be efficiently processed by a pool of "Workers." The system must manage task distribution, monitor completion, and allow for a clean termination.

## Core Components and Solution Overview

The solution is structured into several interacting components, each implemented as a Go package or function, demonstrating advanced concurrency patterns:

1.  **`Task` (in `task.go`):**
    * Defines the unit of work, including `ID`, `Data` (for calculation), `Complexity` (simulated work duration), `Result`, and `Err`.
    * Includes an `isPrime` helper function for CPU-intensive calculations.

2.  **`Producer` (in `producer.go`):**
    * Generates a fixed number of `Task` instances with random `Data` and `Complexity`.
    * Sends these tasks to an **unbuffered channel (`TaskChan`)**. This unbuffered nature is crucial for applying **backpressure**: the producer will block if no worker is ready to receive a task, preventing the producer from overwhelming the system.
    * Closes the `TaskChan` after all tasks are generated, signaling completion.

3.  **`Worker` (in `worker.go`):**
    * Represents an individual worker in the pool.
    * Continuously reads `Task`s from the `TaskChan`.
    * Processes each task by performing a CPU-intensive calculation (e.g., `isPrime`) and simulating work duration (`time.Sleep(task.Complexity)`).
    * Sends the processed `Task` (with `Result` and `Err` populated) to a **buffered channel (`ResultChannel`)**. The buffer allows workers to send results without immediately blocking, improving throughput.

4.  **`Pool` (in `pool.go`):**
    * Manages the lifecycle of the worker pool.
    * Initializes the `TaskChan` (unbuffered) and `ResultChan` (buffered).
    * Launches the specified number of `Worker` goroutines.
    * Includes a `sync.WaitGroup` to track the completion of all workers and ensures that the `ResultChannel` is closed only after all workers have finished processing their tasks.

5.  **`Consumer` (in `consumer.go`):**
    * Reads processed `Task`s from the `ResultChannel`.
    * Logs the task details (`ID`, `Data`, `Result`).
    * Keeps track of the total number of tasks processed.

6.  **`main` (in `cmd/workerpool/main.go`):**
    * Orchestrates the entire system.
    * Initializes the `Pool`, `Producer`, and `Consumer`.
    * Launches the `Producer` and `Consumer` goroutines.
    * Uses a `sync.WaitGroup` to wait for the `Producer` to finish sending tasks and the `Consumer` to finish processing all results, ensuring a graceful system shutdown.
    * Reports a summary of the execution, including total tasks processed, number of workers, and total execution time.

This architecture demonstrates effective use of Go's concurrency primitives to build a scalable and resilient task processing system.

## Project Structure

```bash
Advanced/Exercise02_WorkerPool/
├── cmd/
│   └── workerpool/
│       └── main.go       # Main executable (package main)
├── go.mod                # Go module file for this package
├── task.go               # Task struct definition and isPrime helper (package exercise02workerpool)
├── producer.go           # Producer logic (package exercise02workerpool)
├── worker.go             # Worker logic (package exercise02workerpool)
├── pool.go               # Pool management logic (package exercise02workerpool)
└── consumer.go           # Consumer logic (package exercise02workerpool)
├── README.md             # This file
```

## Acknowledgements

This solution was developed with the aid of AI tools for initial brainstorming and code scaffolding. However, the complex concurrency debugging and resolution of deadlocks were entirely performed through manual analysis and iteration.

## How to Run

1.  **Initialize Go Module (if not already done):**
    Navigate to the root of your repository (`GoroutinesFromBeginningToAdvanced`) and run:
    ```bash
    go mod init [github.com/Daniel-Q-Reis/GoroutinesFromBeginningToAdvanced](https://github.com/Daniel-Q-Reis/GoroutinesFromBeginningToAdvanced)
    ```
    (Replace `Daniel-Q-Reis` with your actual GitHub username if different).

2.  **Navigate to the Exercise Directory:**
    ```bash
    cd Advanced/Exercise02_WorkerPool
    ```

3.  **Run the application:**
    ```bash
    go run ./cmd/workerpool
    ```

## Expected Output

The output will show tasks being processed by different workers, with varying completion times due to the simulated complexity. The order of "Task X: Data=Y, isPrime=Z" messages will be non-deterministic due to concurrency. Finally, a summary will be displayed.

```text
Task 0: Data=12345, isPrime=false
Task 1: Data=67890, isPrime=true
... (many more lines of tasks being processed)
Processed 1000 tasks
System Summary:
Total tasks processed: 1000
Number of workers: X (20, based on my CPU intel core i5 13600k)
Total execution time: Y (3,82s also based on my System)