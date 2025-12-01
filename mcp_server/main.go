package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HiParams 和 SayHi 函数
type HiParams struct {
	Name string `json:"name" jsonschema:"the person to greet"`
}

type HiResult struct {
	Message string `json:"message"`
}

func SayHi(_ context.Context, _ *mcp.CallToolRequest, in HiParams) (*mcp.CallToolResult, HiResult, error) {
	resultText := fmt.Sprintf("Hi %s, this response is from the HTTP server!", in.Name)
	return nil, HiResult{Message: resultText}, nil
}

// AddParams 和 Add 工具的实现
type AddParams struct {
	A int `json:"A" jsonschema:"first number"`
	B int `json:"B" jsonschema:"second number"`
}

type AddResult struct {
	Sum int `json:"sum"`
}

func Add(_ context.Context, _ *mcp.CallToolRequest, in AddParams) (*mcp.CallToolResult, AddResult, error) {
	result := in.A + in.B
	return nil, AddResult{Sum: result}, nil
}

func main() {
	// 1. 创建 Greeter 服务实例
	greeterServer := mcp.NewServer(&mcp.Implementation{Name: "greeter-service", Version: "1.0"}, nil)
	mcp.AddTool(greeterServer, &mcp.Tool{
		Name:        "greet",
		Description: "Say hi to someone",
	}, SayHi)

	// 2. 创建 Math 服务实例
	mathServer := mcp.NewServer(&mcp.Implementation{Name: "math-service", Version: "1.0"}, nil)
	mcp.AddTool(mathServer, &mcp.Tool{
		Name:        "add",
		Description: "Add two integers",
	}, Add)

	// 3. 创建 StreamableHTTPHandler
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		log.Printf("Routing request for URL: %s\n", request.URL.Path)
		switch request.URL.Path {
		case "/greeter":
			return greeterServer
		case "/math":
			return mathServer
		default:
			return nil // 返回 nil 将导致 404 Not Found
		}
	}, nil)

	// 4. 启动标准的 Go HTTP 服务器
	addr := ":8080"
	log.Printf("Multi-service MCP server listening at http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
