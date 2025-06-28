package main

import (
	"container/heap" // Provides heap operations for priority queue
	"fmt"            // For formatted output
	"sync"           // For synchronization primitives
	"time"           // For pauses and simulating processing
)

// Task represents a task to be processed with the following fields:
// - ID: unique identifier for tracking
// - Priority: numerical value where a lower number means higher priority (1 is more important than 5)
// - Payload: the actual content/data of the task
// - Index: internal field required by the heap for its operation (maintains current position in the heap)
type Task struct {
	ID       int    // Must be unique for each task
	Priority int    // Priority: 1 is the highest
	Payload  string // Actual task data
	Index    int    // Managed by the heap, indicates current position in the queue
}

// PriorityQueue is a slice of Task pointers that implements heap.Interface.
// The heap maintains tasks ordered by priority (and ID for tie-breaking).
type PriorityQueue []*Task

// Len returns the number of items in the heap (required by heap.Interface)
func (pq PriorityQueue) Len() int {
	return len(pq) // Returns the current size of the queue
}

// Less compares two items to determine their order in the heap (required by heap.Interface)
// Returns true if the item at index i should come before the item at index j.
// The comparison logic considers:
// 1. Higher priority (lower numerical value) first.
// 2. In case of a tie, uses lower ID as a tie-breaker (FIFO for same priority).
func (pq PriorityQueue) Less(i, j int) bool {
	// Primary comparison by priority
	if pq[i].Priority != pq[j].Priority {
		return pq[i].Priority < pq[j].Priority // Lower Priority value first
	}
	return pq[i].ID < pq[j].ID // Tie-breaking by lower ID (FIFO for same priority)
}

// Swap exchanges the items at indices i and j (required by heap.Interface)
// In addition to swapping positions in the slice, it updates the Index fields of the items
// to reflect their new positions in the heap.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i] // Swap the positions of items in the slice
	pq[i].Index = i             // Update internal index for item at new position i
	pq[j].Index = j             // Update internal index for item at new position j
}

// Push adds an item to the heap (required by heap.Interface)
// The parameter x is of type interface{}, so we need to perform type assertion to *Task.
// It sets the item's Index to its position at the end of the slice before heap reordering.
func (pq *PriorityQueue) Push(x interface{}) {
	task := x.(*Task)       // Convert the interface to *Task
	task.Index = len(*pq)   // Set the index to the last position
	*pq = append(*pq, task) // Add to the end of the slice
}

// Pop removes and returns the highest-priority item from the heap (required by heap.Interface)
// It always removes the last element of the slice (after heap's internal reordering).
// The Index of the removed item is set to -1 to indicate it's no longer in the heap.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq // Store the current slice
	n := len(old)
	task := old[n-1]   // Get the last element (which is the root after heap operations)
	task.Index = -1    // Mark as removed (index -1)
	*pq = old[0 : n-1] // Reduce the slice size, the full slice size is old[0 : n] -> n := len(old)
	return task
}

// TaskProcessor manages the priority queue and coordinates workers.
type TaskProcessor struct {
	queue         *PriorityQueue // Priority queue based on heap
	mu            sync.Mutex     // Mutex to protect concurrent access to the queue
	cond          *sync.Cond     // Condition variable to notify workers about new tasks
	wg            sync.WaitGroup // WaitGroup to wait for all workers to finish
	done          chan struct{}  // Channel to signal shutdown to workers
	tasksInFlight int            // Counter for tasks currently being processed or waiting in queue
}

// NewTaskProcessor creates and initializes a new TaskProcessor.
// It initializes the priority queue, mutex, condition variable, and done channel.
// It also initializes the internal heap using heap.Init.
func NewTaskProcessor() *TaskProcessor {
	tp := &TaskProcessor{
		queue:         &PriorityQueue{},
		done:          make(chan struct{}),
		tasksInFlight: 0,
	}
	tp.cond = sync.NewCond(&tp.mu) // Creates the condition variable bound to the mutex
	heap.Init(tp.queue)            // Initializes the heap structure
	return tp
}

// AddTask adds a new task to the priority queue in a thread-safe manner.
// 1. Locks the mutex for exclusive access to the queue.
// 2. Uses heap.Push to add the task while maintaining priority order.
// 3. Increments the tasksInFlight counter.
// 4. Notifies one worker through the condition variable that a new task is available.
// 5. Releases the mutex when finished (via defer).
func (tp *TaskProcessor) AddTask(task *Task) {
	tp.mu.Lock()         // Ensures exclusive access
	defer tp.mu.Unlock() // Releases the lock at the end

	heap.Push(tp.queue, task) // Adds the task maintaining priority order
	tp.tasksInFlight++        // Increment counter
	tp.cond.Signal()          // Notifies one waiting worker
}

// RunWorkers starts the specified number of worker goroutines.
// Each worker will be responsible for processing tasks from the queue.
// It increments the WaitGroup for each worker launched.
func (tp *TaskProcessor) RunWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		tp.wg.Add(1)        // Increments the WaitGroup counter
		go tp.worker(i + 1) // Starts each worker with a unique ID (1-based)
	}
}

// worker is the function executed by each worker goroutine.
// 1. When finished, it notifies the WaitGroup that it has completed (via defer).
// 2. Enters an infinite loop until it receives a shutdown signal.
// 3. Locks the mutex to safely check the queue.
// 4. If the queue is empty, it waits for a new task or a shutdown signal.
// 5. When tasks are available, it gets the highest priority task using heap.Pop.
// 6. Releases the mutex before processing the task.
// 7. Processes the task (simulated with sleep).
// 8. Decrements the tasksInFlight counter and signals if all tasks are processed.
func (tp *TaskProcessor) worker(workerID int) {
	defer tp.wg.Done() // Ensures WaitGroup is notified upon completion

	for { // Infinite loop until shutdown signal
		tp.mu.Lock() // Locks access to the queue

		// This loop waits for tasks to become available or for a shutdown signal.
		// It continuously checks if the queue is empty.
		for len(*tp.queue) == 0 {
			select {
			case <-tp.done: // Check if shutdown signal received while waiting for tasks
				tp.mu.Unlock() // Release lock before exiting
				fmt.Printf("Worker %d: Shutdown signal received. Exiting.\n", workerID)
				return // Exit goroutine gracefully
			default:
				// If no shutdown signal yet, wait for a new task.
				// cond.Wait() releases the lock and blocks, then re-acquires it upon wake-up.
				// It must be in a loop because it can return spuriously (without a direct signal).
				tp.cond.Wait()
			}
		}

		// After exiting the inner loop, the queue is NOT empty.
		// However, it's good practice to re-check for a shutdown signal
		// in case it arrived just as we were waking up from cond.Wait().
		select {
		case <-tp.done:
			tp.mu.Unlock()
			fmt.Printf("Worker %d: Shutdown signal received after finding tasks. Exiting.\n", workerID)
			return
		default:
			// No shutdown signal, proceed to pop and process the task.
		}

		// Remove the highest priority task from the queue
		task := heap.Pop(tp.queue).(*Task)
		tp.mu.Unlock() // Release the lock before processing (processing can take time)

		// Simulate task processing
		fmt.Printf("Worker %d processing Task ID %d (Priority: %d, Payload: %s)\n",
			workerID, task.ID, task.Priority, task.Payload)
		time.Sleep(500 * time.Millisecond) // Simulate processing time

		tp.mu.Lock() // Re-acquire lock to safely decrement tasksInFlight
		tp.tasksInFlight--
		if tp.tasksInFlight == 0 {
			// If all tasks have been processed (this worker processed the last one),
			// signal the main goroutine which might be waiting for this condition.
			// Using Broadcast() here to ensure all waiting goroutines (including main) are woken up.
			tp.cond.Broadcast()
		}
		tp.mu.Unlock()
	}
}

// Shutdown gracefully terminates the task processor.
// 1. Closes the done channel to signal shutdown.
// 2. Uses Broadcast to wake up all waiting workers.
// 3. Waits for all workers to finish using WaitGroup.
// 4. Prints a confirmation message.
func (tp *TaskProcessor) Shutdown() {
	close(tp.done)      // Sends shutdown signal
	tp.cond.Broadcast() // Wakes up all workers
	tp.wg.Wait()        // Waits for all workers to finish
	fmt.Println("Task processor shutdown completed")
}

//============================================ MAIN FUNC =================================

func main() {
	// Create a new instance of the TaskProcessor
	tp := NewTaskProcessor()

	// Define the number of workers (goroutines that will process tasks)
	numWorkers := 3

	// Start the workers
	tp.RunWorkers(numWorkers)

	// List of tasks to be processed, with different priorities
	tasksToAdd := []*Task{
		{ID: 1, Priority: 5, Payload: "Process order 100"},
		{ID: 2, Priority: 1, Payload: "Generate critical report"},
		{ID: 3, Priority: 10, Payload: "Routine backup"},
		{ID: 4, Priority: 2, Payload: "Update dashboard"},
		{ID: 5, Priority: 1, Payload: "Resolve incident"},
		{ID: 6, Priority: 3, Payload: "Send email notification"},
		{ID: 7, Priority: 5, Payload: "Update user profile"},
	}

	// WaitGroup to ensure all tasks are added by the producer goroutine
	var taskProducerWg sync.WaitGroup
	taskProducerWg.Add(1)

	// Goroutine to add tasks to the queue with intervals.
	// This simulates asynchronous task arrival.
	go func() {
		defer taskProducerWg.Done() // Signal that all tasks have been added when this goroutine exits.
		for i, task := range tasksToAdd {
			// Simulate increasing delay between tasks for asynchronous arrival
			time.Sleep(time.Duration(200*(i+1)) * time.Millisecond)
			tp.AddTask(task) // Add task to the processor (thread-safe method)
			fmt.Printf("Added Task ID %d (Priority: %d)\n", task.ID, task.Priority)
		}
		fmt.Println("All tasks have been submitted to the processor.")
	}()

	// Wait for the task-adding goroutine to finish submitting all tasks.
	taskProducerWg.Wait()

	// Now, wait in the main goroutine until all tasks (those submitted) have been processed.
	// We need to acquire the mutex to safely check the `tasksInFlight` counter.
	tp.mu.Lock()
	// Loop while there are still tasks being processed or waiting in the queue.
	// `tp.cond.Wait()` releases the mutex, blocks, and re-acquires the mutex upon being woken up.
	// This loop is crucial because `cond.Wait()` can return spuriously.
	for tp.tasksInFlight > 0 {
		fmt.Println("Main: Waiting for tasks to complete...") // Debugging message
		tp.cond.Wait()
	}
	tp.mu.Unlock()
	fmt.Println("All submitted tasks have been processed by workers.")

	// Initiate shutdown of the processor (workers will exit cleanly).
	tp.Shutdown()

	// Final message
	fmt.Println("Program execution completed.")
}
