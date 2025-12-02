package api

import (
	"server/global"
	"server/model"
	"server/service"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MCPHandler MCP 处理接口
type MCPHandler struct {
	mcpService *service.MCPService
}

// NewMCPHandler 创建新的 MCPHandler
func NewMCPHandler() *MCPHandler {
	db := global.GetDB()
	mcpService := service.NewMCPService(db)
	return &MCPHandler{mcpService: mcpService}
}

// Chat 处理对话
func (h *MCPHandler) Chat(c *gin.Context) {
	// 从 JWT 中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userIDInt := userID.(int)
	var chatID int

	// 如果是新对话，创建新会话
	if req.IsNewChat {
		var err error
		chatID, err = h.mcpService.CreateChat(userIDInt)
		if err != nil {
			utils.ServerError(c, "创建会话失败: "+err.Error())
			return
		}
	} else {
		// 使用现有会话
		if req.ChatID <= 0 {
			utils.BadRequest(c, "chat_id 必须大于0")
			return
		}
		chatID = req.ChatID
	}

	// 调用服务处理对话
	result, err := h.mcpService.Chat(c.Request.Context(), chatID, req.Content)
	if err != nil {
		utils.ServerError(c, "处理消息失败: "+err.Error())
		return
	}

	utils.Success(c, model.ChatResponse{
		ChatID:   chatID,
		Reply:    result.Reply,
		ImageIDs: result.ImageIDs,
	})
}

// GetChats 获取用户的会话列表
func (h *MCPHandler) GetChats(c *gin.Context) {
	// 从 JWT 中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	userIDInt := userID.(int)
	chats, err := h.mcpService.GetUserChat(userIDInt)
	if err != nil {
		utils.ServerError(c, "获取会话列表失败: "+err.Error())
		return
	}

	chatInfos := make([]model.ChatInfo, len(chats))
	for i, chat := range chats {
		chatInfos[i] = model.ChatInfo{
			ID:        chat.ID,
			CreatedAt: chat.CreatedTime.Format("2006-01-02 15:04:05"),
		}
	}

	utils.Success(c, model.GetChatsResponse{Chats: chatInfos})
}

// GetChatMessages 获取会话消息
func (h *MCPHandler) GetChatMessages(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Query("chatID"))
	if err != nil {
		utils.BadRequest(c, "chatID 必须大于0")
		return
	}
	chatMessages, err := h.mcpService.GetChatMessages(chatID)
	if err != nil {
		utils.ServerError(c, "获取会话消息失败: "+err.Error())
		return
	}
	utils.Success(c, chatMessages)
}
