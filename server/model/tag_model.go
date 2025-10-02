package model

import (
	"time"
)

type Tag struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null;unique;varchar(16)"`
	ImageID     int       `json:"image_id" gorm:"not null;index"`
	Image       Image     `json:"image" gorm:"foreignKey:ImageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`
	UpdatedTime time.Time `json:"updated_time" gorm:"autoUpdateTime"`
}
