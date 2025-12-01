package service

import (
	"server/model"

	"gorm.io/gorm"
)

type MCPService struct {
	db *gorm.DB
}

func NewMCPService(db *gorm.DB) *MCPService {
	return &MCPService{db: db}
}

// 创建一次会话
func (s *MCPService) CreateChat(userID int) (int, error) {
	chat := &model.Chat{
		UserID: userID,
	}
	err := s.db.Create(&chat).Error
	if err != nil {
		return 0, err
	}
	return chat.ID, nil
}
