package service

import (
	"context"
	"server/global"
	"strconv"
)

// NewArticleView 先将浏览记录存储在redis中，后续有定时任务将redis中的数据写入es
func (articleService *ArticleService) NewArticleView() CountDB {
	return CountDB{
		Index: "article_views",
	}
}

type CountDB struct {
	Index string
}

// Set 在原有基础上加一
func (c CountDB) Set(ctx context.Context, id string) error {
	//忽略非法 ID
	if id == "" || len(id) < 10 {
		return nil
	}
	num, _ := global.Redis.HGet(ctx, c.Index, id).Int()
	num++
	err := global.Redis.HSet(ctx, c.Index, id, num).Err()
	return err
}

// GetInfo 取出数据
func (c CountDB) GetInfo(ctx context.Context) map[string]int {
	var Info = map[string]int{}
	maps := global.Redis.HGetAll(ctx, c.Index).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val)
		Info[id] = num
	}
	return Info
}

// Del 删除指定文章的浏览量缓存
func (c CountDB) Del(ctx context.Context, id string) error {
	// 忽略非法 ID
	if id == "" || len(id) < 10 {
		return nil
	}
	return global.Redis.HDel(ctx, c.Index, id).Err()
}

// Clear 清除数据
func (c CountDB) Clear(ctx context.Context) {
	global.Redis.Del(ctx, c.Index)
}
