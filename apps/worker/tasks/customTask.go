package tasks

import (
	"fmt"
	"job-scheduler/worker/models"
	"time"
)

func RunCustomJob(job models.Job) error {
	// check status on job table
	isRunning := checkJobStatusRunning(job)
	if isRunning {
		fmt.Println("Job is already running, skipped exeuction.")
		return nil
	}
	updateJobRunning(job, true)
	executionId, err := createExecution(job)
	if err != nil {
		fmt.Println("Failed to create execution")
		return err
	}

	fmt.Printf("Processing document for Job ID %v...\n", job.ID)
	err = runJob(job)
	defer updateJobRunning(job, false)

	if err != nil {
		// update execution to failed
		updateExecutionStatus(executionId, "failed")
		updateRetryCount(job)
		return err

	}
	fmt.Println("Done executing users' job.")

	// update task
	updateJobExecution(job, executionId)

	return err
}

func runJob(job models.Job) error {
	// TODO: get script and run script
	// Fake task only
	time.Sleep(5 * time.Second)
	return nil
}
