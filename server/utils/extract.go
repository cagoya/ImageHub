package utils

import (
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Exif struct {
	CapturedTime time.Time `json:"captured_time"` // 拍摄时间
	Latitude     float64   `json:"latitude"`      // 纬度
	Longitude    float64   `json:"longitude"`     // 经度
	WidthPixel   int       `json:"width_pixel"`   // 图片宽度像素
	HeightPixel  int       `json:"height_pixel"`  // 图片高度像素
}

// GetExif 获取图片的EXIF信息
func GetExif(file *os.File) (*exif.Exif, error) {
	x, err := exif.Decode(file)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func GetExifInfo(x *exif.Exif) (*Exif, error) {
	result := &Exif{}

	// 检查EXIF数据是否为nil
	if x == nil {
		return result, nil
	}

	// 提取拍摄时间 (DateTimeOriginal)
	if tm, err := x.DateTime(); err == nil {
		result.CapturedTime = tm
	}

	// 提取GPS经纬度信息
	if lat, lon, err := x.LatLong(); err == nil {
		result.Latitude = lat
		result.Longitude = lon
	}

	// 提取图片尺寸 (PixelXDimension, PixelYDimension)
	if width, err := x.Get(exif.PixelXDimension); err == nil {
		if w, err := width.Int(0); err == nil {
			result.WidthPixel = w
		}
	}
	if height, err := x.Get(exif.PixelYDimension); err == nil {
		if h, err := height.Int(0); err == nil {
			result.HeightPixel = h
		}
	}

	return result, nil
}
