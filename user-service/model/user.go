package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Phone     string
	Address   string
	Points    int `gorm:"default:0"`
	CreatedAt time.Time
}
