package model

import (
	"time"
)

// Chat 一次会话
type Chat struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`                                                       // 自增主键
	UserID      int       `json:"user_id" gorm:"not null;index"`                                                            // 用户ID
	User        User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 外键约束
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`                                                       // 创建时间
	UpdatedTime time.Time `json:"updated_time" gorm:"autoUpdateTime"`                                                       // 更新时间
}

// Message 会话中的一条消息
type Message struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`                                                       // 自增主键
	Content     string    `json:"content" gorm:"not null;type:text"`                                                        // 消息内容
	IsReply     bool      `json:"is_reply" gorm:"default:false"`                                                            // 是否为AI的回复，用于区分用户和AI的回复
	ChatID      int       `json:"chat_id" gorm:"not null;index"`                                                            // 聊天ID
	Chat        Chat      `json:"chat" gorm:"foreignKey:ChatID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 外键约束
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`                                                       // 创建时间
	UpdatedTime time.Time `json:"updated_time" gorm:"autoUpdateTime"`                                                       // 更新时间
}

// ChatRequest 对话请求
type ChatRequest struct {
	IsNewChat bool   `json:"is_new_chat"`                // 是否是新对话
	ChatID    int    `json:"chat_id"`                    // 会话ID（新对话时可为0，已有对话时必填）
	Content   string `json:"content" binding:"required"` // 消息内容
}

// ChatResponse 对话响应
type ChatResponse struct {
	ChatID   int      `json:"chat_id"`   // 会话ID（新对话时返回新创建的ID）
	Reply    string   `json:"reply"`     // AI回复
	ImageIDs []uint64 `json:"image_ids"` // 图片ID列表
}

// GetChatsResponse 获取会话列表响应
type GetChatsResponse struct {
	Chats []ChatInfo `json:"chats"`
}

// ChatInfo 会话信息
type ChatInfo struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
}

// ChatResult 对话结果
type ChatResult struct {
	Reply    string   `json:"reply"`
	ImageIDs []uint64 `json:"image_ids"`
}

// GetChatMessagesResponse 获取会话消息响应
type GetChatMessagesResponse struct {
	Content     string    `json:"content"`
	IsReply     bool      `json:"is_reply"`
	CreatedTime time.Time `json:"created_time"`
}
