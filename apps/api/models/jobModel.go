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
}
