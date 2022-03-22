package storage

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique; not null"`
	Password []byte
	Admin    bool `gorm:"default:false"`
}
