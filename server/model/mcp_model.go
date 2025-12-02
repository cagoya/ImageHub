package model

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket 升级器
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 开发环境允许所有资源
		return true
	},
}

// ConnManager WebSocket 连接管理器
// 使用 sync.Map 来保证并发安全
type ConnManager struct {
	connections sync.Map
}

// 单例模式
var connManager = &ConnManager{}

// GetConnManager 获取连接管理器单例
func GetConnManager() *ConnManager {
	return connManager
}

// AddConn 向管理器添加连接
func (m *ConnManager) AddConn(userID int, conn *websocket.Conn) {
	m.connections.Store(userID, conn)
}

// GetConn 从管理器获取连接
func (m *ConnManager) GetConn(userID int) (*websocket.Conn, bool) {
	conn, ok := m.connections.Load(userID)
	if !ok {
		return nil, false
	}
	return conn.(*websocket.Conn), true
}

// RemoveConn 从管理器移除连接
func (m *ConnManager) RemoveConn(userID int) {
	m.connections.Delete(userID)
}

// IsConnValid 检查连接是否有效
func (m *ConnManager) IsConnValid(userID int) bool {
	conn, ok := m.GetConn(userID)
	if !ok {
		return false
	}
	// 尝试发送 ping 消息来检测连接是否有效
	err := conn.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		// 连接无效，从管理器中移除
		m.RemoveConn(userID)
		return false
	}
	return true
}

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
	ChatID int                  `json:"chat_id"` // 会话ID（新对话时返回新创建的ID）
	Reply  string               `json:"reply"`   // AI回复
	Images []ImageQueryResponse `json:"images"`  // 图片信息列表
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
	Reply  string               `json:"reply"`
	Images []ImageQueryResponse `json:"images"`
}

// GetChatMessagesResponse 获取会话消息响应
type GetChatMessagesResponse struct {
	Content     string    `json:"content"`
	IsReply     bool      `json:"is_reply"`
	CreatedTime time.Time `json:"created_time"`
}
