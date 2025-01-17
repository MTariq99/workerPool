# Concurrent Worker Pool for User Email Validation

This project demonstrates the implementation of a **concurrent worker pool** in Go for validating user email addresses. It processes jobs in parallel using workers, validates emails, and aggregates the results.

---

## Features

- Efficiently processes multiple jobs using a **worker pool**.
- Validates email addresses for correctness.
- Tracks and reports the number of valid and invalid emails.
- Implements **graceful cancellation** using Go's `context` package.
- Uses **goroutines** and **channels** for concurrency.

---

## Project Structure

- **`main`**: Entry point of the application.
- **`ValidateUser`**: Validates a single job (email address).
- **`workerpool`**: Processes jobs using multiple workers.
- **`collector`**: Collects and aggregates results from workers.
- **`dispatcher`**: Manages job distribution and worker lifecycle.

---

## Flow of Execution

1. **Initialization**:

   - A slice of jobs (users and their emails) is created.
   - A context with a timeout of 5 seconds is set up.

2. **Dispatcher**:

   - Initializes channels for jobs and results.
   - Launches multiple workers to process jobs.
   - Starts a collector to handle results.

3. **Worker Execution**:

   - Each worker:
     - Reads jobs from the `jobs` channel.
     - Validates the job using `ValidateUser`.
     - Sends the result to the `results` channel.
   - Stops when the `jobs` channel is closed or the context is cancelled.

4. **Collector**:

   - Collects results from the `results` channel.
   - Tracks valid and invalid jobs.
   - Prints validation messages and a summary.

5. **Completion**:
   - The `jobs` channel is closed after all jobs are sent.
   - Workers finish their tasks and signal completion.
   - Results are processed, and a summary is printed.

---

## Functions

### `ValidateUser(job Job) Result`

Validates a single job (email). Returns a `Result` struct indicating whether the email is valid, along with a message.

### `workerpool(ctx context.Context, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup)`

Processes jobs in parallel. Workers validate each job and send results to the `results` channel.

### `collector(result <-chan Result, done chan<- struct{})`

Aggregates results from the `results` channel. Tracks the number of valid and invalid emails.

### `dispatcher(ctx context.Context, workerCount int, data []Job)`

Manages the worker pool and orchestrates job distribution. Waits for workers and results to complete.

---

## Example Usage

```go
func main() {
    ctx := context.Background()
    data := []Job{
        {ID: 1, UserName: "john", Email: "john@example.com"},
        {ID: 2, UserName: "ab", Email: "invalid-email"},
        {ID: 3, UserName: "jane", Email: "1-1@example.com"},
        {ID: 4, UserName: "doe", Email: "doeexample.com"},
        {ID: 5, UserName: "jane", Email: "jane.example.com"},
        {ID: 6, UserName: "x", Email: "x@gmail.com"},
        {ID: 7, UserName: "2e", Email: "2example.com"},
        {ID: 8, UserName: "ca", Email: "invalid"},
        {ID: 9, UserName: "dave", Email: "dave£example.com"},
        {ID: 10, UserName: "sameeralikhan", Email: "salikhan8458@gmail"},
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    fmt.Println("Starting batch processing with synchronized results......")
    workerCount := 3
    dispatcher(ctx, workerCount, data)
}
```
