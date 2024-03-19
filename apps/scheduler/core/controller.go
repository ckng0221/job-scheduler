package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

type Job struct {
	ID          uint
	JobName     string `gorm:"type:varchar(255)"`
	IsRecurring bool   `gorm:"default:false"`
	NextRunTime int64
	UserID      uint
	Cron        string `gorm:"type:varchar(20)"`
}

func getActiveJobs() []Job {
	fmt.Println("Reading active jobs...")

	API_BASE := os.Getenv("API_BASE_URL")
	endpoint := API_BASE + "/scheduler/jobs?active=true"

	res, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Received all active jobs")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jobs []Job
	json.Unmarshal(body, &jobs)

	return jobs
}

func PublishActiveJobs() {
	jobs := getActiveJobs()
	if len(jobs) == 0 {
		fmt.Println("No active jobs.")
		return
	}

	for i, job := range jobs {
		// TODO: add jobs to queue
		fmt.Println(i, job.ID)

		// TODO: move this function on execution level
		updateNextRunTime(job)
	}
}

func updateNextRunTime(job Job) {
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

	// Update
	payloadByte, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(payloadByte))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Updated Job ID", job.ID, ".")
	} else {
		fmt.Println("Failed to update Job ID", job.ID, ".")
	}
}
