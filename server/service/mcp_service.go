package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"server/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type MCPService struct {
	db           *gorm.DB
	openaiClient openai.Client
	mcpServerURL string
}

func getDeepseekAPIKey() string {
	return viper.GetString("deepseek.apiKey")
}

func getMCPServerURL() string {
	url := viper.GetString("mcp.serverURL")
	if url == "" {
		// 默认值，用于向后兼容
		return "http://localhost:7777"
	}
	return url
}

func NewMCPService(db *gorm.DB) *MCPService {
	// 初始化 OpenAI 客户端（使用 DeepSeek API）
	apiKey := getDeepseekAPIKey()
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.deepseek.com/v1"),
	)
	return &MCPService{
		db:           db,
		openaiClient: client,
		mcpServerURL: getMCPServerURL(),
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
func (s *MCPService) SaveAIMessage(chatID int, content string, imagePaths []string) error {
	// 获取用户ID，用于构建完整的图片路径
	var chat model.Chat
	if err := s.db.Where("id = ?", chatID).First(&chat).Error; err != nil {
		return fmt.Errorf("获取会话信息失败: %v", err)
	}
	userID := chat.UserID

	// 复制图片到public/ai目录，并更新路径
	var copiedImagePaths []string
	if len(imagePaths) > 0 {
		// 确保public/ai目录存在
		aiDir := "./public/ai"
		if err := os.MkdirAll(aiDir, 0755); err != nil {
			return fmt.Errorf("创建AI图片目录失败: %v", err)
		}

		for _, imagePath := range imagePaths {
			// 构建源文件路径：public/{userID}/origin/{imagePath}
			srcPath := path.Join("./public", fmt.Sprintf("%d", userID), "origin", imagePath)

			// 检查源文件是否存在
			if _, err := os.Stat(srcPath); os.IsNotExist(err) {
				// 如果源文件不存在，跳过该图片，但记录警告
				fmt.Printf("警告: 源图片文件不存在，跳过: %s\n", srcPath)
				continue
			}

			// 目标文件路径：public/ai/{imagePath}
			dstPath := path.Join(aiDir, imagePath)

			// 复制文件
			if err := s.copyImageFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("复制图片到AI目录失败: %v", err)
			}

			// 更新路径为ai/{imagePath}格式
			copiedImagePaths = append(copiedImagePaths, path.Join("ai", imagePath))
		}
	}

	var imagesJSON string
	if len(copiedImagePaths) > 0 {
		imagesBytes, err := json.Marshal(copiedImagePaths)
		if err != nil {
			return fmt.Errorf("序列化图片路径失败: %v", err)
		}
		imagesJSON = string(imagesBytes)
	}
	message := &model.Message{
		Content: content,
		IsReply: true,
		ChatID:  chatID,
		Images:  imagesJSON,
	}
	return s.db.Create(&message).Error
}

// copyImageFile 复制图片文件
func (s *MCPService) copyImageFile(srcPath, dstPath string) error {
	// 读取源文件
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("读取源文件失败: %v", err)
	}

	// 写入目标文件
	if err := os.WriteFile(dstPath, data, 0644); err != nil {
		return fmt.Errorf("写入目标文件失败: %v", err)
	}

	return nil
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
		var imagePaths []string
		if message.Images != "" {
			if err := json.Unmarshal([]byte(message.Images), &imagePaths); err != nil {
				// 如果解析失败，记录错误但继续处理，不返回空图片列表
				imagePaths = []string{}
			}
		}
		responses = append(responses, model.GetChatMessagesResponse{
			Content:     message.Content,
			IsReply:     message.IsReply,
			CreatedTime: message.CreatedTime,
			Images:      imagePaths,
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
	// 0. 获取用户ID
	var chat model.Chat
	if err := s.db.Where("id = ?", chatID).First(&chat).Error; err != nil {
		return nil, fmt.Errorf("获取会话信息失败: %v", err)
	}
	userID := chat.UserID

	// 1. 加载历史消息（按时间正序排列，在保存当前消息之前加载，避免重复）
	var historyMessages []model.Message
	if err := s.db.Where("chat_id = ?", chatID).Order("created_time ASC").Find(&historyMessages).Error; err != nil {
		return nil, fmt.Errorf("加载历史消息失败: %v", err)
	}

	// 2. 保存用户消息
	if err := s.SaveUserMessage(chatID, userMessage); err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %v", err)
	}

	// 3. 构建 OpenAI 消息列表
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful assistant with access to image search tools. You must call tools by using the tool_calls response format. Don't make assumptions about what values to plug into functions. Ask for clarification if a user request is ambiguous. IMPORTANT: When the search_images tool returns an empty result (count is 0 or image_ids is empty), it means no matching images were found. In this case, you should inform the user that no matching images were found and do NOT continue searching with different queries."),
	}

	// 4. 将历史消息添加到消息列表
	for _, msg := range historyMessages {
		if msg.IsReply {
			// AI 回复
			messages = append(messages, openai.AssistantMessage(msg.Content))
		} else {
			// 用户消息
			messages = append(messages, openai.UserMessage(msg.Content))
		}
	}

	// 5. 添加当前用户消息
	messages = append(messages, openai.UserMessage(userMessage))

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

	// 4. 进行多轮对话（目前是第一轮查询需要调用哪些函数，第二轮捎带上mcp server的结果），避免AI回复过长
	maxTurns := 2
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
		if len(message.ToolCalls) > 0 {
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

				// 如果是搜索图片工具，添加 user_id 参数
				if functionName == "search_images" {
					arguments["user_id"] = userID
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

	// 5. 根据 imageIDs 查询完整的图片信息
	images, err := s.getImagesByIDs(imageIDs)
	if err != nil {
		return nil, fmt.Errorf("获取图片信息失败: %v", err)
	}

	// 6. 提取图片路径列表并复制图片到public/ai目录
	var imagePaths []string
	for _, img := range images {
		imagePaths = append(imagePaths, img.ImagePath)
	}

	// 7. 保存 AI 回复（会复制图片到public/ai目录并更新路径）
	if fullResponse != "" {
		if err := s.SaveAIMessage(chatID, fullResponse, imagePaths); err != nil {
			return nil, fmt.Errorf("保存 AI 消息失败: %v", err)
		}
	}

	// 8. 更新返回的图片路径为ai/{filename}格式，与保存到数据库的路径一致
	var updatedImages []model.ImageQueryResponse
	for _, img := range images {
		updatedImages = append(updatedImages, model.ImageQueryResponse{
			ImageID:    img.ImageID,
			ImageName:  img.ImageName,
			ImagePath:  path.Join("ai", img.ImagePath), // 更新为ai/{filename}格式
			AITags:     img.AITags,
			ManualTags: img.ManualTags,
		})
	}

	// 9. 返回结果
	return &model.ChatResult{
		Reply:  fullResponse,
		Images: updatedImages,
	}, nil
}

// getImagesByIDs 根据图片ID列表获取完整的图片信息（包括标签）
func (s *MCPService) getImagesByIDs(imageIDs []uint64) ([]model.ImageQueryResponse, error) {
	if len(imageIDs) == 0 {
		return []model.ImageQueryResponse{}, nil
	}

	// 将 uint64 转换为 int（因为数据库中的 ID 是 int）
	var ids []int
	for _, id := range imageIDs {
		ids = append(ids, int(id))
	}

	// 查询图片信息
	var images []model.Image
	err := s.db.Where("id IN ?", ids).Find(&images).Error
	if err != nil {
		return nil, err
	}

	// 转换为 ImageQueryResponse
	var responses []model.ImageQueryResponse
	for _, image := range images {
		// 获取标签
		aiTags, manualTags, err := s.getImageTags(image.ID)
		if err != nil {
			// 如果获取标签失败，继续处理其他图片，但记录错误
			aiTags = []string{}
			manualTags = []string{}
		}

		responses = append(responses, model.ImageQueryResponse{
			ImageID:    image.ID,
			ImageName:  image.Name,
			ImagePath:  image.Path,
			AITags:     aiTags,
			ManualTags: manualTags,
		})
	}

	return responses, nil
}

// getImageTags 获取图片的标签，区分AI生成和手动设置
func (s *MCPService) getImageTags(imageID int) ([]string, []string, error) {
	var aiTags []string
	var manualTags []string
	var tags []model.Tag
	err := s.db.Where("image_id = ?", imageID).Find(&tags).Error
	if err != nil {
		return nil, nil, err
	}
	for _, tag := range tags {
		if tag.FromAI {
			aiTags = append(aiTags, tag.TagName)
		} else {
			manualTags = append(manualTags, tag.TagName)
		}
	}
	return aiTags, manualTags, nil
}
