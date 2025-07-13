package exercise02workerpool

import (
	"math/rand"
	"time"
)

type Producer struct {
	TaskCount    int
	TaskChan     chan<- Task
	RandomNumber *rand.Rand
}

// NewProducer creats a new Producer instance
func NewProducer(taskCount int, taskChan chan<- Task) *Producer {
	return &Producer{
		TaskCount:    taskCount,
		TaskChan:     taskChan,
		RandomNumber: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Starts begins task generation
func (p *Producer) Start() {
	for i := 0; i < p.TaskCount; i++ {
		// Create new task with random values
		task := Task{
			ID:         i,
			Data:       p.RandomNumber.Intn(1000000) + 1,
			Complexity: time.Duration(p.RandomNumber.Intn(195)+5) * time.Millisecond,
		}

		p.TaskChan <- task
	}

	close(p.TaskChan)
}
