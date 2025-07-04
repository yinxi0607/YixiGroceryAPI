package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Phone    string
	Address  string
	Points   int `gorm:"default:0"`
}
