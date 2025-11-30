package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PostID    uint           `json:"post_id" gorm:"not null;index"`
	Post      Post           `json:"post,omitempty" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

