package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	JobName     string    `gorm:"type:varchar(255)"`
	IsRecurring bool      `gorm:"default:false"`
	NextRunTime time.Time `json:",omitempty"`
	UserID      uint
	User        *User  `json:",omitempty"`
	Cron        string `gorm:"type:varchar(20)"`
}

type Execution struct {
	ID          uint `gorm:"primarykey"`
	JobID       uint
	Job         *Job
	CreatedAt   time.Time
	CompletedAt time.Time
	Status      string `gorm:"type:enum('in_progress', 'complete'); default:in_progress"`
}
