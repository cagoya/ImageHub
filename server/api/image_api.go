package api

import (
	"path"
	"strconv"
	"time"

	"server/global"
	"server/model"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageHandler struct {
	imageService *service.ImageService
}

func NewImageHandler() *ImageHandler {
	db := global.GetDB()
	imageService := service.NewImageService(db)
	return &ImageHandler{imageService: imageService}
}

func (h *ImageHandler) UploadImage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	image, err := c.FormFile("image")
	if err != nil {
		utils.ServerError(c, "图片获取失败")
		return
	}
	// 存储原图
	storeName := uuid.New().String() + path.Ext(image.Filename)
	filePath := path.Join("./public", strconv.Itoa(userID.(int)), "origin", storeName)
	err = c.SaveUploadedFile(image, filePath)
	if err != nil {
		utils.ServerError(c, "图片保存失败")
		return
	}
	_, err = h.imageService.UploadImage(path.Base(image.Filename), storeName, userID.(int))
	if err != nil {
		utils.ServerError(c, "图片上传失败")
		return
	}
	utils.Success(c, "图片上传成功")
}

func (h *ImageHandler) SetImageTags(c *gin.Context) {
	var setTagRequest model.SetTagRequest
	if err := c.ShouldBindJSON(&setTagRequest); err != nil {
		utils.ServerError(c, "设置图片标签请求绑定失败")
		return
	}
	imageID := setTagRequest.ImageID
	tags := setTagRequest.Tags
	err := h.imageService.SetImageTags(imageID, tags)
	if err != nil {
		utils.ServerError(c, "设置图片标签失败")
		return
	}
	utils.Success(c, "设置图片标签成功")
}

func (h *ImageHandler) GetImageFile(c *gin.Context) {
	imageID, err := strconv.Atoi(c.Query("imageID"))
	if err != nil {
		utils.ServerError(c, "图片ID获取失败")
		return
	}
	filePath, err := h.imageService.GetImageFile(imageID)
	if err != nil {
		utils.ServerError(c, "获取图片文件失败")
		return
	}
	c.File(filePath)
}

func (h *ImageHandler) GetAllImages(c *gin.Context) {
	isAdmin, exists := c.Get("isAdmin")
	if !exists {
		utils.ServerError(c, "判断用户是否为管理员失败")
		return
	}
	if !isAdmin.(bool) {
		utils.Forbidden(c, "无权限访问该接口")
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.ServerError(c, "页码获取失败")
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		utils.ServerError(c, "页码大小获取失败")
		return
	}
	images, pageTotal, err := h.imageService.GetAllImages(page, pageSize)
	if err != nil {
		utils.ServerError(c, "获取所有图片失败")
		return
	}
	utils.Success(c, gin.H{
		"images":    images,
		"pageTotal": pageTotal,
	})
}

func (h *ImageHandler) GetImagesByUserID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.ServerError(c, "页码获取失败")
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		utils.ServerError(c, "页码大小获取失败")
		return
	}
	images, pageTotal, err := h.imageService.GetImagesByUserID(userID.(int), page, pageSize)
	if err != nil {
		utils.ServerError(c, "获取用户图片失败")
		return
	}
	utils.Success(c, gin.H{
		"images":    images,
		"pageTotal": pageTotal,
	})
}

func (h *ImageHandler) GetImagesByTag(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.ServerError(c, "页码获取失败")
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		utils.ServerError(c, "页码大小获取失败")
		return
	}
	tag := c.Query("tag")
	if tag == "" {
		utils.ServerError(c, "标签获取失败")
		return
	}
	images, pageTotal, err := h.imageService.GetImagesByTag(userID.(int), tag, page, pageSize)
	if err != nil {
		utils.ServerError(c, "获取图片标签失败")
		return
	}
	utils.Success(c, gin.H{
		"images":    images,
		"pageTotal": pageTotal,
	})
}

func (h *ImageHandler) GetImagesByTimeRange(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.ServerError(c, "页码获取失败")
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		utils.ServerError(c, "页码大小获取失败")
		return
	}
	startTime, err := time.Parse(time.DateTime, c.Query("startTime"))
	if err != nil {
		utils.ServerError(c, "开始时间解析失败")
		return
	}
	endTime, err := time.Parse(time.DateTime, c.Query("endTime"))
	if err != nil {
		utils.ServerError(c, "结束时间解析失败")
		return
	}
	images, pageTotal, err := h.imageService.GetImagesByTimeRange(userID.(int), startTime, endTime, page, pageSize)
	if err != nil {
		utils.ServerError(c, "获取图片时间范围失败")
		return
	}
	utils.Success(c, gin.H{
		"images":    images,
		"pageTotal": pageTotal,
	})
}

func (h *ImageHandler) GetImagesByProvince(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.ServerError(c, "页码获取失败")
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		utils.ServerError(c, "页码大小获取失败")
		return
	}
	province := c.Query("province")
	images, pageTotal, err := h.imageService.GetImagesByProvince(userID.(int), province, page, pageSize)
	if err != nil {
		utils.ServerError(c, "获取图片省份失败")
		return
	}
	utils.Success(c, gin.H{
		"images":    images,
		"pageTotal": pageTotal,
	})
}

func (h *ImageHandler) GetImagesByCity(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.ServerError(c, "页码获取失败")
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		utils.ServerError(c, "页码大小获取失败")
		return
	}
	city := c.Query("city")
	images, pageTotal, err := h.imageService.GetImagesByCity(userID.(int), city, page, pageSize)
	if err != nil {
		utils.ServerError(c, "获取图片城市失败")
		return
	}
	utils.Success(c, gin.H{
		"images":    images,
		"pageTotal": pageTotal,
	})
}

func (h *ImageHandler) GetImagesByDistrict(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.ServerError(c, "页码获取失败")
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		utils.ServerError(c, "页码大小获取失败")
		return
	}
	district := c.Query("district")
	images, pageTotal, err := h.imageService.GetImagesByDistrict(userID.(int), district, page, pageSize)
	if err != nil {
		utils.ServerError(c, "获取图片区县失败")
		return
	}
	utils.Success(c, gin.H{
		"images":    images,
		"pageTotal": pageTotal,
	})
}

func (h *ImageHandler) UpdateImage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	imageID, err := strconv.Atoi(c.Query("imageID"))
	if err != nil {
		utils.ServerError(c, "图片ID获取失败")
		return
	}
	image, err := c.FormFile("image")
	if err != nil {
		utils.ServerError(c, "图片获取失败")
		return
	}
	storeName := uuid.New().String() + path.Ext(image.Filename)
	filePath := path.Join("./public", strconv.Itoa(userID.(int)), "origin", storeName)
	err = c.SaveUploadedFile(image, filePath)
	if err != nil {
		utils.ServerError(c, "图片保存失败")
		return
	}
	err = h.imageService.UpdateImage(storeName, userID.(int), imageID)
	if err != nil {
		utils.ServerError(c, "更新图片失败")
		return
	}
	utils.Success(c, "更新图片成功")
}

func (h *ImageHandler) DeleteImage(c *gin.Context) {
	operatorID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "操作者信息获取失败")
		return
	}
	userID, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		utils.ServerError(c, "用户ID获取失败")
		return
	}
	isAdmin, exists := c.Get("isAdmin")
	if !exists {
		utils.ServerError(c, "判断用户是否为管理员失败")
		return
	}
	if !isAdmin.(bool) && operatorID.(int) != userID {
		utils.Forbidden(c, "无权限访问该接口")
		return
	}
	imageID, err := strconv.Atoi(c.Query("imageID"))
	if err != nil {
		utils.ServerError(c, "图片ID获取失败")
		return
	}
	err = h.imageService.DeleteImage(userID, imageID)
	if err != nil {
		utils.ServerError(c, "删除图片失败")
		return
	}
	utils.Success(c, "删除图片成功")
}

func (h *ImageHandler) GetImageDetail(c *gin.Context) {
	imageID, err := strconv.Atoi(c.Query("imageID"))
	if err != nil {
		utils.ServerError(c, "图片ID获取失败")
		return
	}
	imageDetail, err := h.imageService.GetImageDetail(imageID)
	if err != nil {
		utils.ServerError(c, "获取图片详细信息失败")
		return
	}
	utils.Success(c, imageDetail)
}
