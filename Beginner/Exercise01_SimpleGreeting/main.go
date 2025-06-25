package main

import (
	"fmt"
	"sync"
	"time"
)

func greet(name string) {
	fmt.Printf("Hello, my name is %s!\n", name)
	time.Sleep(time.Millisecond * 100)
}

func main() {
	names := []string{"Alice", "Peter", "James", "Jordan", "Rob"}

	var wg sync.WaitGroup

	for _, name := range names {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()

			greet(n)
		}(name)
	}

	wg.Wait()

	fmt.Println("All greetings are done!")

}
