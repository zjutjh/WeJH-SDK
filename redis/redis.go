package redis

import (
	Redis "github.com/go-redis/redis/v8"
)

type InfoConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
}

func Init(info *InfoConfig) *Redis.Client {
	client := Redis.NewClient(&Redis.Options{
		Addr:     info.Host + ":" + info.Port,
		Password: info.Password,
		DB:       info.DB,
	})
	return client
}
