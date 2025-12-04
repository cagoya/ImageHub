package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

// getGaodeApiKey 获取高德地图API密钥，在函数调用时动态读取配置
func getGaodeApiKey() string {
	return viper.GetString("gaode.apiKey")
}

// 完整的响应结构体，这个只在本文件内使用
type apiResponse struct {
	Status    string `json:"status"`
	Regeocode struct {
		AddressComponent struct {
			Province string      `json:"province"`
			City     interface{} `json:"city"`
			District string      `json:"district"`
		} `json:"addressComponent"`
		FormattedAddress string `json:"formatted_address"`
	} `json:"regeocode"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
}

// 地理信息，只取需要的部分
type GeoInfo struct {
	Province         string `json:"province"`
	City             string `json:"city"`
	District         string `json:"district"`
	FormattedAddress string `json:"formatted_address"`
}

// ReverseGeocode 逆地理编码API，这里选用高德地图的，仅支持中国境内
// 失败时返回空的GeoInfo，不返回错误
func ReverseGeocode(longitude float64, latitude float64) *GeoInfo {
	// 检查坐标是否有效
	if longitude == 0 && latitude == 0 {
		return &GeoInfo{}
	}

	baseURL := "https://restapi.amap.com/v3/geocode/regeo"

	// 构建查询参数
	params := url.Values{}
	params.Add("key", getGaodeApiKey())
	params.Add("location", fmt.Sprintf("%f,%f", longitude, latitude))
	params.Add("extensions", "base")

	// 发送HTTP GET请求
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	httpReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return &GeoInfo{}
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return &GeoInfo{}
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &GeoInfo{}
	}

	var result apiResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &GeoInfo{}
	}

	city, ok := result.Regeocode.AddressComponent.City.(string)
	if !ok {
		city = ""
	}

	return &GeoInfo{
		Province:         result.Regeocode.AddressComponent.Province,
		City:             city,
		District:         result.Regeocode.AddressComponent.District,
		FormattedAddress: result.Regeocode.FormattedAddress,
	}
}
