package models

import "time"

type Job struct {
	ID          uint
	JobName     string
	IsRecurring bool
	NextRunTime int64
	UserID      uint
	Cron        string
	IsCompleted bool
	IsRunning   bool
	IsDisabled  bool
	RetryCount  uint16
}

type Execution struct {
	ID          uint
	JobID       uint
	Job         *Job
	CreatedAt   time.Time
	CompletedAt *time.Time
	Status      string
}
