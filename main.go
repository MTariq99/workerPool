package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Job struct {
	ID       int
	UserName string
	Email    string
}

type Result struct {
	JobID   int
	Valid   bool
	Email   string
	Message string
}

func ValidateUser(job Job) Result {
	var issues []string
	if !strings.Contains(job.Email, "@") {
		issues = append(issues, "Invalid Email")
	}
	if len(issues) > 0 {
		return Result{
			JobID:   job.ID,
			Valid:   false,
			Email:   job.Email,
			Message: strings.Join(issues, ","),
		}
	}
	return Result{
		JobID:   job.ID,
		Valid:   true,
		Email:   job.Email,
		Message: "Valid Email",
	}

}

func workerpool(ctx context.Context, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done() //Signilizing it finished
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			result := ValidateUser(job)
			select {
			case results <- result:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// collector recieves and processes the results
func collector(results <-chan Result, done chan<- struct{}) {
	valid, invalid := 0, 0
	for result := range results {
		if result.Valid {
			valid++
			fmt.Printf("job %d:%s\n", result.JobID, result.Message)
		} else {
			invalid++
			fmt.Printf("job %d:%s Email:%s\n", result.JobID, result.Message, result.Email)
		}
	}
	fmt.Printf("\nprocessing completed:\n-Valid: %d\n-Invalid: %d\n", valid, invalid)
	close(done)
}

func dispatcher(ctx context.Context, workerCount int, data []Job) {
	jobCount := len(data)
	jobs := make(chan Job, jobCount)
	results := make(chan Result, jobCount)
	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(workerCount)
	for w := 1; w <= workerCount; w++ {
		go workerpool(ctx, jobs, results, wg)
	}

	go collector(results, done)

	fmt.Printf("Starting user validation with worker:%d\n\n", workerCount)

	for _, job := range data {
		jobs <- job
	}
	close(jobs)
	wg.Wait()
	close(results)
	<-done
}

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
	fmt.Println("Starting batch processing with syncronized results......")
	workerCount := 3
	dispatcher(ctx, workerCount, data)
}
