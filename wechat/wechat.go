package wechat

import (
	"context"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/zjutjh/WeJH-SDK/redis"
)

type InfoConfig struct {
	AppId       string
	AppSecret   string
	RedisConfig *redis.InfoConfig
}

// Init 初始化微信小程序
func Init(config *InfoConfig) *miniprogram.MiniProgram {
	redisConfig := config.RedisConfig

	wc := wechat.NewWechat()
	wcCache := cache.NewRedis(context.Background(), &cache.RedisOpts{
		Host:        redisConfig.Host + ":" + redisConfig.Port,
		Database:    redisConfig.DB,
		MaxActive:   10,
		MaxIdle:     10,
		IdleTimeout: 60,
	})

	cfg := &miniConfig.Config{
		AppID:     config.AppId,
		AppSecret: config.AppSecret,
		Cache:     wcCache,
	}

	return wc.GetMiniProgram(cfg)
}
