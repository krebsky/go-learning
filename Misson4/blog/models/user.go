package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"` // 密码不返回给客户端
	Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Posts     []Post         `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	Comments  []Comment      `json:"comments,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

