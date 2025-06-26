package model

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID        uint   `gorm:"not null"`
	ReceiverName  string `gorm:"not null"`
	Phone         string `gorm:"not null"`
	AddressDetail string `gorm:"not null"`
	IsDefault     bool   `gorm:"default:false"`
}
