package initialize

import (
	"os"
	"server/global"

	"github.com/go-redis/redis/v8"

	"go.uber.org/zap"
)

// InitRedis 初始化并返回一个 Redis 客户端，支持集群或单节点配置
func InitRedis() *redis.Client {
	redisCfg := global.Configs.Redis

	// 使用单节点配置创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	// Ping Redis 服务器以检查连接是否正常
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		global.Log.Error("redis connect error:", zap.Error(err))
		os.Exit(1)
	}
	global.Log.Info("redis connect success")
	return client
}
