package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	JobName            string `gorm:"type:varchar(255)"`
	IsRecurring        bool   `gorm:"default:false"`
	FirstScheduledTime *int64 ``
	NextRunTime        *int64 ``
	UserID             uint
	User               *User  `json:",omitempty"`
	Cron               string `gorm:"type:varchar(20)"`
	IsCompleted        bool   `gorm:"default:false"`
	IsRunning          bool   `gorm:"default:false"`
	IsDisabled         bool   `gorm:"default:false"`
	RetryCount         uint16 `gorm:"default:0"`
	TaskPath           string
}

type JobUpdate struct {
	JobName            *string
	IsRecurring        *bool
	FirstScheduledTime *int64
	NextRunTime        *int64
	UserID             *uint
	Cron               string
	IsCompleted        *bool
	IsRunning          *bool
	IsDisabled         *bool
	RetryCount         uint16
}
