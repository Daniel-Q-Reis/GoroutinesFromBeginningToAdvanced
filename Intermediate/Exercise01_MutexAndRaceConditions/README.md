 # Exercise 01: Mutex and Race Conditions

## Objective

This exercise demonstrates a common concurrency problem known as a "race condition" and how to solve it effectively using Go's `sync.Mutex`. It highlights the importance of protecting shared resources from simultaneous access by multiple goroutines.

## Problem

Multiple goroutines concurrently increment a shared integer variable (`counter`). Without proper synchronization, the final value of the counter will be less than the expected sum of all increments, due to read-modify-write operations not being atomic.

## Solution

The solution showcases two phases:
1.  **Demonstrating Race Condition (Implicit):** Initially, the `counter` is incremented by multiple goroutines without any locking mechanism. Running the program multiple times will likely result in different and incorrect final `counter` values, clearly indicating a race condition. (To explicitly demonstrate this phase, you would run the code *before* adding the mutex, then run *after*.)
2.  **Resolving with `sync.Mutex`:** A `sync.Mutex` (`mu`) is introduced. Before each increment operation (`counter++`), `mu.Lock()` is called to acquire a lock, ensuring that only one goroutine can modify `counter` at a time. After the increment, `mu.Unlock()` releases the lock. This guarantees atomicity of the increment operation, leading to the correct final `counter` value consistently.

`sync.WaitGroup` is used to ensure the `main` goroutine waits for all worker goroutines to complete their incrementing tasks before printing the final result.

## How to Run

1. Navigate to this directory:
   ```bash
   cd GoroutinesFromBeginningToAdvanced/Intermediate/Exercise01_MutexAndRaceConditions
   ```

2. Run the Go program:
   ```bash
   go run main.go
   ```

## Expected Output

The Worker finished incrementing message will appear multiple times (once per worker). The crucial part is the Final Counter Value.

Without Mutex (Race Condition):
The final counter value would be inconsistent and less than the expected 100000 (100 workers * 1000 increments). For example:

```text
Worker finished incrementing
... (many more lines)
Final Counter Value (expected 100000): 98765
```
(Note: the exact incorrect value would vary on each run)

With Mutex (Correct):
The final counter value will consistently be the expected 100000.

```text
Worker finished incrementing
... (many more lines, as many as numWorkers due the loop)
Final Counter Value (expected 100000): 100000
```