package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

type Exif struct {
	CapturedTime time.Time `json:"captured_time"` // 拍摄时间
	Latitude     float64   `json:"latitude"`      // 纬度
	Longitude    float64   `json:"longitude"`     // 经度
	WidthPixel   int       `json:"width_pixel"`   // 图片宽度像素
	HeightPixel  int       `json:"height_pixel"`  // 图片高度像素
}

// GetExif 获取图片的EXIF信息，失败时返回nil但不报错
func GetExif(filePath string) *exif.Exif {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()
	x, err := exif.Decode(file)
	// 某些图片可能没有EXIF信息，忽略所有错误
	if err != nil {
		return nil
	}
	return x
}

// ParseExif 解析EXIF数据，无论如何都不会返回错误
func ParseExif(x *exif.Exif) *Exif {
	result := &Exif{}

	// 检查EXIF数据是否为nil
	if x == nil {
		return result
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

	return result
}

func GenerateThumbnail(inputPath string) error {
	// 打开原始图片
	src, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开图片失败: %v", err)
	}

	// 将图片缩放至宽度为200像素，高度按比例自动计算
	dst := imaging.Resize(src, 200, 0, imaging.Lanczos)

	// 保存缩略图，名称为原始名称前加thumbnail_
	thumbnailPath := path.Join(filepath.Dir(inputPath), "thumbnail_"+filepath.Base(inputPath))
	err = imaging.Save(dst, thumbnailPath)
	if err != nil {
		return fmt.Errorf("保存缩略图失败: %v", err)
	}

	return nil
}
