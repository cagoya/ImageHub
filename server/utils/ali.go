package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/spf13/viper"
)

// getAliApiKey 获取阿里云API密钥，在函数调用时动态读取配置
func getAliApiKey() string {
	return viper.GetString("ali.apikey")
}

const API_URL = "https://dashscope.aliyuncs.com/api/v1/services/embeddings/multimodal-embedding/multimodal-embedding"
const MODEL = "qwen2.5-vl-embedding"

// DashScope API 响应结构体
type DashScopeResponse struct {
	Output struct {
		Embeddings []struct {
			Index     int       `json:"index"`
			Embedding []float64 `json:"embedding"`
			Type      string    `json:"type"`
		} `json:"embeddings"`
	} `json:"output"`
	Usage struct {
		InputTokens int `json:"input_tokens"`
		ImageTokens int `json:"image_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
}

// DashScope API 请求结构体
type DashScopeRequest struct {
	Model string `json:"model"`
	Input struct {
		Contents []map[string]interface{} `json:"contents"`
	} `json:"input"`
}

/**
 * 将本地图片文件转换为 base64 data URL
 *
 * 图片编码要求：
 * - Base64格式为 data:image/{format};base64,{data}
 * - 其中 {format} 是图片格式（如 jpeg, png, gif, webp）
 * - {data} 是Base64编码字符串
 *
 * @param filePath 本地文件路径
 * @return string base64 data URL，格式：data:image/{format};base64,{data}
 * @return error 错误信息
 */
func ImageFileToDataURL(filePath string) (string, error) {
	// 读取文件
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 获取文件扩展名以确定图片格式
	ext := strings.ToLower(filepath.Ext(filePath))
	format := "jpeg" // 默认格式
	switch ext {
	case ".jpg", ".jpeg":
		format = "jpeg"
	case ".png":
		format = "png"
	case ".gif":
		format = "gif"
	case ".webp":
		format = "webp"
	}

	// 转换为 base64，格式：data:image/{format};base64,{data}
	base64String := base64.StdEncoding.EncodeToString(fileBytes)
	return fmt.Sprintf("data:image/%s;base64,%s", format, base64String), nil
}

/**
 * 压缩 data URL 格式的图片
 * @param dataURL base64 data URL，格式：data:image/{format};base64,{data}
 * @param maxWidth 最大宽度，0表示不限制宽度，只按质量压缩
 * @param quality JPEG质量，1-100，仅对JPEG有效
 * @return string 压缩后的 data URL
 * @return error 错误信息
 */
func compressDataURLImage(dataURL string, maxWidth int, quality int) (string, error) {
	// 解析 data URL 格式：data:image/{format};base64,{data}
	if !strings.HasPrefix(dataURL, "data:image/") {
		return "", fmt.Errorf("无效的 data URL 格式")
	}

	parts := strings.SplitN(dataURL, ";base64,", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("无效的 data URL 格式：缺少 ;base64, 分隔符")
	}

	// 提取格式：data:image/{format}
	formatPart := strings.TrimPrefix(parts[0], "data:image/")
	format := strings.ToLower(formatPart)

	// 解码 base64
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("解码 base64 失败: %w", err)
	}

	// 从字节数组解码图片
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("解码图片失败: %w", err)
	}

	// 压缩图片
	var compressedImg image.Image
	if maxWidth > 0 && img.Bounds().Dx() > maxWidth {
		// 按宽度缩放
		compressedImg = imaging.Resize(img, maxWidth, 0, imaging.Lanczos)
	} else {
		compressedImg = img
	}

	// 编码为字节数组
	var compressedData []byte
	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		err = imaging.Encode(&buf, compressedImg, imaging.JPEG, imaging.JPEGQuality(quality))
	case "png":
		err = imaging.Encode(&buf, compressedImg, imaging.PNG)
	case "gif":
		err = imaging.Encode(&buf, compressedImg, imaging.GIF)
	case "webp":
		// imaging 库不支持 webp，转换为 JPEG
		err = imaging.Encode(&buf, compressedImg, imaging.JPEG, imaging.JPEGQuality(quality))
		format = "jpeg"
	default:
		// 默认使用 JPEG
		err = imaging.Encode(&buf, compressedImg, imaging.JPEG, imaging.JPEGQuality(quality))
		format = "jpeg"
	}

	if err != nil {
		return "", fmt.Errorf("编码图片失败: %w", err)
	}

	compressedData = buf.Bytes()

	// 重新编码为 base64 data URL
	compressedBase64 := base64.StdEncoding.EncodeToString(compressedData)
	return fmt.Sprintf("data:image/%s;base64,%s", format, compressedBase64), nil
}

/**
 * 向量化一张图片（从本地文件）
 * @param filePath 本地文件路径
 * @return []float64 1024维的向量
 * @return error 错误信息
 */
func VectorizeImageFromFile(filePath string) ([]float64, error) {
	dataURL, err := ImageFileToDataURL(filePath)
	if err != nil {
		return nil, err
	}
	return VectorizeImage(dataURL)
}

// 向量化一张图片
func VectorizeImage(imageURL string) ([]float64, error) {
	return vectorizeImageWithRetry(imageURL, false)
}

// vectorizeImageWithRetry 内部函数，支持重试标志
func vectorizeImageWithRetry(imageURL string, isRetry bool) ([]float64, error) {
	// 验证输入：如果以 data: 开头，验证格式是否正确
	if strings.HasPrefix(imageURL, "data:") {
		if !strings.HasPrefix(imageURL, "data:image/") {
			return nil, fmt.Errorf("Base64 Data URI格式错误：必须以 data:image/ 开头")
		}
		if !strings.Contains(imageURL, ";base64,") {
			return nil, fmt.Errorf("Base64 Data URI格式错误：必须包含 ;base64, 分隔符")
		}
	}
	// 构建请求体
	request := DashScopeRequest{
		Model: MODEL,
		Input: struct {
			Contents []map[string]interface{} `json:"contents"`
		}{
			Contents: []map[string]interface{}{
				{
					"image": imageURL,
				},
			},
		},
	}

	// 序列化请求体
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", getAliApiKey()))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		// 如果是 400 错误且是 data URL 格式且还未重试过，尝试压缩后重试
		if resp.StatusCode == http.StatusBadRequest && strings.HasPrefix(imageURL, "data:") && !isRetry {
			// 压缩图片：最大宽度 1920px，JPEG 质量 85
			compressedURL, compressErr := compressDataURLImage(imageURL, 1920, 85)
			if compressErr != nil {
				return nil, fmt.Errorf("API请求返回400错误，压缩图片失败: %w, 响应: %s", compressErr, string(body))
			}
			// 使用压缩后的图片重试（标记为已重试）
			return vectorizeImageWithRetry(compressedURL, true)
		}
		return nil, fmt.Errorf("API请求返回非成功状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result DashScopeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %w", err)
	}

	// 检查是否有返回的向量
	if len(result.Output.Embeddings) == 0 {
		return nil, fmt.Errorf("API响应中没有向量数据")
	}

	// 返回第一个向量（应该只有一个）
	return result.Output.Embeddings[0].Embedding, nil
}

/**
 * 向量化一段文本
 * @param text 文本内容
 * @return []float64 1024维的向量
 * @return error 错误信息
 */
func VectorizeText(text string) ([]float64, error) {
	// 构建请求体
	request := DashScopeRequest{
		Model: MODEL,
		Input: struct {
			Contents []map[string]interface{} `json:"contents"`
		}{
			Contents: []map[string]interface{}{
				{
					"text": text,
				},
			},
		},
	}

	// 序列化请求体
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", getAliApiKey()))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求返回非成功状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result DashScopeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %w", err)
	}

	// 检查是否有返回的向量
	if len(result.Output.Embeddings) == 0 {
		return nil, fmt.Errorf("API响应中没有向量数据")
	}

	// 返回第一个向量（应该只有一个）
	return result.Output.Embeddings[0].Embedding, nil
}
