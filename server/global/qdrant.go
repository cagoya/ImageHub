package global

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
	"github.com/spf13/viper"
)

var qdrantClient *qdrant.Client

func InitQdrant() error {
	var err error
	host := viper.GetString("qdrant.host")
	port := viper.GetInt("qdrant.port")
	qdrantClient, err = qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		return fmt.Errorf("初始化qdrant失败: %v", err)
	}
	exist, err := qdrantClient.CollectionExists(context.Background(), "images")
	if err != nil {
		return fmt.Errorf("检查集合是否存在失败: %v", err)
	}
	if !exist {
		err = qdrantClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
			CollectionName: "images",
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     1024,
				Distance: qdrant.Distance_Cosine,
			}),
		})
		if err != nil {
			return fmt.Errorf("创建集合失败: %v", err)
		}
	}
	return nil
}

func GetQdrantClient() *qdrant.Client {
	return qdrantClient
}
