package main

import (
	"job-scheduler/utils"

	"github.com/robfig/cron/v3"
)

func init() {
	requiredEnv := []string{"API_BASE_URL", "RABBIT_MQ_HOST", "JOB_QUEUE_NAME", "ADMIN_API_KEY"}
	utils.LoadEnv(requiredEnv)
}

func main() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("*/5 * * * * *", PublishActiveJobs)
	c.Start()

	select {}
}

// CompileDaemon -command="./scheduler"
