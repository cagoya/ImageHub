package model

import (
	"time"
)

type Image struct {
	ID           int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string     `json:"name" gorm:"not null;varchar(256)"`
	Path         string     `json:"path" gorm:"not null;unique;varchar(256)"`
	CapturedTime *time.Time `json:"captured_time"`
	// 因为这几个都有可能为空
	// 极端情况是经过社交媒体的照片会失去大部分exif信息
	// 所以这里用指针类型
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	// 这个只有大陆地点能解析
	FormattedAddress *string   `json:"formatted_address" gorm:"type:varchar(128)"`
	Province         *string   `json:"province" gorm:"type:varchar(32)"`
	City             *string   `json:"city" gorm:"type:varchar(32)"`
	District         *string   `json:"district" gorm:"type:varchar(32)"`
	WidthPixel       *int      `json:"width_pixel"`
	HeightPixel      *int      `json:"height_pixel"`
	UserID           int       `json:"user_id" gorm:"not null;index"`
	User             User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedTime      time.Time `json:"created_time" gorm:"autoCreateTime"`
	UpdatedTime      time.Time `json:"updated_time" gorm:"autoUpdateTime"`
}

// 图片查询响应
// 标签是要一并返回
// 但是来自exif的信息仅在点击图片时显示
type ImageQueyResponse struct {
	ImageID   int      `json:"image_id"`
	ImageName string   `json:"image_name"`
	ImagePath string   `json:"image_path"`
	Tags      []string `json:"tags"`
}
