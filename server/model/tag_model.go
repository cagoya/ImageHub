package model

import (
	"time"
)

type Tag struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TagName     string    `json:"tag_name" gorm:"not null;varchar(16)"`
	ImageID     int       `json:"image_id" gorm:"not null;index"`
	Image       Image     `json:"image" gorm:"foreignKey:ImageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`
	UpdatedTime time.Time `json:"updated_time" gorm:"autoUpdateTime"`
}
