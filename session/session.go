package session

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	sessionRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/zjutjh/WeJH-SDK/redis"
	"strconv"
)

type InfoConfig struct {
	Name        string
	SecretKey   string
	RedisConfig *redis.InfoConfig
}

// Init 使用 Redis 初始化会话管理
func Init(config *InfoConfig, r *gin.Engine) error {
	redisConfig := config.RedisConfig

	store, err := sessionRedis.NewStoreWithDB(10, "tcp",
		redisConfig.Host+":"+redisConfig.Port, redisConfig.Password,
		strconv.Itoa(redisConfig.DB),
		[]byte(config.SecretKey),
	)
	if err != nil {
		return fmt.Errorf("session init failed: %w", err)
	}
	r.Use(sessions.Sessions(config.Name, store))
	return nil
}