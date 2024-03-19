package tasks

import (
	"fmt"
	"job-scheduler/worker/models"
	"time"
)

func ProcessDocument(job models.Job) {
	// check status on job table
	isRunning := checkJobStatusRunning(job)
	if isRunning {
		fmt.Println("Job is already running, skipped exeuction.")
		return
	}

	updateJobRunning(job)
	executionId := createExecution(job)

	// Fake task only
	fmt.Printf("Processing document for Job ID %v...\n", job.ID)
	time.Sleep(5 * time.Second)
	fmt.Println("Done processing document.")

	// update task
	updateJobExecution(job, executionId)
}
