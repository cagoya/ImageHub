package service

import (
	"fmt"
	"os"
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

// 上传图片，并生成缩略图
// 这里返回值包含图片ID主要是为了方便测试
func (s *ImageService) UploadImage(imageName string, storeName string, userID int) (int, error) {
	// 生成缩略图，图片原图存储在 ./public/userID/orgin/ 目录下
	// 缩略图存储在 ./public/userID/thumbnail/ 目录下
	filePath := path.Join("./public", strconv.Itoa(userID), "origin", storeName)
	err := utils.GenerateThumbnail(filePath, path.Join("./public", strconv.Itoa(userID), "thumbnail"))
	if err != nil {
		return 0, fmt.Errorf("生成缩略图失败: %v", err)
	}

	// 保存到数据库，处理可选字段
	image := &model.Image{
		Name:   imageName,
		Path:   storeName,
		UserID: userID,
	}

	// 解析exif信息
	exif := utils.GetExif(filePath)
	exifInfo := utils.ParseExif(exif)
	// 逆地理编码
	geoInfo := utils.ReverseGeocode(exifInfo.Longitude, exifInfo.Latitude)

	image.CapturedTime = utils.GetTimePointer(exifInfo.CapturedTime)
	image.Latitude = utils.GetFloat64Pointer(exifInfo.Latitude)
	image.Longitude = utils.GetFloat64Pointer(exifInfo.Longitude)
	image.WidthPixel = utils.GetIntPointer(exifInfo.WidthPixel)
	image.HeightPixel = utils.GetIntPointer(exifInfo.HeightPixel)
	image.Province = utils.GetStringPointer(geoInfo.Province)
	image.City = utils.GetStringPointer(geoInfo.City)
	image.District = utils.GetStringPointer(geoInfo.District)
	image.FormattedAddress = utils.GetStringPointer(geoInfo.FormattedAddress)

	err = s.db.Create(&image).Error
	if err != nil {
		return 0, fmt.Errorf("保存到数据库失败: %v", err)
	}
	return image.ID, nil
}

// 设置图片标签
func (s *ImageService) SetImageTags(imageID int, tags []string) error {
	// 先删除旧的标签
	err := s.db.Where("image_id = ?", imageID).Delete(&model.Tag{}).Error
	if err != nil {
		return fmt.Errorf("删除旧的标签失败: %v", err)
	}
	// 再创建新的标签
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

// 由AI生成图片标签
func (s *ImageService) GenerateImageTags(imageID int) ([]string, error) {
	// 从数据库中读出图片元信息
	var image model.Image
	err := s.db.Where("id = ?", imageID).First(&image).Error
	if err != nil {
		return nil, err
	}
	// 读出图片文件，因为接口有图片大小限制，所以使用缩略图
	filePath := path.Join("./public", strconv.Itoa(image.UserID), "thumbnail", image.Path)
	imageFile64 := utils.GetFileContentAsBase64(filePath)
	// 调用百度AI生成图片标签
	tags, err := utils.ImageClassification(imageFile64)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// 获取图片的详细信息
func (s *ImageService) GetImageDetail(imageID int) (*model.ImageDetailResponse, error) {
	var image model.Image
	err := s.db.Where("id = ?", imageID).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &model.ImageDetailResponse{
		ID:               image.ID,
		Name:             image.Name,
		Path:             image.Path,
		CapturedTime:     image.CapturedTime,
		Latitude:         image.Latitude,
		Longitude:        image.Longitude,
		FormattedAddress: image.FormattedAddress,
		Province:         image.Province,
		City:             image.City,
		District:         image.District,
		WidthPixel:       image.WidthPixel,
		HeightPixel:      image.HeightPixel,
		UserID:           image.UserID,
		CreatedTime:      image.CreatedTime,
		UpdatedTime:      image.UpdatedTime,
	}, nil
}

// 获取图片的文件本身
func (s *ImageService) GetImageFile(imageID int) (string, error) {
	// 先从数据库中读出图片路径
	var image model.Image
	err := s.db.Where("id = ?", imageID).First(&image).Error
	if err != nil {
		return "", err
	}
	filePath := path.Join("./public", strconv.Itoa(image.UserID), "origin", image.Path)
	return filePath, nil
}

// 获取图片的标签
func (s *ImageService) GetImageTags(imageID int) ([]string, error) {
	var tags []string
	err := s.db.Model(&model.Tag{}).Where("image_id = ?", imageID).Pluck("tag_name", &tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// 获取所有图片，分页
func (s *ImageService) GetAllImages(page, pageSize int) ([]*model.ImageQueryResponse, int64, error) {
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

	responses, err := s.convertToImageQueryResponseList(images)
	if err != nil {
		return nil, 0, err
	}
	return responses, count, nil
}

// 根据用户ID获取图片，分页
func (s *ImageService) GetImagesByUserID(userID int, page, pageSize int) ([]*model.ImageQueryResponse, int64, error) {
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

	responses, err := s.convertToImageQueryResponseList(images)
	if err != nil {
		return nil, 0, err
	}
	return responses, count, nil
}

// 根据图片Tag获取图片，分页
func (s *ImageService) GetImagesByTag(userID int, tag string, page, pageSize int) ([]*model.ImageQueryResponse, int64, error) {
	var images []*model.Image
	// tag 这里认为是精确的
	err := s.db.Model(&model.Image{}).Joins("JOIN tags ON images.id = tags.image_id").
		Where("tags.tag_name = ? AND images.user_id = ?", tag, userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Joins("JOIN tags ON images.id = tags.image_id").
		Where("tags.tag_name = ? AND images.user_id = ?", tag, userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	responses, err := s.convertToImageQueryResponseList(images)
	if err != nil {
		return nil, 0, err
	}
	return responses, count, nil
}

// 根据时间范围获取图片，分页
func (s *ImageService) GetImagesByTimeRange(userID int, startTime, endTime time.Time, page, pageSize int) ([]*model.ImageQueryResponse, int64, error) {
	var images []*model.Image
	err := s.db.Where("captured_time BETWEEN ? AND ? AND user_id = ?", startTime, endTime, userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = s.db.Model(&model.Image{}).Where("captured_time BETWEEN ? AND ? AND user_id = ?", startTime, endTime, userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	responses, err := s.convertToImageQueryResponseList(images)
	if err != nil {
		return nil, 0, err
	}
	return responses, count, nil
}

// 根据地理位置获取图片，分页
func (s *ImageService) GetImagesByArea(userID int, province string, city string,
	district string, page, pageSize int) ([]*model.ImageQueryResponse, int64, error) {
	// 获取满足这些省市区的图片
	var images []*model.Image
	err := s.db.Where("(province = ? OR city = ? OR district = ?) AND user_id = ?",
		province, city, district, userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error
	if err != nil {
		return nil, 0, err
	}
	// 获取满足这些省市区的图片数量
	var count int64
	err = s.db.Model(&model.Image{}).Where("(province = ? OR city = ? OR district = ?) AND user_id = ?",
		province, city, district, userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	responses, err := s.convertToImageQueryResponseList(images)
	if err != nil {
		return nil, 0, err
	}
	return responses, count, nil
}

// 由于用户的编辑需要更新图片，exif信息只更新像素和名称，其它的不变
func (s *ImageService) UpdateImage(storeName string, userID int, imageID int) error {
	// 删除旧的图片
	var oldImage model.Image
	err := s.db.Where("id = ?", imageID).First(&oldImage).Error
	if err != nil {
		return fmt.Errorf("读出图片失败: %v", err)
	}
	oldImagePath := path.Join("./public", strconv.Itoa(userID), "origin", oldImage.Path)
	oldThumbnailPath := path.Join("./public", strconv.Itoa(userID), "thumbnail", oldImage.Path)
	os.Remove(oldImagePath)
	os.Remove(oldThumbnailPath)
	// 生成新的缩略图
	newImagePath := path.Join("./public", strconv.Itoa(userID), "origin", storeName)
	outputFolder := path.Join("./public", strconv.Itoa(userID), "thumbnail")
	err = utils.GenerateThumbnail(newImagePath, outputFolder)
	if err != nil {
		return fmt.Errorf("生成缩略图失败: %v", err)
	}

	// 更新数据库
	image := &model.Image{
		Path: storeName,
	}
	err = s.db.Model(&model.Image{}).Where("id = ?", imageID).Updates(image).Error
	if err != nil {
		return fmt.Errorf("更新图片失败: %v", err)
	}
	return nil
}

// 删除图片
func (s *ImageService) DeleteImage(userID int, imageID int) error {
	// 先从数据库中读出图片名称
	var image model.Image
	err := s.db.Where("id = ?", imageID).First(&image).Error
	if err != nil {
		return fmt.Errorf("读出图片失败: %v", err)
	}
	// 删除旧的图片
	oldImagePath := path.Join("./public", strconv.Itoa(userID), "origin", image.Path)
	oldThumbnailPath := path.Join("./public", strconv.Itoa(userID), "thumbnail", image.Path)
	// 删除旧的图片
	os.Remove(oldImagePath)
	os.Remove(oldThumbnailPath)
	// 删除数据库中的图片
	err = s.db.Where("id = ?", imageID).Delete(&model.Image{}).Error
	if err != nil {
		return fmt.Errorf("删除图片失败: %v", err)
	}
	return nil
}

// 将Image转换为ImageQueryResponse
func (s *ImageService) convertToImageQueryResponse(image *model.Image) (*model.ImageQueryResponse, error) {
	tags, err := s.GetImageTags(image.ID)
	if err != nil {
		return nil, err
	}
	return &model.ImageQueryResponse{
		ImageID:   image.ID,
		ImageName: image.Name,
		ImagePath: image.Path,
		Tags:      tags,
	}, nil
}

// 将Image数组转换为ImageQueryResponse数组
func (s *ImageService) convertToImageQueryResponseList(images []*model.Image) ([]*model.ImageQueryResponse, error) {
	var responses []*model.ImageQueryResponse
	for _, image := range images {
		response, err := s.convertToImageQueryResponse(image)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}
