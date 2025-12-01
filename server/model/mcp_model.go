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
