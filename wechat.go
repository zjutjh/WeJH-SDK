package sdk

import (
	"context"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

type WechatConfig struct {
	AppId     string
	AppSecret string
}

func WeChatInitInMemory(config WechatConfig) *miniprogram.MiniProgram {
	wc := wechat.NewWechat()
	var wcCache cache.Cache

	wcCache = cache.NewMemory()

	cfg := &miniConfig.Config{
		AppID:     config.AppId,
		AppSecret: config.AppSecret,
		Cache:     wcCache,
	}
	return wc.GetMiniProgram(cfg)
}

func WeChatInitInRedis(config WechatConfig, redisConfig RedisInfoConfig) *miniprogram.MiniProgram {
	wc := wechat.NewWechat()
	var wcCache cache.Cache

	wcCache = setWechatRedis(wcCache, redisConfig)

	cfg := &miniConfig.Config{
		AppID:     config.AppId,
		AppSecret: config.AppSecret,
		Cache:     wcCache,
	}
	return wc.GetMiniProgram(cfg)
}

func setWechatRedis(wcCache cache.Cache, config RedisInfoConfig) cache.Cache {
	redisOpts := &cache.RedisOpts{
		Host:        config.Host + ":" + config.Port,
		Database:    config.DB,
		MaxActive:   10,
		MaxIdle:     10,
		IdleTimeout: 60,
	}
	wcCache = cache.NewRedis(context.Background(), redisOpts)
	return wcCache
}
