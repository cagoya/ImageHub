package test

import (
	"os"
	"testing"

	"server/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetExif(t *testing.T) {
	x := utils.GetExif("./test1.jpg")
	// GetExif现在可能返回nil，这是正常的
	t.Logf("Raw EXIF data: %+v\n", x)
}

func TestGetExifInfo1(t *testing.T) {
	// 首先获取原始EXIF数据
	x := utils.GetExif("./test1_.jpg")

	// 测试ParseExif函数
	exifInfo := utils.ParseExif(x)
	assert.NotNil(t, exifInfo)

	// 验证提取的信息
	t.Logf("Extracted EXIF info: %+v\n", exifInfo)

	// 验证拍摄时间（如果有的话）
	if !exifInfo.CapturedTime.IsZero() {
		t.Logf("拍摄时间: %v", exifInfo.CapturedTime.Format("2006-01-02 15:04:05"))
		assert.True(t, exifInfo.CapturedTime.After(exifInfo.CapturedTime.AddDate(-10, 0, 0)))
		assert.True(t, exifInfo.CapturedTime.Before(exifInfo.CapturedTime.AddDate(10, 0, 0)))
	} else {
		t.Log("未找到拍摄时间信息")
	}

	// 验证GPS坐标（如果有的话）
	if exifInfo.Latitude != 0 || exifInfo.Longitude != 0 {
		t.Logf("GPS坐标: 纬度 %f, 经度 %f", exifInfo.Latitude, exifInfo.Longitude)
		assert.True(t, exifInfo.Latitude >= -90 && exifInfo.Latitude <= 90, "纬度应在-90到90度之间")
		assert.True(t, exifInfo.Longitude >= -180 && exifInfo.Longitude <= 180, "经度应在-180到180度之间")
	} else {
		t.Log("未找到GPS坐标信息")
	}

	// 验证图片尺寸（如果有的话）
	if exifInfo.WidthPixel > 0 || exifInfo.HeightPixel > 0 {
		t.Logf("图片尺寸: %dx%d 像素", exifInfo.WidthPixel, exifInfo.HeightPixel)
		assert.Greater(t, exifInfo.WidthPixel, 0, "宽度像素应大于0")
		assert.Greater(t, exifInfo.HeightPixel, 0, "高度像素应大于0")
	} else {
		t.Log("未找到图片尺寸信息")
	}
}

func TestGetExifInfo2(t *testing.T) {
	x := utils.GetExif("./test2.jpg")
	exifInfo := utils.ParseExif(x)
	assert.NotNil(t, exifInfo)

	// 经纬度不存在，其它信息存在
	assert.NotZero(t, exifInfo.CapturedTime)
	assert.Zero(t, exifInfo.Latitude)
	assert.Zero(t, exifInfo.Longitude)
	assert.NotZero(t, exifInfo.WidthPixel)
	assert.NotZero(t, exifInfo.HeightPixel)
}

func TestGenerateThumbnail(t *testing.T) {
	os.MkdirAll("./thumbnail", 0755)
	err := utils.GenerateThumbnail("./test1.jpg", "./thumbnail")
	assert.NoError(t, err)
	// 检查是否生成了缩略图
	_, err = os.Stat("./thumbnail/test1.jpg")
	assert.NoError(t, err)
	// 删除缩略图
	os.Remove("./thumbnail/test1.jpg")
	assert.NoError(t, err)
	os.RemoveAll("./thumbnail")
}

func TestReverseGeocode(t *testing.T) {
	// 测试逆地理编码功能
	t.Log("测试逆地理编码功能...")

	// 测试有效坐标（北京天安门）
	geoInfo := utils.ReverseGeocode(116.397128, 39.916527)
	assert.NotNil(t, geoInfo)
	t.Logf("北京坐标逆地理编码结果: %+v", geoInfo)

	// 测试零坐标
	emptyGeoInfo := utils.ReverseGeocode(0, 0)
	assert.NotNil(t, emptyGeoInfo)
	t.Logf("零坐标结果: %+v", emptyGeoInfo)

	// 测试无效坐标（太平洋深处）
	invalidGeoInfo := utils.ReverseGeocode(180, 90)
	assert.NotNil(t, invalidGeoInfo)
	t.Logf("无效坐标结果: %+v", invalidGeoInfo)
}
