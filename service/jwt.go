package service

import (
	"server/model/database"
	"server/utils"

	"server/global"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

// JwtService 提供与JWT相关的服务
type JwtService struct {
}

// SetRedisJWT 将JWT存储到Redis中
func (jwtService *JwtService) SetRedisJWT(c *gin.Context, jwt string, uuid uuid.UUID) error {
	// 解析配置中的JWT过期时间
	ept, err := utils.ParseDuration(global.Configs.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		return err
	}
	// 设置JWT在Redis中的过期时间
	return global.Redis.Set(c, uuid.String(), jwt, ept).Err()
}

// GetRedisJWT 从Redis中获取JWT
func (jwtService *JwtService) GetRedisJWT(c *gin.Context, uuid uuid.UUID) (string, error) {
	// 从Redis获取指定uuid对应的JWT
	return global.Redis.Get(c, uuid.String()).Result()
}

// JoinInBlacklist 将JWT添加到黑名单
func (jwtService *JwtService) JoinInBlacklist(jwtList database.JwtBlacklist) error {
	// 将JWT记录插入到数据库中的黑名单表，传入指针 &jwtList 允许 GORM 直接操作内存中的原始对象（回填，如修改gorm.Model中的数据）。
	if err := global.DB.Create(&jwtList).Error; err != nil {
		return err
	}
	// 将JWT添加到内存中的黑名单缓存，本来是传入键值对的，这里只传入键，值用空结构体代替
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return nil
}

// IsInBlacklist 检查JWT是否在黑名单中
//
//	func (jwtService *JwtService) IsInBlacklist(jwt string) bool {
//		// 从黑名单缓存中检查JWT是否存在
//		_, ok := global.BlackCache.Get(jwt)
//		return ok
//	}
func (jwtService *JwtService) IsInBlacklist(jwt string) bool {
	// 1. 优先从黑名单缓存中检查JWT是否存在
	if _, ok := global.BlackCache.Get(jwt); ok {
		return true
	}

	// 2. 缓存未命中，查询数据库进行兜底
	var blacklist database.JwtBlacklist
	// 使用 First 查找第一条匹配记录，如果找不到会返回 gorm.ErrRecordNotFound
	err := global.DB.Where("jwt = ?", jwt).First(&blacklist).Error

	// 3. 如果数据库中存在该黑名单记录
	if err == nil {
		// 将结果回填到缓存，避免下次重复查库
		// 注意：这里使用 SetDefault，过期时间由全局配置决定
		global.BlackCache.SetDefault(jwt, struct{}{})
		return true
	}

	return false
}

// LoadAll 从数据库加载所有的JWT黑名单并加入缓存
func LoadAll() {
	var data []string
	// 从数据库中获取所有的黑名单JWT
	if err := global.DB.Model(&database.JwtBlacklist{}).Pluck("jwt", &data).Error; err != nil {
		// 如果获取失败，记录错误日志
		global.Log.Error("Failed to load JWT blacklist from the database", zap.Error(err))
		return
	}
	// 将所有JWT添加到BlackCache缓存中
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	}
}
