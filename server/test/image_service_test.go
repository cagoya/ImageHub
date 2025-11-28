package test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"server/model"

	"github.com/stretchr/testify/assert"
)

func TestUploadImage(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users", "images"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)
	// 模拟存储图片到本地了
	err = CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test1.jpg")
	assert.NoError(t, err)
	imageID, err := imageService.UploadImage("test1.jpg", "test1.jpg", user.ID, true)
	assert.NoError(t, err)
	// 检查缩略图是否存在
	_, err = os.Stat("./public/" + strconv.Itoa(user.ID) + "/thumbnail/test1.jpg")
	assert.NoError(t, err)
	// 检查数据库中图片是否存在
	var image model.Image
	err = db.Where("id = ?", imageID).First(&image).Error
	assert.NoError(t, err)
	assert.Equal(t, image.Name, "test1.jpg")
	assert.Equal(t, image.UserID, user.ID)
}

func TestSetImageTags(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users", "images"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)
	// 模拟存储图片到本地了
	err = CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test1.jpg")
	assert.NoError(t, err)
	imageID, err := imageService.UploadImage("test1.jpg", "test1.jpg", user.ID, true)
	assert.NoError(t, err)
	// 检查数据库中图片是否存在
	var image model.Image
	err = db.Where("id = ?", imageID).First(&image).Error
	assert.NoError(t, err)
	assert.Equal(t, image.Name, "test1.jpg")
	assert.Equal(t, image.UserID, user.ID)
	// 设置图片标签
	err = imageService.SetImageTags(imageID, []string{"tag1", "tag2"})
	assert.NoError(t, err)
	// 检查数据库中图片标签是否存在
	var tag model.Tag
	err = db.Where("image_id = ? AND tag_name = ?", imageID, "tag1").Find(&tag).Error
	assert.NoError(t, err)
	err = db.Where("image_id = ? AND tag_name = ?", imageID, "tag2").Find(&tag).Error
	assert.NoError(t, err)
}

// 有些东西没必要反复测
// 所以后面重复的部分就不写断言了
// 让代码更简洁一些
func TestGetAllImages(t *testing.T) {
	err := ClearTables(dbConfig, []string{"tags", "images", "users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, _ := userService.RegisterUser(registerRequest)
	// 模拟存储图片到本地了
	for i := 1; i <= 4; i++ {
		CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test"+strconv.Itoa(i)+".jpg")
		imageService.UploadImage("test"+strconv.Itoa(i)+".jpg", "test"+strconv.Itoa(i)+".jpg", user.ID, true)
	}
	// 获取所有图片，第一页
	images, totalCount, err := imageService.GetAllImages(1, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 3)
	assert.Equal(t, totalCount, int64(4))
	// 获取第二页
	images, totalCount, err = imageService.GetAllImages(2, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 1)
	assert.Equal(t, totalCount, int64(4))
	// 获取超出范围页
	images, totalCount, err = imageService.GetAllImages(3, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 0)
	assert.Equal(t, totalCount, int64(4))
}

// 根据用户ID获取图片，分页
func TestGetImagesByUserID(t *testing.T) {
	err := ClearTables(dbConfig, []string{"tags", "images", "users"})
	if err != nil {
		panic(err)
	}
	registerRequest1 := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user1, _ := userService.RegisterUser(registerRequest1)
	registerRequest2 := &model.RegisterRequest{
		Username: "李四",
		Email:    "lisi@gmail.com",
		Password: "abcdef",
	}
	user2, _ := userService.RegisterUser(registerRequest2)
	// 模拟存储图片到本地了
	for i := 1; i <= 2; i++ {
		err := CopyFile("./images", "./public/"+strconv.Itoa(user1.ID)+"/origin", "test"+strconv.Itoa(i)+".jpg")
		assert.NoError(t, err)
		_, err = imageService.UploadImage("test"+strconv.Itoa(i)+".jpg", "test"+strconv.Itoa(i)+".jpg", user1.ID, true)
		assert.NoError(t, err)
	}
	for i := 3; i <= 4; i++ {
		err := CopyFile("./images", "./public/"+strconv.Itoa(user2.ID)+"/origin", "test"+strconv.Itoa(i)+".jpg")
		assert.NoError(t, err)
		_, err = imageService.UploadImage("test"+strconv.Itoa(i)+".jpg", "test"+strconv.Itoa(i)+".jpg", user2.ID, true)
		assert.NoError(t, err)
	}
	// 根据用户ID获取图片，用户1
	images, totalCount, err := imageService.GetImagesByUserID(user1.ID, 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 2)
	assert.Equal(t, totalCount, int64(2))
	// 根据用户ID获取图片，用户1，超出范围页
	images, totalCount, err = imageService.GetImagesByUserID(user1.ID, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 0)
	assert.Equal(t, totalCount, int64(2))
	// 根据用户ID获取图片，用户2
	images, totalCount, err = imageService.GetImagesByUserID(user2.ID, 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 2)
	assert.Equal(t, totalCount, int64(2))
	// 根据用户ID获取图片，用户2，超出范围页
	images, totalCount, err = imageService.GetImagesByUserID(user2.ID, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 0)
	assert.Equal(t, totalCount, int64(2))
}

// 根据图片Tag获取图片，分页
func TestGetImagesByTag(t *testing.T) {
	err := ClearTables(dbConfig, []string{"tags", "images", "users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, _ := userService.RegisterUser(registerRequest)
	// 模拟存储图片到本地了
	for i := 1; i <= 4; i++ {
		err := CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test"+strconv.Itoa(i)+".jpg")
		assert.NoError(t, err)
		imageID, err := imageService.UploadImage("test"+strconv.Itoa(i)+".jpg", "test"+strconv.Itoa(i)+".jpg", user.ID, true)
		assert.NoError(t, err)
		err = imageService.SetImageTags(imageID, []string{fmt.Sprintf("tag%d", i), fmt.Sprintf("tag%d", i+1)})
		assert.NoError(t, err)
	}
	// 根据图片Tag获取图片，第一页
	images, totalCount, err := imageService.GetImagesByTag(user.ID, "tag2", 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 2)
	assert.Equal(t, totalCount, int64(2))
	// 根据图片Tag获取图片，超出范围页
	images, totalCount, err = imageService.GetImagesByTag(user.ID, "tag2", 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 0)
	assert.Equal(t, totalCount, int64(2))
}

// 根据时间范围获取图片，分页
// test1.jpg 2024:08:07 13:54:29
// test2.jpg 2025:10:01 03:01:21
// test3.jpg 2025:10:02 12:19:29
// test4.jpg None
func TestGetImagesByTimeRange(t *testing.T) {
	err := ClearTables(dbConfig, []string{"tags", "images", "users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, _ := userService.RegisterUser(registerRequest)
	// 模拟存储图片到本地了
	for i := 1; i <= 4; i++ {
		CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test"+strconv.Itoa(i)+".jpg")
		imageService.UploadImage("test"+strconv.Itoa(i)+".jpg", "test"+strconv.Itoa(i)+".jpg", user.ID, true)
	}
	// 根据时间范围获取图片，第一页
	images, totalCount, err := imageService.GetImagesByTimeRange(user.ID, time.Date(2025, 9, 1, 3, 1, 21, 0, time.UTC),
		time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC), 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 2)
	assert.Equal(t, totalCount, int64(2))
	// 根据时间范围获取图片，超出范围页
	images, totalCount, err = imageService.GetImagesByTimeRange(user.ID, time.Date(2025, 9, 1, 3, 1, 21, 0, time.UTC),
		time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC), 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 0)
	assert.Equal(t, totalCount, int64(2))
}

// 根据省份获取图片，分页
func TestGetImagesByProvince(t *testing.T) {
	err := ClearTables(dbConfig, []string{"tags", "images", "users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, _ := userService.RegisterUser(registerRequest)
	for i := 1; i <= 4; i++ {
		CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test"+strconv.Itoa(i)+".jpg")
		imageService.UploadImage("test"+strconv.Itoa(i)+".jpg", "test"+strconv.Itoa(i)+".jpg", user.ID, true)
	}
	// 根据省份获取图片
	images, totalCount, err := imageService.GetImagesByArea(user.ID, "江苏省", "", "", 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 1)
	assert.Equal(t, totalCount, int64(1))
}

func TestGetImagesByCity(t *testing.T) {
	err := ClearTables(dbConfig, []string{"tags", "images", "users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, _ := userService.RegisterUser(registerRequest)
	for i := 1; i <= 4; i++ {
		CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test"+strconv.Itoa(i)+".jpg")
		imageService.UploadImage("test"+strconv.Itoa(i)+".jpg", "test"+strconv.Itoa(i)+".jpg", user.ID, true)
	}
	// 根据城市获取图片
	images, totalCount, err := imageService.GetImagesByArea(user.ID, "", "南京市", "", 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 1)
	assert.Equal(t, totalCount, int64(1))
}

func TestGetImagesByDistrict(t *testing.T) {
	err := ClearTables(dbConfig, []string{"tags", "images", "users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, _ := userService.RegisterUser(registerRequest)
	for i := 1; i <= 4; i++ {
		CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test"+strconv.Itoa(i)+".jpg")
		imageService.UploadImage("test"+strconv.Itoa(i)+".jpg", "test"+strconv.Itoa(i)+".jpg", user.ID, true)
	}
	// 根据区县获取图片
	images, totalCount, err := imageService.GetImagesByArea(user.ID, "", "", "鼓楼区", 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, len(images), 1)
	assert.Equal(t, totalCount, int64(1))
}

func TestUpdateImage(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users", "images"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, _ := userService.RegisterUser(registerRequest)
	CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test1.jpg")
	imageID, err := imageService.UploadImage("test1.jpg", "test1.jpg", user.ID, true)
	assert.NoError(t, err)
	// 获取旧的exif信息
	var image model.Image
	err = db.Where("id = ?", imageID).First(&image).Error
	assert.NoError(t, err)
	// 模拟存储test1_.jpg到本地了
	CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test1_.jpg")
	// 更新图片
	imageService.UpdateImage("test1_.jpg", user.ID, imageID)
	assert.NoError(t, err)
	// 检查旧的图片是否存在，应该不存在了
	_, err = os.Stat("./public/" + strconv.Itoa(user.ID) + "/origin/test1.jpg")
	assert.Error(t, err)
	_, err = os.Stat("./public/" + strconv.Itoa(user.ID) + "/thumbnail/test1.jpg")
	assert.Error(t, err)
	// 检查新的图片是否存在
	_, err = os.Stat("./public/" + strconv.Itoa(user.ID) + "/origin/test1_.jpg")
	assert.NoError(t, err)
	_, err = os.Stat("./public/" + strconv.Itoa(user.ID) + "/thumbnail/test1_.jpg")
	assert.NoError(t, err)
	// 获取新的exif信息
	var newImage model.Image
	err = db.Where("id = ?", imageID).First(&newImage).Error
	assert.NoError(t, err)
	assert.Equal(t, newImage.Name, "test1.jpg")
	assert.Equal(t, newImage.UserID, user.ID)
	assert.Equal(t, *newImage.CapturedTime, *image.CapturedTime)
	assert.Equal(t, *newImage.Longitude, *image.Longitude)
	assert.Equal(t, *newImage.Latitude, *image.Latitude)
	assert.Equal(t, *newImage.Province, *image.Province)
	assert.Equal(t, *newImage.City, *image.City)
	assert.Equal(t, *newImage.District, *image.District)
}

func TestDeleteImage(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users", "images"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, _ := userService.RegisterUser(registerRequest)
	CopyFile("./images", "./public/"+strconv.Itoa(user.ID)+"/origin", "test1.jpg")
	imageID, err := imageService.UploadImage("test1.jpg", "test1.jpg", user.ID, true)
	assert.NoError(t, err)
	// 删除图片
	err = imageService.DeleteImage(user.ID, imageID)
	assert.NoError(t, err)
	// 检查旧的图片是否存在，应该不存在了
	_, err = os.Stat("./public/" + strconv.Itoa(user.ID) + "/origin/test1.jpg")
	assert.Error(t, err)
	_, err = os.Stat("./public/" + strconv.Itoa(user.ID) + "/thumbnail/test1.jpg")
	assert.Error(t, err)
	// 检查数据库中图片是否存在
	var image model.Image
	err = db.Where("id = ?", imageID).First(&image).Error
	assert.Error(t, err)
}
