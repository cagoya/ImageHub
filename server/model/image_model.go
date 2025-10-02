package model

import (
	"time"
)

type Image struct {
	ID           int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string     `json:"name" gorm:"not null;unique"`
	CapturedTime *time.Time `json:"captured_time"`
	// 因为这几个都有可能为空
	// 极端情况是经过社交媒体的照片会失去大部分exif信息
	// 所以这里用指针类型
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	WidthPixel  *int      `json:"width_pixel"`
	HeightPixel *int      `json:"height_pixel"`
	UserID      int       `json:"user_id" gorm:"not null;index"`
	User        User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // 外键关联
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`
	UpdatedTime time.Time `json:"updated_time" gorm:"autoUpdateTime"`
}
