package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"server/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// 向量化任务队列的 Redis key
	VectorQueueKey = "imagehub:vector:queue"
	// 队列消费超时时间（秒）
	QueueTimeout = 5
)

// VectorTask 向量化任务结构
type VectorTask struct {
	ImageID  int    `json:"image_id"`
	FilePath string `json:"file_path"`
}

// EnqueueVectorTask 将向量化任务推入 Redis 队列
func EnqueueVectorTask(imageID int, filePath string) error {
	redisClient := global.GetRedisClient()
	if redisClient == nil {
		return fmt.Errorf("redis 客户端未初始化")
	}

	task := VectorTask{
		ImageID:  imageID,
		FilePath: filePath,
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("序列化任务失败: %v", err)
	}

	ctx := context.Background()
	err = redisClient.LPush(ctx, VectorQueueKey, taskJSON).Err()
	if err != nil {
		return fmt.Errorf("推入队列失败: %v", err)
	}

	logger := global.GetLogger()
	if logger != nil {
		logger.Info("向量化任务已入队",
			zap.Int("image_id", imageID),
			zap.String("file_path", filePath),
		)
	}

	return nil
}

// StartVectorWorker 启动向量化任务消费 worker
func StartVectorWorker(imageService *ImageService) {
	logger := global.GetLogger()
	if logger == nil {
		fmt.Println("警告: logger 未初始化，worker 将使用 fmt 输出日志")
	}

	redisClient := global.GetRedisClient()
	if redisClient == nil {
		if logger != nil {
			logger.Error("Redis 客户端未初始化，worker 无法启动")
		} else {
			fmt.Println("错误: Redis 客户端未初始化，worker 无法启动")
		}
		return
	}

	if logger != nil {
		logger.Info("向量化任务 worker 已启动")
	} else {
		fmt.Println("向量化任务 worker 已启动")
	}

	ctx := context.Background()

	for {
		// 使用 BRPop 阻塞式从队列右侧弹出任务（FIFO）
		result, err := redisClient.BRPop(ctx, time.Duration(QueueTimeout)*time.Second, VectorQueueKey).Result()
		if err != nil {
			if err == redis.Nil {
				// 超时，继续循环
				continue
			}
			if logger != nil {
				logger.Error("从队列获取任务失败", zap.Error(err))
			} else {
				fmt.Printf("从队列获取任务失败: %v\n", err)
			}
			// 发生错误时等待一段时间再重试
			time.Sleep(time.Second)
			continue
		}

		if len(result) < 2 {
			continue
		}

		taskJSON := result[1]
		var task VectorTask
		err = json.Unmarshal([]byte(taskJSON), &task)
		if err != nil {
			if logger != nil {
				logger.Error("反序列化任务失败",
					zap.Error(err),
					zap.String("task_json", taskJSON),
				)
			} else {
				fmt.Printf("反序列化任务失败: %v, task: %s\n", err, taskJSON)
			}
			continue
		}

		// 处理任务
		if logger != nil {
			logger.Info("开始处理向量化任务",
				zap.Int("image_id", task.ImageID),
				zap.String("file_path", task.FilePath),
			)
		} else {
			fmt.Printf("开始处理向量化任务: image_id=%d, file_path=%s\n", task.ImageID, task.FilePath)
		}

		err = imageService.storeImageVector(task.ImageID, task.FilePath)
		if err != nil {
			if logger != nil {
				logger.Error("处理向量化任务失败",
					zap.Int("image_id", task.ImageID),
					zap.String("file_path", task.FilePath),
					zap.Error(err),
				)
			} else {
				fmt.Printf("处理向量化任务失败: image_id=%d, error=%v\n", task.ImageID, err)
			}
		} else {
			if logger != nil {
				logger.Info("向量化任务处理成功",
					zap.Int("image_id", task.ImageID),
				)
			} else {
				fmt.Printf("向量化任务处理成功: image_id=%d\n", task.ImageID)
			}
		}
	}
}
