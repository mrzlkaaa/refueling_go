package storage

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Surname   string `gorm:"not null"`
	Email     string `gorm:"unique; not null"`
	Username  string `gorm:"unique; not null"`
	Password  []byte
	Moderator bool `gorm:"default:false"`
	Admin     bool `gorm:"default:false"`
}
