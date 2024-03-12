package sdk

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	sessionRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
)

type driver string

const (
	memory driver = "memory"
	redis  driver = "redis"
)

var defaultName = "wejh-sdk-session"

type SessionInfoConfig struct {
	Driver      string
	Name        string
	SecretKey   string
	RedisConfig RedisInfoConfig
}

func SessionInit(r *gin.Engine, config SessionInfoConfig) {
	switch config.Driver {
	case string(redis):
		setSessionRedis(r, config)
		break
	case string(memory):
		setSessionMemory(r, config.Name)
		break
	default:
		log.Fatal("ConfigError")
	}
}

func DefaultSessionConfig() SessionInfoConfig {
	return SessionInfoConfig{
		Driver:      "Memory",
		Name:        defaultName,
		SecretKey:   "secret",
		RedisConfig: DefaultRedisConfig(),
	}
}

func setSessionMemory(r *gin.Engine, name string) {
	store := memstore.NewStore()
	r.Use(sessions.Sessions(name, store))
}

func setSessionRedis(r *gin.Engine, config SessionInfoConfig) {
	Info := config.RedisConfig
	store, _ := sessionRedis.NewStore(10, "tcp", Info.Host+":"+Info.Port, Info.Password, []byte(config.SecretKey))
	r.Use(sessions.Sessions(config.Name, store))
}
