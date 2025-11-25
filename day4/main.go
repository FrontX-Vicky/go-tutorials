package main

import (
	"fmt"
	"sync"
	"time"
)

// TODO: Implement PrintNumbers function
func PrintNumbers(n int) {
	for i := 1; i <= n; i++ {
		fmt.Println(i)
	}
}

// TODO: Implement Sum function that sends result to a channel
func Sum(numbers []int, resultChan chan int) {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	resultChan <- sum
}

// TODO: Demonstrate unbuffered channel (blocks)
func UnbufferedChannelDemo(resultChan chan<- int) {
	go func() {
		// simulate some work
		time.Sleep(2 * time.Second)
		resultChan <- 42
	}()
}

// TODO: Demonstrate buffered channel (doesn't block until full)
func BufferedChannelDemo(value int, resultChan chan<- int) {
	go func() {
		// simulate some work
		time.Sleep(1 * time.Second)
		resultChan <- value
	}()
}

// TODO: Implement worker pool pattern
// - Define a Job struct
type Job struct {
	ID   int
	Data int
}

// - Create processJob function
func processJob(job Job) int {
	// Simulate processing
	time.Sleep(1000 * time.Millisecond)
	return job.Data * 2 // Example processing: just double the data
}

// - Launch worker goroutines
func worker(id int, jobs <-chan Job, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job.ID)
		result := processJob(job)
		results <- result
	}
}

// - Send jobs and collect results
func workerPoolDemo(numWorkers, numJobs int) { // just expain the flow of data
	jobs := make(chan Job, numJobs)
	results := make(chan int, numJobs)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Data: j * 10}
	}
	close(jobs)
	fmt.Println("jobs channel closed")
	for a := 1; a <= numJobs; a++ {
		result := <-results
		fmt.Printf("Result received: %d\n", result)
	}
}

func workerPoolDemo2(numWorkers, numJobs int) {
    jobs := make(chan Job)
    results := make(chan int)

    var wg sync.WaitGroup
    // Launch workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for job := range jobs {
                result := processJob(job)
                results <- result
            }
        }(w)
    }

    // producer: send jobs then close jobs
    go func() {
        for j := 1; j <= numJobs; j++ {
            jobs <- Job{ID: j, Data: j * 10}
        }
        close(jobs)
    }()

    // collector: wait for workers then close results
    go func() {
        wg.Wait()
        close(results)
    }()

    // consume results until channel closed
    for result := range results {
        fmt.Printf("Result received: %d\n", result)
    }
}

// TODO: Implement select with timeout function
func selectWithTimeoutDemo(wait time.Duration) {
	responseChan := make(chan string)
	go func() {
		// Simulate variable response time
		time.Sleep(3 * time.Second)
		responseChan <- "Response received"
	}()

	select {
	case res := <-responseChan:
		fmt.Println(res)
	case <-time.After(wait * time.Second):
		fmt.Println("Timeout: no response received within 2 seconds")
	}
}

// TODO (Optional): Implement pipeline (generate -> square -> print)
func calcSquare(num int, squared chan int) {
	squared <- num * num
}

func pipelineDemo() {
	squared := make(chan int, 10)

	for i := 1; i <= 10; i++ {
		go calcSquare(i, squared)
	}

	for i := 1; i <= 10; i++ {
		result := <-squared
		fmt.Println("Square:", result)
	}

}

func main() {
	fmt.Println("=== Day 4: Concurrency - Goroutines, Channels & Sync ===\n")

	// TODO: Test basic goroutines
	// - Launch PrintNumbers as goroutine
	go PrintNumbers(10)
	// - Use time.Sleep to wait (not ideal, but for learning)
	time.Sleep(1 * time.Second)

	// TODO: Test channel communication
	// - Create result channel
	resultChan := make(chan int)
	// - Launch multiple Sum goroutines with different slices
	go Sum([]int{1, 2, 3, 4, 5}, resultChan)
	go Sum([]int{6, 7, 8, 9, 10}, resultChan)
	// - Collect results
	sum1 := <-resultChan
	sum2 := <-resultChan
	fmt.Println("Sum1:", sum1)
	fmt.Println("Sum2:", sum2)

	// why channel assigning and retriving syntax is same?

	// TODO: Test buffered vs unbuffered channels
	// - Show blocking behavior of unbuffered
	resultChan = make(chan int)

	fmt.Println("Starting unbuffered channel and waiting for result to be recieved")
	UnbufferedChannelDemo(resultChan)
	result := <-resultChan
	fmt.Println("Result from unbuffered channel after 2 seconds:", result)

	// - Show non-blocking behavior of buffered (until full)
	resultChan = make(chan int, 2)
	BufferedChannelDemo(42, resultChan)
	BufferedChannelDemo(84, resultChan)
	fmt.Println("Result from buffered channel 1:", <-resultChan)
	fmt.Println("Result from buffered channel 2:", <-resultChan)
	//channel is buffered so it will not block until buffer is full
	// making channel full so it will block on next send if any
	BufferedChannelDemo(84, resultChan)
	BufferedChannelDemo(36, resultChan)
	BufferedChannelDemo(29, resultChan)

	// TODO: Test worker pool
	// - Create job and result channels
	// - Launch workers
	// - Send jobs
	// - Collect results
	workerPoolDemo(10, 100)

	workerPoolDemo2(10, 100)

	// TODO: Test select with timeout
	// - Try with fast response
	selectWithTimeoutDemo(4)
	// - Try with slow response (timeout)
	selectWithTimeoutDemo(2)

	// TODO (Optional): Test pipeline
	// - Generate numbers
	// - Square them
	// - Print results
	pipelineDemo()

	fmt.Println("\nDay 4 tasks completed!")
}
