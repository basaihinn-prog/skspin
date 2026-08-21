package gcdb

import (
	"context"
	"encoding/json"
	"game/comm"
	"game/duck/lazy"
	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"log/slog"
	"time"
)

var redisClient *redis.Client

func InitRedis() {
	redisURL := lo.Must(lazy.RouteFile.Get("redis"))
	opt, err := redis.ParseURL(redisURL)
	redisClient = redis.NewClient(opt)
	// 测试 Redis 连接
	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		slog.Error("Redis 连接失败:", err)
		return
	}
	slog.Info("Redis 连接成功:", pong)
}

func StoreNextMult(key string, v *comm.NextMult) error {
	jsonStore, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return redisClient.SetEX(context.Background(), key, jsonStore, time.Hour*24*15).Err()
}
