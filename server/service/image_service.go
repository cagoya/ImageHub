package service

import (
	"fmt"
	"math"
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
func (s *ImageService) UploadImage(imageName string, userID int) (int, error) {
	// 生成缩略图
	filePath := path.Join("./public", strconv.Itoa(userID), "origin", imageName)
	err := utils.GenerateThumbnail(filePath, path.Join("./public", strconv.Itoa(userID), "thumbnail"))
	if err != nil {
		return 0, fmt.Errorf("生成缩略图失败: %v", err)
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
	err = s.db.Create(&image).Error
	if err != nil {
		return 0, fmt.Errorf("保存到数据库失败: %v", err)
	}
	return image.ID, nil
}

// 设置图片标签
func (s *ImageService) SetImageTags(imageID int, tags []string) error {
	for _, tag := range tags {
		tag := &model.Tag{
			ImageID: imageID,
			TagName: tag,
		}
		err := s.db.Create(tag).Error
		if err != nil {
			return fmt.Errorf("保存到数据库失败: %v", err)
		}
	}
	return nil
}

// 获取所有图片，分页
func (s *ImageService) GetAllImages(page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	err := s.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	pageTotal := int64(math.Ceil(float64(count) / float64(pageSize)))
	return images, pageTotal, nil
}

// 根据用户ID获取图片，分页
func (s *ImageService) GetImagesByUserID(userID int, page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	err := s.db.Where("user_id = ?", userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	pageTotal := int64(math.Ceil(float64(count) / float64(pageSize)))
	return images, pageTotal, nil
}

// 根据图片Tag获取图片，分页
func (s *ImageService) GetImagesByTag(tag string, page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	// tag 这里认为是精确的
	err := s.db.Model(&model.Image{}).Joins("JOIN tags ON images.id = tags.image_id").
		Where("tags.tag_name = ?", tag).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Joins("JOIN tags ON images.id = tags.image_id").
		Where("tags.tag_name = ?", tag).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	pageTotal := int64(math.Ceil(float64(count) / float64(pageSize)))
	return images, pageTotal, nil
}

// 根据时间范围获取图片，分页
func (s *ImageService) GetImagesByTimeRange(startTime, endTime time.Time, page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	err := s.db.Where("captured_time BETWEEN ? AND ?", startTime, endTime).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Where("captured_time BETWEEN ? AND ?", startTime, endTime).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	pageTotal := int64(math.Ceil(float64(count) / float64(pageSize)))
	return images, pageTotal, nil
}

// 根据省份获取图片，分页
func (s *ImageService) GetImagesByProvince(province string, page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	err := s.db.Where("province = ?", province).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Where("province = ?", province).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	pageTotal := int64(math.Ceil(float64(count) / float64(pageSize)))
	return images, pageTotal, nil
}

// 根据城市获取图片，分页
func (s *ImageService) GetImagesByCity(city string, page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	err := s.db.Where("city = ?", city).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Where("city = ?", city).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	pageTotal := int64(math.Ceil(float64(count) / float64(pageSize)))
	return images, pageTotal, nil
}

// 根据区县获取图片，分页
func (s *ImageService) GetImagesByDistrict(district string, page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	err := s.db.Where("district = ?", district).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Where("district = ?", district).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	pageTotal := int64(math.Ceil(float64(count) / float64(pageSize)))
	return images, pageTotal, nil
}

// 更新图片，由于用户的编辑，exif信息只更新像素和名称，其它的不变
func (s *ImageService) UpdateImage(imageName string, userID int, imageID int) error {
	// 读取exif信息
	filePath := path.Join("./public", strconv.Itoa(userID), "origin", imageName)
	exif := utils.GetExif(filePath)
	exifInfo := utils.ParseExif(exif)

	// 更新数据库
	image := &model.Image{
		Name:        imageName,
		WidthPixel:  getIntPointer(exifInfo.WidthPixel),
		HeightPixel: getIntPointer(exifInfo.HeightPixel),
	}
	err := s.db.Model(&model.Image{}).Where("id = ?", imageID).Updates(image).Error
	if err != nil {
		return fmt.Errorf("更新图片失败: %v", err)
	}
	return nil
}

// 删除图片
func (s *ImageService) DeleteImage(imageName string, userID int) error {
	err := s.db.Where("name = ? AND user_id = ?", imageName, userID).Delete(&model.Image{}).Error
	if err != nil {
		return fmt.Errorf("删除图片失败: %v", err)
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
