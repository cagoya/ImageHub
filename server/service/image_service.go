package service

import (
	"fmt"
	"path"
	"strconv"
	"time"

	"server/model"
	"server/utils"

	"gorm.io/gorm"
)

type ImageService struct {
	db *gorm.DB
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{db: db}
}

// 上传图片，自动解析exif信息
// 并生成缩略图
// 图片原图存储在 ./public/userID/orgin/ 目录下
// 缩略图存储在 ./public/userID/thumbnail/ 目录下
func (s *ImageService) UploadImage(imageName string, userID int) error {
	// 生成缩略图
	filePath := path.Join("./public", strconv.Itoa(userID), "origin", imageName)
	err := utils.GenerateThumbnail(filePath)
	if err != nil {
		return fmt.Errorf("生成缩略图失败: %v", err)
	}

	// 解析exif信息，失败也不会影响图片存储
	exif := utils.GetExif(filePath)
	exifInfo := utils.ParseExif(exif)

	// 逆地理编码，失败也不会影响图片存储
	geoInfo := utils.ReverseGeocode(exifInfo.Longitude, exifInfo.Latitude)

	// 保存到数据库，处理可选字段
	image := &model.Image{
		Name:             imageName,
		UserID:           userID,
		CapturedTime:     getTimePointer(exifInfo.CapturedTime),
		Latitude:         getFloat64Pointer(exifInfo.Latitude),
		Longitude:        getFloat64Pointer(exifInfo.Longitude),
		WidthPixel:       getIntPointer(exifInfo.WidthPixel),
		HeightPixel:      getIntPointer(exifInfo.HeightPixel),
		Province:         getStringPointer(geoInfo.Province),
		City:             getStringPointer(geoInfo.City),
		District:         getStringPointer(geoInfo.District),
		FormattedAddress: getStringPointer(geoInfo.FormattedAddress),
	}
	err = s.db.Create(image).Error
	if err != nil {
		return fmt.Errorf("保存到数据库失败: %v", err)
	}
	return nil
}

// 辅助函数，用于处理可选字段的指针类型
func getTimePointer(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func getFloat64Pointer(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

func getIntPointer(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func getStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// 设置图片标签

// 获取所有图片缩略图，分页

// 根据用户ID获取图片缩略图，分页

// 根据图片Tag获取图片缩略图，分页

// 根据时间范围获取图片缩略图，分页

// 根据地点获取图片缩略图，分页

// 根据地址获取图片缩略图，分页

// 根据图片ID获取图片原图

// 更新图片，由于用户的编辑

// 删除图片
