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
	Email      string `gorm:"type:varchar(100)"`
	Password   string `gorm:"type:varchar(255)" json:"-"`
	Role       Role   `gorm:"type:enum('admin', 'member'); default:member"`
	ProfilePic string `gorm:"type:varchar(255)"`
	Sub        string `gorm:"type:varchar(100); unique"`
}

type GoogleProfile struct {
	Sub            string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
}

var Roles = [2]string{string(Admin), string(Member)}
