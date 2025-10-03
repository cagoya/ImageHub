package api

import (
	"path"
	"strconv"

	"server/global"
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
	fileName := uuid.New().String() + "_" + path.Base(image.Filename)
	filePath := path.Join("./public", strconv.Itoa(userID.(int)), "origin", fileName)
	err = c.SaveUploadedFile(image, filePath)
	if err != nil {
		utils.ServerError(c, "图片保存失败")
		return
	}
	err = h.imageService.UploadImage(fileName, userID.(int))
	if err != nil {
		utils.ServerError(c, "图片上传失败")
		return
	}
	utils.Success(c, "图片上传成功")
}
