package models

import (
	"time"
)

type Execution struct {
	ID          uint `gorm:"primarykey"`
	JobID       uint
	Job         *Job `json:",omitempty"`
	CreatedAt   time.Time
	CompletedAt *time.Time
	Status      string `gorm:"type:enum('in_progress', 'complete'); default:in_progress"`
}
