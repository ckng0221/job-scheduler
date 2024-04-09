package models

import "time"

type Status string

const (
	Complete Status = "complete"
	Failed   Status = "failed"
)

type Job struct {
	ID                 uint
	JobName            string
	IsRecurring        bool
	FirstScheduledTime int64
	NextRunTime        int64
	UserID             uint
	Cron               string
	IsCompleted        bool
	IsRunning          bool
	IsDisabled         bool
	RetryCount         uint16
	TaskPath           string
}

type Execution struct {
	ID          uint
	JobID       uint
	Job         *Job
	CreatedAt   time.Time
	CompletedAt *time.Time
	Status      Status
}
