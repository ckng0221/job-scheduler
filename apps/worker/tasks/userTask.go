package tasks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"job-scheduler/worker/models"
	"job-scheduler/worker/utils"
	"os"
	"os/exec"
	"strings"
)

func RunUserTask(job models.Job) error {
	// check status on job table
	isRunning, err := checkJobStatusRunning(job)
	if err != nil {
		fmt.Println("Failed to check status", err)
		return err
	}

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
	API_BASE := os.Getenv("API_BASE_URL")
	endpoint := fmt.Sprintf("%s/scheduler/jobs/%s", API_BASE, fmt.Sprint(job.ID))
	resp, err := utils.GetRequest(endpoint)

	if err != nil {
		fmt.Println(err)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.Unmarshal(body, &job)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if job.TaskPath == "" {
		fmt.Println("No script to run. Skipped.")
		return nil
	}

	// Hardcoded relative filelocation
	err = runScript(job.TaskPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// FIXME: For POC only, dangerous to direcly run user-submitted script on the server.
func runScript(filePath string) error {
	if filePath == "" {
		fmt.Println("No script to run. Skipped.")
		return nil
	}
	fmt.Println("Running script...")

	// get file type
	filePathSplit := strings.Split(filePath, ".")
	fileExtension := filePathSplit[len(filePathSplit)-1]
	program, err := getProgramName(fileExtension)
	if err != nil {
		fmt.Println(err)
		return err

	}

	cmd := exec.Command(program, filePath)
	var out bytes.Buffer
	// define the process standard output
	cmd.Stdout = &out
	// Run the command
	err = cmd.Run()
	if err != nil {
		// error case : status code of command is different from 0
		fmt.Println("Shell err:", err)
		return err

	}
	fmt.Println(out.String())
	return nil
}

func getProgramName(fileExtension string) (string, error) {
	var program string

	// Supported file type
	switch fileExtension {
	case "sh":
		program = "sh"
	case "js":
		program = "node"
	}
	if program == "" {
		return "", errors.New("invalid file type")
	}
	return program, nil
}
