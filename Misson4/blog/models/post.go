package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 文章模型
type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"type:varchar(200);not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comments  []Comment      `json:"comments,omitempty" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

