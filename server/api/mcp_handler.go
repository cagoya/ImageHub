package api

import (
	"fmt"
	"server/global"
	"server/model"
	"server/service"
	"server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

// Chat 处理对话，使用 WebSocket
func (h *MCPHandler) Chat(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		utils.BadRequest(c, "用户ID不能为空")
		return
	}
	// 升级 HTTP 连接到 WebSocket
	conn, err := model.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.ServerError(c, "WebSocket 升级失败: "+err.Error())
		return
	}
	defer conn.Close()

	// 将新连接添加到连接管理器
	connManager := model.GetConnManager()
	connManager.AddConn(userID, conn)

	// 监听连接断开事件，并从管理器中移除
	defer func() {
		connManager.RemoveConn(userID)
	}()

	// 设置读取超时和关闭处理
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// 启动心跳检测
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	// 当前会话ID，初始为0表示未创建
	var currentChatID int

	// 循环读取客户端发送的消息
	for {
		var req model.ChatRequest
		// 读取JSON消息
		err := conn.ReadJSON(&req)
		if err != nil {
			// 检查是否是连接关闭错误
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 连接异常关闭
				fmt.Printf("WebSocket连接异常关闭: %v\n", err)
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				// 正常关闭
				fmt.Printf("WebSocket连接正常关闭\n")
			} else {
				// 其他错误（如超时、格式错误等）
				fmt.Printf("WebSocket读取错误: %v\n", err)
			}
			break
		}

		// 重置读取超时
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// 验证消息内容
		if req.Content == "" {
			conn.WriteJSON(gin.H{"error": "消息内容不能为空"})
			continue
		}

		var chatID int

		// 如果是新对话，创建新会话
		if req.IsNewChat {
			var err error
			chatID, err = h.mcpService.CreateChat(userID)
			if err != nil {
				conn.WriteJSON(gin.H{"error": "创建会话失败: " + err.Error()})
				continue
			}
			currentChatID = chatID
		} else {
			// 使用现有会话
			if req.ChatID > 0 {
				chatID = req.ChatID
				currentChatID = chatID
			} else if currentChatID > 0 {
				// 如果没有指定 chatID 但已有当前会话，使用当前会话
				chatID = currentChatID
			} else {
				conn.WriteJSON(gin.H{"error": "chat_id 必须大于0或先创建新会话"})
				continue
			}
		}

		// 调用服务处理对话
		result, err := h.mcpService.Chat(c.Request.Context(), chatID, req.Content)
		if err != nil {
			conn.WriteJSON(gin.H{"error": "处理消息失败: " + err.Error()})
			continue
		}

		// 发送响应给客户端
		response := model.ChatResponse{
			ChatID: chatID,
			Reply:  result.Reply,
			Images: result.Images,
		}
		if err := conn.WriteJSON(response); err != nil {
			fmt.Printf("发送消息失败: %v\n", err)
			break
		}
	}
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
