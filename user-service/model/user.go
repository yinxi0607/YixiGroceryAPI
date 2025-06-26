package model

import (
	"gorm.io/gorm"
	"time"
)

// User 定义用户模型
type User struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Address   string `gorm:"type:varchar(255)"` // 新增地址字段
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
