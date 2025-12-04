package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"server/global"

	"github.com/spf13/viper"
)

// getBaiduApiKey 获取百度AI API密钥，在函数调用时动态读取配置
func getBaiduApiKey() string {
	return viper.GetString("baidu.apiKey")
}

// getBaiduSecretKey 获取百度AI Secret Key，在函数调用时动态读取配置
func getBaiduSecretKey() string {
	return viper.GetString("baidu.secretKey")
}

var SECRET_KEY = viper.GetString("baidu.secretKey")

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func GetAccessToken() string {
	// 先试着从Redis中获取AccessToken
	redisClient := global.GetRedisClient()
	accessToken, err := redisClient.Get(context.Background(), "baidu_access_token").Result()
	if err != redis.Nil {
		return accessToken
	}

	// 如果Redis中没有AccessToken，则重新获取
	tokenUrl := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", getBaiduApiKey(), getBaiduSecretKey())

	resp, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		fmt.Println("获取 Access Token 失败:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取 Access Token 响应体失败:", err)
		return ""
	}

	accessTokenObj := map[string]any{}
	if err := json.Unmarshal(body, &accessTokenObj); err != nil {
		fmt.Println("解析 Access Token 失败:", err)
		return ""
	}

	if token, ok := accessTokenObj["access_token"].(string); ok {
		// 将AccessToken保存到Redis中
		err = redisClient.Set(context.Background(), "baidu_access_token", token, time.Hour*48).Err()
		if err != nil {
			fmt.Println("保存 Access Token 到 Redis 失败:", err)
			return ""
		}
		return token
	}

	fmt.Println("Access Token 字段不存在或类型错误")
	return ""
}

// 百度AI识别结果结构体
type BaiduClassificationResult struct {
	Result []struct {
		Score   float64 `json:"score"`
		Root    string  `json:"root"`
		Keyword string  `json:"keyword"`
	} `json:"result"`
	ResultNum int   `json:"result_num"`
	LogID     int64 `json:"log_id"`
}

/**
 * 调用百度AI高级通用图像识别接口
 * @param imageBase64 图像文件的 Base64 编码字符串（不含文件头）
 * @return []string 识别结果中的关键词数组
 * @return error 错误信息
 */
func ImageClassification(imageBase64 string) ([]string, error) {
	accessToken := GetAccessToken()
	if accessToken == "" {
		return nil, fmt.Errorf("无法获取 Access Token")
	}

	// 1. 拼接 URL
	apiUrl := "https://aip.baidubce.com/rest/2.0/image-classify/v2/advanced_general?access_token=" + accessToken

	// 2. 对 Base64 编码进行 URL 编码
	encodedImage := url.QueryEscape(imageBase64)
	payload := strings.NewReader("image=" + encodedImage)

	// 3. 发送请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, payload)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer res.Body.Close()

	// 4. 读取响应
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 5. 检查 HTTP 状态码
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求返回非成功状态码: %d, 响应: %s", res.StatusCode, string(body))
	}

	// 6. 解析 JSON 响应
	var result BaiduClassificationResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应 JSON 失败: %w", err)
	}

	// 7. 提取关键词数组
	keywords := make([]string, 0, len(result.Result))
	for _, item := range result.Result {
		if item.Keyword != "" {
			keywords = append(keywords, item.Keyword)
		}
	}

	return keywords, nil
}

/**
 * 获取文件 base64 编码
 * @param path 文件路径
 * @return string base64 编码信息，不带文件头
 */
func GetFileContentAsBase64(path string) string {
	// 替换 ioutil.ReadFile(path) 为 os.ReadFile(path)
	srcByte, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(srcByte)
}
