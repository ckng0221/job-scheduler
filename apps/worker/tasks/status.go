package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"job-scheduler/worker/models"
	"job-scheduler/worker/utils"
	"net/http"
	"os"
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

func updateNextRunTime(job models.Job) {
	fmt.Println("Updating Job ID", job.ID, "...")
	if job.ID == 0 {
		fmt.Println("Job ID cannot be null")
		return
	}

	API_BASE := os.Getenv("API_BASE_URL")
	endpoint := API_BASE + "/scheduler/jobs/" + fmt.Sprint(job.ID)
	payload := map[string]interface{}{}

	// Recuring job
	if job.IsRecurring {
		if job.Cron == "" {
			fmt.Println("Cron cannot be empty for recurring job.")
			return
		}
		cronExp, err := cronexpr.Parse(job.Cron)
		if err != nil {
			fmt.Println("Invalid cron expression")
			return
		}
		nextTime := cronExp.Next(time.Now().UTC()).Unix()
		nextTime2 := cronExp.Next(time.Now().UTC())
		payload["NextRunTime"] = nextTime
		fmt.Println("Next run time", nextTime2, "unix:", nextTime)
	} else {
		payload["IsCompleted"] = true
	}
	payload["IsRunning"] = false

	// Update
	payloadByte, _ := json.Marshal(payload)
	resp, err := utils.PatchRequest(endpoint, payloadByte)
	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Updated Job ID", job.ID, ".")
	} else {
		fmt.Println("Failed to update Job ID", job.ID, ".")
	}
}

func updateExecutionStatus(executionId uint) {
	API_BASE := os.Getenv("API_BASE_URL")
	endpoint := API_BASE + "/scheduler/executions/" + fmt.Sprint(executionId)
	payload := map[string]interface{}{
		"Status":      "complete",
		"CompletedAt": time.Now().UTC(),
	}
	payloadByte, _ := json.Marshal(payload)

	resp, err := utils.PatchRequest(endpoint, payloadByte)
	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode == 200 {
		fmt.Printf("Updated Execution ID %v.\n", executionId)
	} else {
		fmt.Printf("Failed to Execution ID %v.\n", executionId)
	}
}

func updateJobExecution(job models.Job, executionId uint) {
	updateExecutionStatus(executionId)
	updateNextRunTime(job)
}

func checkJobStatusRunning(job models.Job) bool {
	API_BASE := os.Getenv("API_BASE_URL")
	endpoint := API_BASE + "/scheduler/jobs/" + fmt.Sprint(job.ID)
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &job)
	fmt.Println("status", job.IsRunning)

	return job.IsRunning
}

func createExecution(job models.Job) uint {
	API_BASE := os.Getenv("API_BASE_URL")
	endpoint := API_BASE + "/scheduler/executions"

	payload := []map[string]interface{}{
		{"JobID": job.ID},
	}

	payloadByte, _ := json.Marshal(payload)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payloadByte))
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var executions []models.Execution
	json.Unmarshal(body, &executions)

	if resp.StatusCode == 201 {
		fmt.Println("Created Execution ID", executions[0].ID, ".")
	} else {
		fmt.Println("Failed to create execution for Job", job.ID, ".")
	}
	return executions[0].ID
}

func updateJobRunning(job models.Job) {
	API_BASE := os.Getenv("API_BASE_URL")
	endpoint := API_BASE + "/scheduler/jobs/" + fmt.Sprint(job.ID)
	payload := map[string]interface{}{
		"IsRunning": true,
	}
	payloadByte, _ := json.Marshal(payload)

	resp, _ := utils.PatchRequest(endpoint, payloadByte)

	if resp.StatusCode == 200 {
		fmt.Println("Updated Job ID running status", job.ID, ".")
	} else {
		fmt.Println("Failed to update Job ID running status", job.ID, ".")
	}
}