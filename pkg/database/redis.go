package database

import (
	"DBP/config"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHandler struct {
	Client *redis.Client
}

var (
	redisInstance *RedisHandler
	redisOnce     sync.Once
)

// Redis 인스턴스를 싱글톤 패턴으로 관리
func GetRedisInstance() *RedisHandler {
	redisOnce.Do(func() {
		config.LoadEnv()
		client := redis.NewClient(&redis.Options{
			Addr:         os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			DB:           0,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		})
		redisInstance = &RedisHandler{Client: client}
	})
	return redisInstance
}
