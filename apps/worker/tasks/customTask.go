package tasks

import (
	"fmt"
	"job-scheduler/worker/models"
	"time"
)

func RunCustomJob(job models.Job) {
	// check status on job table
	isRunning := checkJobStatusRunning(job)
	if isRunning {
		fmt.Println("Job is already running, skipped exeuction.")
		return
	}
	updateJobRunning(job)
	executionId := createExecution(job)

	fmt.Printf("Processing document for Job ID %v...\n", job.ID)
	runJob(job)
	fmt.Println("Done processing document.")

	// update task
	updateJobExecution(job, executionId)
}

func runJob(job models.Job) {
	// TODO: get script and run script
	// Fake task only
	time.Sleep(5 * time.Second)
}
