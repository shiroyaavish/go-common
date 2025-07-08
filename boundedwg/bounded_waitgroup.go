/*
Package boundedwg provides a BoundedWaitGroup which limits the number of concurrently running goroutines.

It is designed to be a drop-in replacement for sync.WaitGroup with the added functionality of bounding concurrency.
BoundedWaitGroup works similarly to sync.WaitGroup, but it blocks when adding a new goroutine
if the concurrency limit has been reached, ensuring that no more than the specified number of goroutines
are running at the same time.

# Usage

To create a new BoundedWaitGroup with a concurrency limit of 5:

	bwg := boundedwg.New(5)

# Example

	package main
	import (
	    "fmt"
	    "time"

	    "*go-common/boundedwg"
	)

	func main() {
	    bwg := boundedwg.New(5)

	    for i := 0; i < 10; i++ {
	        bwg.Add(1)
	        go func(i int) {
	            defer bwg.Done()
	            time.Sleep(time.Second)
	            fmt.Printf("Goroutine %d finished\n", i)
	        }(i)
	    }

	    bwg.Wait()
	    fmt.Println("All goroutines finished")
	}

# Use Cases

  - Limiting the number of concurrently running goroutines to prevent resource exhaustion.

  - Controlling the concurrency level when dealing with external systems (e.g., API rate limits, database connections).

  - Replacing sync.WaitGroup in situations where bounded concurrency is desired.

# Considerations

  - **Initialization**: BoundedWaitGroup must be initialized using the New function. Do not attempt to instantiate it directly. The zero value of BoundedWaitGroup is nil and not usable. This design prevents accidental misuse and ensures that the semaphore channel is properly initialized.

  - **Blocking Behavior**: The Add method will block if adding goroutines would exceed the concurrency limit. Ensure that this behavior is acceptable in your application.

  - **Concurrency Limit**: Choose an appropriate concurrency limit based on your application's needs and the resources available. Setting a limit that is too high may lead to increased resource usage and contention. Conversely setting a limit too low may impact performance of the code.

# Best Practices

  - **Avoiding Direct Instantiation**: Always use the New function to create a BoundedWaitGroup. This ensures proper initialization and prevents nil pointer dereferences.

  - **Handling Blocking in Add**: Be aware that the Add method can block if the concurrency limit is reached. Design your application accordingly to avoid deadlocks or performance bottlenecks.

  - **Resource Management**: Consider the resources available on the system when setting the concurrency limit. In environments with limited resources, a lower concurrency limit may be appropriate.
*/
package boundedwg

import (
	"sync"
)

// boundedWaitGroup is unexported to prevent direct instantiation.
type boundedWaitGroup struct {
	wg        sync.WaitGroup
	semaphore chan struct{}
}

// BoundedWaitGroup limits the number of concurrently running goroutines.
//
// It panics if the limit <= 0
//
// It provides an interface similar to sync.WaitGroup but with bounded concurrency.
// BoundedWaitGroup must be created using the New function.
type BoundedWaitGroup = *boundedWaitGroup

// New creates a new BoundedWaitGroup with the specified concurrency limit.
//
// Example:
//
//	bwg := boundedwg.New(5) // Limit concurrency to 5
func New(limit int) BoundedWaitGroup {
	if limit <= 0 {
		panic("BoundedWaitGroup go routine limit must be greater than 0")
	}
	return &boundedWaitGroup{
		semaphore: make(chan struct{}, limit),
	}
}

// Add increments the WaitGroup counter by delta.
//
// Add blocks if adding goroutines would exceed the concurrency limit.
// It panics if delta is negative.
//
// Example:
//
//	bwg.Add(1)
//	go func() {
//	    defer bwg.Done()
//	    // Goroutine work
//	}()
func (bwg BoundedWaitGroup) Add(delta int) {
	for i := 0; i < delta; i++ {
		bwg.semaphore <- struct{}{} // Will block if limit is reached
	}
	bwg.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one.
//
// It should be called when a goroutine has completed its work.
//
// Example:
//
//	defer bwg.Done()
func (bwg BoundedWaitGroup) Done() {
	<-bwg.semaphore // Release a slot
	bwg.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
//
// Example:
//
//	bwg.Wait()
func (bwg BoundedWaitGroup) Wait() {
	bwg.wg.Wait()
}
