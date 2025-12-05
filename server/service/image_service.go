package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"server/global"
	"server/model"
	"server/utils"

	"github.com/qdrant/go-client/qdrant"
	"gorm.io/gorm"
)

type ImageService struct {
	db           *gorm.DB
	qdrantClient *qdrant.Client
}

func NewImageService(db *gorm.DB, qdrantClient *qdrant.Client) *ImageService {
	return &ImageService{db: db, qdrantClient: qdrantClient}
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
	// 确保缺省字段，尤其是数的默认值是null而不是NaN
	image := &model.Image{
		Name:             imageName,
		Path:             storeName,
		UserID:           userID,
		CapturedTime:     nil,
		Latitude:         nil,
		Longitude:        nil,
		FormattedAddress: nil,
		Province:         nil,
		City:             nil,
		District:         nil,
		WidthPixel:       nil,
		HeightPixel:      nil,
	}

	// 检查用户是否开启了EXIF解析
	var user model.User
	err = s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return 0, fmt.Errorf("读出用户失败: %v", err)
	}
	// 如果用户开启了EXIF解析，则解析EXIF信息
	if user.ParseExif {
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
	}

	err = s.db.Create(&image).Error
	if err != nil {
		return 0, fmt.Errorf("保存到数据库失败: %v", err)
	}

	// 将向量化任务推入队列（异步处理）
	err = EnqueueVectorTask(image.ID, filePath)
	if err != nil {
		// 入队失败不影响主流程，只记录错误
		fmt.Printf("向量化任务入队失败: %v\n", err)
	}

	return image.ID, nil
}

// 设置图片标签
func (s *ImageService) SetImageTags(imageID int, aiTags []string, manualTags []string) error {
	// 先删除旧的标签
	err := s.db.Where("image_id = ?", imageID).Delete(&model.Tag{}).Error
	if err != nil {
		return fmt.Errorf("删除旧的标签失败: %v", err)
	}
	// 再创建新的AI生成标签
	for _, tag := range aiTags {
		tag := &model.Tag{
			ImageID: imageID,
			TagName: tag,
			FromAI:  true,
		}
		err := s.db.Create(tag).Error
		if err != nil {
			return fmt.Errorf("保存AI标签到数据库失败: %v", err)
		}
	}
	// 再创建新的手动设置标签
	for _, tag := range manualTags {
		tag := &model.Tag{
			ImageID: imageID,
			TagName: tag,
			FromAI:  false,
		}
		err := s.db.Create(tag).Error
		if err != nil {
			return fmt.Errorf("保存手动标签到数据库失败: %v", err)
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

// 获取图片的标签，区分AI生成和手动设置
func (s *ImageService) GetImageTags(imageID int) ([]string, []string, error) {
	var aiTags []string
	var manualTags []string
	var tags []model.Tag
	err := s.db.Where("image_id = ?", imageID).Find(&tags).Error
	if err != nil {
		return nil, nil, err
	}
	for _, tag := range tags {
		if tag.FromAI {
			aiTags = append(aiTags, tag.TagName)
		} else {
			manualTags = append(manualTags, tag.TagName)
		}
	}
	return aiTags, manualTags, nil
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

// 根据图片ID获取图片
func (s *ImageService) GetImageByID(imageID int) (*model.ImageDetailResponse, error) {
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

	// 将向量化任务推入队列（异步处理）
	err = EnqueueVectorTask(imageID, newImagePath)
	if err != nil {
		// 入队失败不影响主流程，只记录错误
		fmt.Printf("向量化任务入队失败: %v\n", err)
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

	// 删除 Qdrant 中的向量
	qdrantClient := global.GetQdrantClient()
	if qdrantClient != nil {
		_, err = qdrantClient.Delete(context.Background(), &qdrant.DeletePoints{
			CollectionName: "images",
			Points: &qdrant.PointsSelector{
				PointsSelectorOneOf: &qdrant.PointsSelector_Points{
					Points: &qdrant.PointsIdsList{
						Ids: []*qdrant.PointId{
							qdrant.NewIDNum(uint64(imageID)),
						},
					},
				},
			},
		})
		if err != nil {
			// Qdrant 删除失败不影响主流程，只记录错误
			fmt.Printf("删除 Qdrant 向量失败: %v\n", err)
		}
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
	aiTags, manualTags, err := s.GetImageTags(image.ID)
	if err != nil {
		return nil, err
	}
	return &model.ImageQueryResponse{
		ImageID:    image.ID,
		ImageName:  image.Name,
		ImagePath:  image.Path,
		AITags:     aiTags,
		ManualTags: manualTags,
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

// storeImageVector 生成图片向量并存储到 Qdrant
func (s *ImageService) storeImageVector(imageID int, filePath string) error {
	// 生成图片向量
	vectors, err := utils.VectorizeImageFromFile(filePath)
	if err != nil {
		return fmt.Errorf("生成图片向量失败: %v", err)
	}

	// 获取 Qdrant 客户端
	qdrantClient := global.GetQdrantClient()
	if qdrantClient == nil {
		return fmt.Errorf("qdrant 客户端未初始化")
	}

	// 将 []float64 转换为 []float32（Qdrant 需要 float32）
	vector32 := make([]float32, len(vectors))
	for i, v := range vectors {
		vector32[i] = float32(v)
	}

	// 构建 payload，包含图片ID和用户ID等信息
	// 从数据库获取图片信息以获取用户ID
	var image model.Image
	err = s.db.Where("id = ?", imageID).First(&image).Error
	if err != nil {
		return fmt.Errorf("获取图片信息失败: %v", err)
	}

	// 存储向量到 Qdrant
	operationInfo, err := qdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "images",
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(uint64(imageID)),
				Vectors: qdrant.NewVectors(vector32...),
				Payload: qdrant.NewValueMap(map[string]any{
					"user_id": image.UserID,
					"path":    image.Path,
				}),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("存储向量到 Qdrant 失败: %v", err)
	}
	if operationInfo != nil {
		fmt.Printf("Qdrant操作信息为: %v\n", operationInfo)
	}

	return nil
}
