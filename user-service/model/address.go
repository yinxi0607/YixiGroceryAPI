package model

import "time"

type Address struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `gorm:"not null"`
	ReceiverName  string `gorm:"not null"`
	Phone         string `gorm:"not null"`
	AddressDetail string `gorm:"not null"`
	IsDefault     bool   `gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
