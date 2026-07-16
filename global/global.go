package global

import (
	"server/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Configs    *config.Config //用于存储项目的配置信息
	Log        *zap.Logger
	DB         *gorm.DB
	ESClient   *elasticsearch.TypedClient
	Redis      *redis.Client
	BlackCache local_cache.Cache
)
