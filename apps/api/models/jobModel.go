package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	JobName     string `gorm:"type:varchar(255)"`
	IsRecurring bool   `gorm:"default:false"`
	NextRunTime *int64 ``
	UserID      uint
	User        *User  `json:",omitempty"`
	Cron        string `gorm:"type:varchar(20)"`
	IsCompleted bool   `gorm:"default:false"`
	IsRunning   bool   `gorm:"default:false"`
	IsDisabled  bool   `gorm:"default:false"`
	RetryCount  uint16 `gorm:"default:0"`
}

type JobUpdate struct {
	JobName     *string
	IsRecurring *bool
	NextRunTime *int64
	UserID      *uint
	Cron        string
	IsCompleted *bool
	IsRunning   *bool
	IsDisabled  *bool
	RetryCount  uint16
}
