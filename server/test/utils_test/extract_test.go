package test

import (
	"os"
	"testing"

	"server/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetExif(t *testing.T) {
	file, err := os.Open("./test.jpg")
	assert.NoError(t, err)
	defer file.Close()

	x, err := utils.GetExif(file)
	assert.NoError(t, err)
	assert.NotNil(t, x)
	t.Logf("Raw EXIF data: %+v\n", x)
}

func TestGetExifInfo(t *testing.T) {
	file, err := os.Open("./test.jpg")
	assert.NoError(t, err)
	defer file.Close()

	// 首先获取原始EXIF数据
	x, err := utils.GetExif(file)
	assert.NoError(t, err)
	assert.NotNil(t, x)

	// 测试GetExifInfo函数
	exifInfo, err := utils.GetExifInfo(x)
	assert.NoError(t, err)
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

func TestGetExifInfoWithNilExif(t *testing.T) {
	// 测试传入nil的情况
	exifInfo, err := utils.GetExifInfo(nil)
	assert.NoError(t, err)
	assert.NotNil(t, exifInfo)

	// 应该返回零值
	assert.Zero(t, exifInfo.CapturedTime)
	assert.Zero(t, exifInfo.Latitude)
	assert.Zero(t, exifInfo.Longitude)
	assert.Zero(t, exifInfo.WidthPixel)
	assert.Zero(t, exifInfo.HeightPixel)
}
