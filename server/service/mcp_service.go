package service

import (
	"context"
	"encoding/json"
	"fmt"

	"server/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"gorm.io/gorm"
)

type MCPService struct {
	db           *gorm.DB
	openaiClient openai.Client
	mcpServerURL string
}

func NewMCPService(db *gorm.DB) *MCPService {
	// 初始化 OpenAI 客户端（使用 DeepSeek API）
	apiKey := "sk-884a6ecd756644cdb4ead4240b0b73cd"
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.deepseek.com/v1"),
	)
	return &MCPService{
		db:           db,
		openaiClient: client,
		mcpServerURL: "http://localhost:7777",
	}
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

// 获取用户的会话
func (s *MCPService) GetUserChat(userID int) ([]model.Chat, error) {
	chats := []model.Chat{}
	err := s.db.Where("user_id = ?", userID).Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}

// 保存用户消息
func (s *MCPService) SaveUserMessage(chatID int, content string) error {
	message := &model.Message{
		Content: content,
		IsReply: false,
		ChatID:  chatID,
	}
	return s.db.Create(&message).Error
}

// 保存AI回复消息
func (s *MCPService) SaveAIMessage(chatID int, content string) error {
	message := &model.Message{
		Content: content,
		IsReply: true,
		ChatID:  chatID,
	}
	return s.db.Create(&message).Error
}

// 获取会话历史消息
func (s *MCPService) GetChatMessages(chatID int) ([]model.GetChatMessagesResponse, error) {
	var messages []model.Message
	err := s.db.Where("chat_id = ?", chatID).Order("created_time DESC").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	var responses []model.GetChatMessagesResponse
	for _, message := range messages {
		responses = append(responses, model.GetChatMessagesResponse{
			Content:     message.Content,
			IsReply:     message.IsReply,
			CreatedTime: message.CreatedTime,
		})
	}
	return responses, nil
}

// 调用 MCP 工具搜索图片
func (s *MCPService) invokeMCPTool(toolAlias string, arguments map[string]interface{}) (string, error) {
	// 1. 创建 MCP 客户端实例
	client := mcp.NewClient(&mcp.Implementation{Name: "go-agent", Version: "1.0"}, nil)

	// 2. 创建 HTTP Transport，连接到 /search 路径
	// StreamableClientTransport 需要完整的 URL，包括路径
	transport := &mcp.StreamableClientTransport{Endpoint: s.mcpServerURL + "/search"}

	// 3. 连接到服务器，建立会话
	ctx := context.Background()
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return "", fmt.Errorf("failed to connect to MCP server for tool %s: %w", toolAlias, err)
	}
	defer session.Close() // 每次调用都是一个独立的会话，确保关闭

	// 4. 执行调用并处理结果
	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolAlias,
		Arguments: arguments,
	})
	if err != nil {
		return "", fmt.Errorf("CallTool failed: %w", err)
	}

	if res.IsError {
		if len(res.Content) > 0 {
			if textContent, ok := res.Content[0].(*mcp.TextContent); ok {
				return "", fmt.Errorf("tool execution failed: %s", textContent.Text)
			}
		}
		return "", fmt.Errorf("tool execution failed")
	}

	// 5. 提取结果文本
	if len(res.Content) == 0 {
		return "", fmt.Errorf("tool returned empty content")
	}

	if textContent, ok := res.Content[0].(*mcp.TextContent); ok {
		return textContent.Text, nil
	}

	return "", fmt.Errorf("unexpected content type in tool response")
}

// Chat 处理对话，返回 AI 的最终回答和图片 ID
func (s *MCPService) Chat(ctx context.Context, chatID int, userMessage string) (*model.ChatResult, error) {
	// 1. 保存用户消息
	if err := s.SaveUserMessage(chatID, userMessage); err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %v", err)
	}

	// 2. 构建 OpenAI 消息列表（仅使用当前问题，不加载历史对话）
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful assistant with access to image search tools. You must call tools by using the tool_calls response format. Don't make assumptions about what values to plug into functions. Ask for clarification if a user request is ambiguous."),
		openai.UserMessage(userMessage),
	}

	// 3. 定义可用工具
	tools := []openai.ChatCompletionToolParam{
		{
			Function: openai.FunctionDefinitionParam{
				Name:        "search_images",
				Description: openai.String("Search images in qdrant by text query. Returns image IDs, similarity scores, and count."),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"text": map[string]string{
							"type":        "string",
							"description": "The text query to search for images",
						},
					},
					"required": []string{"text"},
				},
			},
		},
	}

	// 4. 进行多轮对话（最多5轮），避免AI回复过长
	maxTurns := 5
	var fullResponse string
	var imageIDs []uint64 // 收集所有工具调用返回的图片 ID

	for i := 0; i < maxTurns; i++ {
		// 调用非流式 API
		completion, err := s.openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    "deepseek-chat",
			Messages: messages,
			Tools:    tools,
		})
		if err != nil {
			return nil, fmt.Errorf("调用 AI API 失败: %v", err)
		}

		if len(completion.Choices) == 0 {
			return nil, fmt.Errorf("AI API 返回空响应")
		}

		choice := completion.Choices[0]
		message := choice.Message

		// 处理内容
		if message.Content != "" {
			fullResponse = message.Content
		}

		// 检查是否有工具调用
		hasToolCalls := len(message.ToolCalls) > 0
		if hasToolCalls {
			// 转换工具调用类型
			toolCalls := make([]openai.ChatCompletionMessageToolCallParam, len(message.ToolCalls))
			for i, tc := range message.ToolCalls {
				toolCalls[i] = openai.ChatCompletionMessageToolCallParam{
					ID:   tc.ID,
					Type: tc.Type,
					Function: openai.ChatCompletionMessageToolCallFunctionParam{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
				}
			}

			// 创建包含工具调用的助手消息
			assistantMsg := openai.ChatCompletionAssistantMessageParam{
				ToolCalls: toolCalls,
			}
			currentMessage := openai.ChatCompletionMessageParamUnion{
				OfAssistant: &assistantMsg,
			}
			messages = append(messages, currentMessage)

			// 处理工具调用
			for _, toolCall := range message.ToolCalls {
				functionName := toolCall.Function.Name
				var arguments map[string]interface{}
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
					return nil, fmt.Errorf("解析工具参数失败: %v", err)
				}

				// 调用 MCP 工具
				toolResult, err := s.invokeMCPTool(functionName, arguments)
				if err != nil {
					toolResult = fmt.Sprintf("Error executing tool: %v", err)
				} else {
					// 解析工具结果，提取图片 ID
					var toolResultData struct {
						ImageIDs []uint64  `json:"image_ids"`
						Scores   []float32 `json:"scores"`
						Count    int       `json:"count"`
					}
					if err := json.Unmarshal([]byte(toolResult), &toolResultData); err == nil {
						// 成功解析，添加图片 ID 到列表
						imageIDs = append(imageIDs, toolResultData.ImageIDs...)
					}
				}

				// 添加工具结果到消息列表
				messages = append(messages, openai.ToolMessage(toolResult, toolCall.ID))
			}

			// 重置 fullResponse，准备接收新的回复
			fullResponse = ""
			continue
		}

		// 没有工具调用，对话结束
		break
	}

	// 5. 保存 AI 回复
	if fullResponse != "" {
		if err := s.SaveAIMessage(chatID, fullResponse); err != nil {
			return nil, fmt.Errorf("保存 AI 消息失败: %v", err)
		}
	}

	// 6. 返回结果
	return &model.ChatResult{
		Reply:    fullResponse,
		ImageIDs: imageIDs,
	}, nil
}
