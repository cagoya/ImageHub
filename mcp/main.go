// agent/main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// toolRegistry 映射对 LLM 友好的工具别名到其 HTTP 服务地址
var toolRegistry = map[string]string{
	"add": "http://localhost:8080/math",
}

// invokeMCPTool 是 Agent 的核心函数，负责直接与 MCP 服务通信（支持 HTTP）
func invokeMCPTool(toolAlias string, arguments map[string]interface{}) (string, error) {
	httpAddr, ok := toolRegistry[toolAlias]
	if !ok {
		return "", fmt.Errorf("unknown tool alias: %s", toolAlias)
	}

	// 1. 创建 MCP 客户端实例
	client := mcp.NewClient(&mcp.Implementation{Name: "go-agent", Version: "1.0"}, nil)

	// 2. 创建 HTTP Transport
	transport := &mcp.StreamableClientTransport{Endpoint: httpAddr}

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
		return "", fmt.Errorf("tool execution failed: %s", res.Content[0].(*mcp.TextContent).Text)
	}
	resultText := res.Content[0].(*mcp.TextContent).Text

	return resultText, nil
}

func main() {
	apiKey := "sk-884a6ecd756644cdb4ead4240b0b73cd"
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY not set.")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.deepseek.com/v1"),
	)

	// 定义可用的工具（仅 HTTP 服务）
	tools := []openai.ChatCompletionToolParam{
		{
			Function: openai.FunctionDefinitionParam{
				Name:        "add",
				Description: openai.String("Add two integers."),
				Parameters: openai.FunctionParameters{
					"type": "object", "properties": map[string]interface{}{"A": map[string]string{"type": "integer"}, "B": map[string]string{"type": "integer"}}, "required": []string{"A", "B"},
				},
			},
		},
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful assistant with access to local tools. You must call tools by using the tool_calls response format. Don't make assumptions about what values to plug into functions. Ask for clarification if a user request is ambiguous."),
		openai.UserMessage("Hi, can you add 5 and 7?"),
	}

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		log.Println("--- Sending request to DeepSeek ---")

		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{Model: "deepseek-chat", Messages: messages, Tools: tools})
		if err != nil {
			log.Fatalf("ChatCompletion error: %v\n", err)
		}
		if len(resp.Choices) == 0 {
			log.Fatal("No choices returned from API")
		}

		msg := resp.Choices[0].Message
		messages = append(messages, msg.ToParam())

		if msg.ToolCalls != nil {
			for _, toolCall := range msg.ToolCalls {
				functionName := toolCall.Function.Name
				var arguments map[string]interface{}
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
					log.Fatalf("Failed to unmarshal function arguments: %v", err)
				}

				log.Printf("--- LLM wants to call tool: %s with args: %v ---\n", functionName, arguments)

				// 直接调用我们的 Go 函数，该函数内建了 MCP 客户端逻辑
				toolResult, err := invokeMCPTool(functionName, arguments)
				if err != nil {
					log.Printf("Tool call failed: %v\n", err)
					toolResult = fmt.Sprintf("Error executing tool: %v", err)
				}

				log.Printf("--- Tool result: ---\n%s\n---------------------\n", toolResult)

				messages = append(messages, openai.ToolMessage(toolResult, toolCall.ID))
			}
			continue
		}

		log.Println("--- Final response from LLM ---")
		log.Println(msg.Content)
		return
	}
	log.Println("Reached max conversation turns.")
}
