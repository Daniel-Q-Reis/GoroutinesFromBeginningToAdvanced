# Exercise 02: Simple Producer-Consumer with Channels

## Objective

This exercise demonstrates the fundamental producer-consumer pattern using Go's concurrency primitives: goroutines and channels. It highlights how goroutines can communicate and synchronize by sending and receiving data through a channel.

## Problem

Implement a system with multiple producers sending integer numbers to a shared channel, and multiple consumers receiving and processing these numbers from the same channel. Ensure proper synchronization and graceful shutdown.

## Solution

The solution involves:
-   `producer` goroutines: Generate numbers and send them to the channel. Each producer signals its completion via a `sync.WaitGroup` dedicated to producers (`wgProducers`).
-   `consumer` goroutines: Read numbers from the channel using a `for-range` loop. Each consumer signals its completion via a separate `sync.WaitGroup` for consumers (`wgConsumers`).
-   A dedicated orchestrator goroutine in `main`: This goroutine waits for all `wgProducers` to signal completion. Once all producers are done, it safely closes the shared channel. This is crucial as a channel must only be closed once, and only after all senders have finished.
-   The `main` function then waits for all `wgConsumers` to finish their work, ensuring all numbers are processed before the program exits.

This setup effectively manages data flow and lifecycle in a concurrent environment, preventing common issues like panics from closing an already closed channel or deadlocks from improper synchronization.

## How to Run

1. Navigate to this directory:
   ```bash
   cd GoroutinesFromBeginningToAdvanced/Beginner/Exercise02_ProducerConsumer
   ```

2. Run the Go program:
   ```bash
   go run main.go
   ```

## Expected Output
The output will show producers sending numbers and consumers receiving them, with intermingled messages due to concurrent execution. The order of processing by consumers might vary. Eventually, all producers will finish, the channel will be closed, and consumers will drain the remaining values before exiting.

```text
Producer 1 sending: 1
Producer 2 sending: 1
Consumer 1 received: 1
Producer 1 sending: 2
Consumer 2 received: 1
Producer 2 sending: 2
... (output will vary)
Channel closed by main orchestrator.
Consumer 1 finished.
Consumer 2 finished.
All jobs are done
```