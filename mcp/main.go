package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"mcp_server/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/qdrant/go-client/qdrant"
)

// QdrantSearchParams 查询参数
type QdrantSearchParams struct {
	Text   string `json:"text" jsonschema:"the text to search for"`
	UserID int    `json:"user_id" jsonschema:"the user ID to filter by"`
}

// QdrantSearchResult 查询结果
type QdrantSearchResult struct {
	ImageIDs []uint64  `json:"image_ids"`
	Scores   []float32 `json:"scores"`
	Count    int       `json:"count"`
}

// QdrantSearch 在 qdrant 中搜索
func QdrantSearch(ctx context.Context, _ *mcp.CallToolRequest, in QdrantSearchParams) (*mcp.CallToolResult, QdrantSearchResult, error) {
	// 1. 将文本转换为向量
	vectors, err := utils.VectorizeText(in.Text)
	if err != nil {
		return nil, QdrantSearchResult{}, fmt.Errorf("向量化文本失败: %v", err)
	}

	// 2. 将 []float64 转换为 []float32（Qdrant 需要 float32）
	vector32 := make([]float32, len(vectors))
	for i, v := range vectors {
		vector32[i] = float32(v)
	}

	// 3. 连接 Qdrant
	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		return nil, QdrantSearchResult{}, fmt.Errorf("连接 Qdrant 失败: %v", err)
	}

	// 4. 在 qdrant 中搜索
	limit := uint64(3) // 返回前3个最相似的结果

	// 构建过滤器，根据 user_id 过滤
	filter := &qdrant.Filter{
		Must: []*qdrant.Condition{
			qdrant.NewMatchInt("user_id", int64(in.UserID)),
		},
	}

	searchResult, err := qdrantClient.Query(ctx, &qdrant.QueryPoints{
		CollectionName: "images",
		Query:          qdrant.NewQueryDense(vector32),
		Filter:         filter,
		Limit:          &limit,
		WithPayload: &qdrant.WithPayloadSelector{
			SelectorOptions: &qdrant.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
	})
	if err != nil {
		return nil, QdrantSearchResult{}, fmt.Errorf("搜索 Qdrant 失败: %v", err)
	}

	// 5. 提取结果
	var imageIDs []uint64
	var scores []float32
	for _, result := range searchResult {
		if result.Id != nil {
			imageIDs = append(imageIDs, result.Id.GetNum())
		}
		scores = append(scores, result.Score)
	}

	// 6. 返回结果
	if len(imageIDs) == 0 {
		log.Printf("No images found for user ID: %d", in.UserID)
		return nil, QdrantSearchResult{
			ImageIDs: []uint64{},
			Scores:   []float32{},
			Count:    0,
		}, nil
	} else {
		log.Printf("Images found for user ID: %d", in.UserID)
		log.Printf("Image IDs: %v", imageIDs)
		log.Printf("Scores: %v", scores)
		log.Printf("Count: %d", len(imageIDs))
		return nil, QdrantSearchResult{
			ImageIDs: imageIDs,
			Scores:   scores,
			Count:    len(imageIDs),
		}, nil
	}
}

func main() {
	// 创建 Qdrant 搜索服务实例
	searchServer := mcp.NewServer(&mcp.Implementation{Name: "qdrant-search-service", Version: "1.0"}, nil)
	mcp.AddTool(searchServer, &mcp.Tool{
		Name:        "search_images",
		Description: "Search images in qdrant by text query",
	}, QdrantSearch)

	// 创建 StreamableHTTPHandler
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		log.Printf("Routing request for URL: %s\n", request.URL.Path)
		switch request.URL.Path {
		case "/search":
			return searchServer
		default:
			return nil // 返回 nil 将导致 404 Not Found
		}
	}, nil)

	// 启动标准的 Go HTTP 服务器
	addr := ":7777"
	log.Printf("Qdrant search MCP server listening at http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
