# Exercise 01: Prioritized Task Processor with `container/heap`

## Objective

This exercise demonstrates how to implement a concurrent task processor that handles tasks with varying priorities. It ensures that higher-priority tasks are processed before lower-priority ones, even if they arrive in a different order. This setup utilizes Go's `sync` primitives (`sync.Mutex`, `sync.Cond`, `sync.WaitGroup`) and the `container/heap` package for efficient priority queue management.

## Concepts Explored

-   **`container/heap`**: Understanding and implementing `heap.Interface` to create a custom min-heap based on task priority.
-   **Interfaces**: Applying Go interfaces to enable `container/heap`'s generic operations on a custom data type.
-   **`sync.Mutex`**: Protecting shared resources (the priority queue) from concurrent read/write access.
-   **`sync.Cond`**: Efficiently signaling worker goroutines when new tasks are available, preventing busy-waiting and conserving CPU cycles.
-   **Goroutine Management**: Launching and gracefully shutting down worker goroutines.
-   **Concurrency Patterns**: Implementing a sophisticated producer-consumer model where consumers (workers) prioritize tasks.

## Problem

Design a system where tasks, each with a unique ID, a priority level (lower number = higher priority), and a payload, are submitted asynchronously. Multiple worker goroutines should consume these tasks from a central queue. The key requirement is that workers must always pick and process the task with the highest priority available in the queue. The system should also support graceful shutdown, ensuring all in-flight tasks are completed before termination.

## Solution

The solution is structured around a `TaskProcessor` which orchestrates a pool of worker goroutines and a `PriorityQueue` implemented using `container/heap`:

1.  **`Task` Struct**: Defines the task with `ID`, `Priority`, `Payload`, and an `Index` field crucial for `container/heap`'s internal management.
2.  **`PriorityQueue`**: A custom type that implements `heap.Interface` (`Len`, `Less`, `Swap`, `Push`, `Pop`). The `Less` method is critical here; it orders tasks first by `Priority` (lower value is higher priority) and then by `ID` for tie-breaking (ensuring FIFO for same-priority tasks).
3.  **`TaskProcessor`**:
    * Holds the `PriorityQueue`, a `sync.Mutex` to protect the queue, a `sync.Cond` to signal workers, a `sync.WaitGroup` to track worker completion, a `done` channel for shutdown signals, and a `tasksInFlight` counter to track active tasks.
    * **`AddTask(task *Task)`**: Adds a task to the queue, increments `tasksInFlight`, and signals a waiting worker via `cond.Signal()`.
    * **`RunWorkers(numWorkers int)`**: Launches the specified number of worker goroutines.
    * **`worker(workerID int)`**: Each worker goroutine continuously loops:
        * It acquires the mutex to check the queue.
        * If the queue is empty, it calls `tp.cond.Wait()`, which atomically unlocks the mutex and puts the worker to sleep until `cond.Signal()` or `cond.Broadcast()` is called. Upon waking, the mutex is re-acquired. This prevents busy-waiting.
        * A `select` case `<-tp.done` is used within the worker's loop to allow for graceful exit if a shutdown signal arrives while waiting or processing.
        * Once a task is available (queue not empty), it `heap.Pop`s the highest priority task, releases the mutex, processes the task (simulated with `time.Sleep`), and then decrements `tasksInFlight`.
    * **`Shutdown()`**: Closes the `done` channel to signal all workers, calls `tp.cond.Broadcast()` to wake up any sleeping workers, and then `tp.wg.Wait()` to ensure all workers have finished before the program truly exits.
4.  **`main` Function**:
    * Initializes the `TaskProcessor` and starts workers.
    * Launches a goroutine to add tasks with varying priorities and delays, simulating asynchronous arrival.
    * Crucially, the `main` goroutine waits for `tasksInFlight` to reach zero, ensuring all submitted tasks are processed before initiating the `Shutdown`. This is done using `tp.cond.Wait()` on `tp.mu` while `tp.tasksInFlight > 0`.

This implementation provides a robust and efficient way to handle prioritized work in a concurrent Go application, highlighting advanced concurrency patterns and data structures.

## How to Run

1. Navigate to this directory:
   ```bash
   cd GoroutinesFromBeginningToAdvanced/Advanced/Exercise01_PrioritizedTaskProcessor
   ```

2. Run the Go program:
   ```bash
   go run main.go
   ```

## Expected Output

The output will demonstrate tasks being added in a seemingly random order but processed strictly according to their priority (lower Priority value means processed first). You'll see workers picking up higher priority tasks even if they were added later.

```text
Added Task ID 1 (Priority: 5)
Added Task ID 2 (Priority: 1)
Worker 2 processing Task ID 2 (Priority: 1, Payload: Gerar relatório crítico)
Added Task ID 3 (Priority: 10)
Worker 1 processing Task ID 5 (Priority: 1, Payload: Resolver incidente)
Added Task ID 4 (Priority: 2)
Worker 3 processing Task ID 4 (Priority: 2, Payload: Atualizar dashboard)
Added Task ID 5 (Priority: 1)
Added Task ID 6 (Priority: 3)
Worker 2 processing Task ID 6 (Priority: 3, Payload: Send email notification)
Added Task ID 7 (Priority: 5)
Worker 1 processing Task ID 1 (Priority: 5, Payload: Processar pedido 100)
Worker 3 processing Task ID 7 (Priority: 5, Payload: Update user profile)
Worker 2 processing Task ID 3 (Priority: 10, Payload: Backup de rotina)
Task processor shutdown completed
All tasks processed. Exiting.
```