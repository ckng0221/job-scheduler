package main

import (
	"job-scheduler/scheduler/core"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func init() {
	godotenv.Load()
}

func main() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("*/5 * * * * *", core.PublishActiveJobs)
	c.Start()

	select {}
}

// CompileDaemon -command="./scheduler"
