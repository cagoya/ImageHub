package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const ALI_APIKEY = "sk-9ef56cbc5d944b65b1bb612fd7864877"
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

/**
 * 向量化一张图片
 *
 * 图片编码要求：
 * - key为 "image"
 * - value可以是：
 *   1. 公开可访问的URL（如：https://example.com/image.jpg）
 *   2. Base64编码的Data URI，格式为 data:image/{format};base64,{data}
 *      - {format} 是图片格式（如 jpeg, png, gif, webp）
 *      - {data} 是Base64编码字符串
 *
 * @param imageURL 图片的URL地址或base64 data URL（格式：data:image/{format};base64,{data}）
 * @return []float64 1024维的向量
 * @return error 错误信息
 */
func VectorizeImage(imageURL string) ([]float64, error) {
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ALI_APIKEY))
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ALI_APIKEY))
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
