# Exercise 03: Simple Pipeline Processing

## Objective

This exercise demonstrates how to build a basic processing pipeline using Go's goroutines and channels. Data flows sequentially through multiple stages, with each stage performed by a dedicated goroutine.

## Problem

Implement a three-stage pipeline:
1.  **Generator:** Produces a sequence of numbers.
2.  **Duplicator:** Receives numbers from the generator, doubles them.
3.  **Adder:** Receives doubled numbers, adds five to them.
The main function then consumes and prints the final processed results.

## Solution

The solution leverages Go channels to connect each stage of the pipeline.
-   The `generateNumbers` goroutine creates initial integer values and sends them to an output channel. Crucially, it closes this channel once all numbers are sent, signaling the next stage that no more data is coming.
-   The `duplicateNumbers` goroutine reads from its input channel (which is the output of the generator), processes each number by doubling it, and then sends the result to its own output channel. It also closes its output channel when its input channel is closed and all data has been processed.
-   Similarly, the `addFive` goroutine reads from its input channel (the output of the duplicator), adds five to each number, and sends the final result to its output channel, closing it when done.
-   The `main` function acts as the final consumer, reading all results from the last stage's output channel using a `for-range` loop. This loop naturally terminates when the channel is closed by the `addFive` goroutine.

This pattern effectively uses channel closures as a signal for graceful shutdown across the pipeline stages, eliminating the need for explicit `sync.WaitGroup` in this particular sequential flow.

## How to Run

1. Navigate to this directory:
   ```bash
   cd GoroutinesFromBeginningToAdvanced/Beginner/Exercise03_SimplePipeline
   ```

2. Run the Go program:
   ```bash
   go run main.go
   ```

## Expected Output

The output will show the flow of numbers through each stage of the pipeline, from generation to final processing, demonstrating the sequential nature of data transformation across concurrent goroutines.

Generator sent: 1
Generator sent: 2
Duplicator received: 1, sending: 2
Generator sent: 3
Adder received: 2, sending: 7
Duplicator received: 2, sending: 4
Generator sent: 4
Adder received: 4, sending: 9
Duplicator received: 3, sending: 6
... (output will vary due to concurrency, but the flow is sequential)
Final Result: 7
Final Result: 9
...
All pipeline stages are completed