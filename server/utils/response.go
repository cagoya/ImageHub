package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 全局统一响应结构
// Code HTTP 状态码
// Message 响应消息
// Data 响应数据
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "成功",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = "请求参数错误"
	}
	Fail(c, http.StatusBadRequest, message)
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "认证令牌无效"
	}
	Fail(c, http.StatusUnauthorized, message)
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	Fail(c, http.StatusNotFound, message)
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "无权限访问"
	}
	Fail(c, http.StatusForbidden, message)
}

// ServerError 500错误响应
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	Fail(c, http.StatusInternalServerError, message)
}
