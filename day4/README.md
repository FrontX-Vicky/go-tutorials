# Day 4 — Concurrency: Goroutines, Channels & Sync

## Concepts to Learn

### 1. Goroutines

Lightweight threads managed by the Go runtime.

- **Launch with `go`**:
  ```go
  go someFunction()
  go func() {
      fmt.Println("Anonymous goroutine")
  }()
  ```
- **Important**: The main function doesn't wait for goroutines to finish unless you coordinate them.

### 2. Channels

Channels enable communication between goroutines.

- **Unbuffered channels** (synchronous):
  ```go
  ch := make(chan int)
  go func() {
      ch <- 42  // send
  }()
  value := <-ch  // receive
  ```
  
- **Buffered channels** (asynchronous):
  ```go
  ch := make(chan int, 3)  // buffer size 3
  ch <- 1  // won't block until buffer is full
  ch <- 2
  ch <- 3
  ```

- **Close channels**:
  ```go
  close(ch)
  value, ok := <-ch  // ok is false when channel is closed
  ```

- **Range over channels**:
  ```go
  for value := range ch {
      fmt.Println(value)
  }
  ```

### 3. Select Statement

Choose between multiple channel operations.

```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
}
```

### 4. WaitGroups

Coordinate multiple goroutines.

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Printf("Worker %d done\n", id)
    }(i)
}

wg.Wait()  // blocks until all goroutines call Done()
```

### 5. Common Patterns

- **Worker Pool**: Multiple workers processing from a shared channel.
- **Pipeline**: Chain of stages where each stage receives from one channel and sends to another.
- **Fan-out/Fan-in**: Distribute work to multiple goroutines, then collect results.

## Tasks

### 1. Basic Goroutines

Write a function `PrintNumbers(n int)` that prints numbers 1 to n. Launch it as a goroutine and use `time.Sleep` to wait for it to complete (not ideal, but for learning).

### 2. Channel Communication

Write a function `Sum(nums []int, resultChan chan int)` that:
- Takes a slice of numbers and a channel
- Calculates the sum
- Sends the result to the channel

Test it by launching multiple goroutines to sum different slices.

### 3. Buffered vs Unbuffered

Create examples showing:
- An unbuffered channel (blocks until receiver is ready)
- A buffered channel (doesn't block until full)

### 4. Worker Pool

Implement a worker pool pattern:
- Create a job channel and a results channel
- Launch N workers that read from the job channel, process jobs, and send results
- Send jobs to the job channel
- Collect results

### 5. Select with Timeout

Write a function that:
- Sends data to a channel after a random delay
- Uses `select` to receive with a timeout
- Reports whether data arrived or timeout occurred

## Extra Challenge

Implement a **Pipeline**:
- Stage 1: Generate numbers 1-10, send to channel
- Stage 2: Read from channel, square each number, send to another channel
- Stage 3: Read from channel, print results

Use goroutines for each stage and proper channel closing.
