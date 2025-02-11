package redisHelper

import (
	"github.com/redis/go-redis/v9"
)

// InfoConfig Redis 配置
type InfoConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
}

// Init 初始化 Redis 客户端
func Init(info *InfoConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     info.Host + ":" + info.Port,
		Password: info.Password,
		DB:       info.DB,
	})
	return client
}
