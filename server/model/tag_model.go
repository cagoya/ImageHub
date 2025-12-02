package model

import (
	"time"
)

type Tag struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TagName     string    `json:"tag_name" gorm:"not null;varchar(16)"`
	ImageID     int       `json:"image_id" gorm:"not null;index"`
	FromAI      bool      `json:"from_ai" gorm:"default:false"` // 是否为AI生成的标签
	Image       Image     `json:"image" gorm:"foreignKey:ImageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`
	UpdatedTime time.Time `json:"updated_time" gorm:"autoUpdateTime"`
}

type SetTagRequest struct {
	ImageID    int      `json:"image_id"`
	AITags     []string `json:"ai_tags"`     // AI生成的标签
	ManualTags []string `json:"manual_tags"` // 手动设置的标签
}
