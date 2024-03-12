package sdk

import Redis "github.com/go-redis/redis/v8"

type RedisInfoConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
}

func GetRedisClient(info RedisInfoConfig) *Redis.Client {
	redisClient := Redis.NewClient(&Redis.Options{
		Addr:     info.Host + ":" + info.Port,
		Password: info.Password,
		DB:       info.DB,
	})
	return redisClient
}

func DefaultRedisConfig() RedisInfoConfig {
	return RedisInfoConfig{
		Host:     "localhost",
		Port:     "6379",
		DB:       0,
		Password: "",
	}
}
