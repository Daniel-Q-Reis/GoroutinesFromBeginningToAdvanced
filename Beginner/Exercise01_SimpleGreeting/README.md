# Exercise 01: Simple Concurrent Greeting

## Objective

This exercise demonstrates the basic usage of Goroutines to perform concurrent tasks and `sync.WaitGroup` to wait for all goroutines to complete.

## Problem

Create a Go program that prints greetings from a list of names concurrently. Each greeting should be handled by a separate goroutine.

## Solution

The `greet` function prints a personalized message and simulates some work with a short sleep. In `main`, a list of names is iterated, and for each name, a new goroutine is launched to call `greet`. `sync.WaitGroup` is used to ensure the main program waits for all goroutines to finish their greetings before exiting.

A key learning point in this exercise is understanding how to correctly pass loop variables to goroutines to avoid unexpected behavior (e.g., all goroutines printing the last value of the loop variable).

## How to Run

1. Navigate to this directory:
   ```bash
   cd GoroutinesFromBeginningToAdvanced/Beginner/Exercise01_SimpleGreeting
   ```
2. Run the Go program:
   ```bash
   go run main.go
   ```

## Expected Output

The output will show greetings from each name, potentially in a non-deterministic order due to concurrency, followed by a completion message. For example:

```text
Hello, my name is Alice!
Hello, my name is Peter! 
Hello, my name is James! 
Hello, my name is Jordan! 
Hello, my name is Rob! 
All greetings are done!

(The order of "Hello, my name is..." lines may vary.)
```