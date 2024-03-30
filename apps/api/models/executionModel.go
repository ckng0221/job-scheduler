package models

import (
	"time"
)

type Status string

const (
	Complete Status = "complete"
	Failed   Status = "failed"
)

type Execution struct {
	ID          uint `gorm:"primarykey"`
	JobID       uint
	Job         *Job `json:",omitempty"`
	CreatedAt   time.Time
	CompletedAt *time.Time
	Status      Status `gorm:"type:enum('in_progress', 'complete', 'failed'); default:in_progress"`
}

type ExecutionUpdate struct {
	JobID       *uint
	CreatedAt   *time.Time
	CompletedAt *time.Time
	Status      *string
}
