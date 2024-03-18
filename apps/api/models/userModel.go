package models

import "gorm.io/gorm"

type Role string

const (
	Admin  Role = "admin"
	Member Role = "member"
)

type User struct {
	gorm.Model
	Name       string `gorm:"type:varchar(100)"`
	Email      string `gorm:"type:varchar(100); unique"`
	Password   string `gorm:"type:varchar(255)" json:"-"`
	Role       Role   `gorm:"type:enum('admin', 'member'); default:member"`
	ProfilePic string `gorm:"type:varchar(255)"`
}

var Roles = [2]string{string(Admin), string(Member)}
