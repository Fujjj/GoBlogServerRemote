package task

import (
	"context"
	"server/global"
	"server/model/elasticsearch"
	"server/service"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"go.uber.org/zap"
)

// UpdateArticleViewsSyncTask 将 Redis 中的文章浏览量（增量），同步到 Elasticsearch
func UpdateArticleViewsSyncTask(ctx context.Context) error {
	// 获取redis中的缓存数据
	articleView := service.ServiceGroupApp.ArticleService.NewArticleView()

	viewsInfo := articleView.GetInfo(ctx)
	var lastErr error
	for id, num := range viewsInfo {
		// 无变化就跳过
		if num == 0 {
			continue
		}

		//先判断文档是否存在
		exists, err := global.ESClient.Exists(
			elasticsearch.ArticleIndex(),
			id,
		).Do(context.TODO())
		if err != nil {
			global.Log.Error("check exists failed", zap.String("id", id), zap.Error(err))
			lastErr = err
			continue
		}
		// 文档不存在则  直接跳过，并清理 Redis
		if !exists {
			global.Log.Warn("article not found in es, skip update", zap.String("id", id))
			articleView.Del(ctx, id) // 清理单个 key
			continue
		}

		// 更新数据 之前的数据+缓存中的数据
		source := "ctx._source.views += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err = global.ESClient.Update(elasticsearch.ArticleIndex(), id).Script(&script).Do(context.TODO())
		if err != nil {
			global.Log.Error("update failed", zap.String("id", id), zap.Error(err))
			lastErr = err
			continue //继续处理下一个
		}
	}

	// 清除redis中的数据
	articleView.Clear(ctx)
	return lastErr
}
