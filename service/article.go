package service

import (
	"context"
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/elasticsearch"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleService struct {
}

func (articleService *ArticleService) ArticleLike(req request.ArticleLike) error {
	//启动事务，事务是为了保证一组数据库操作要么全成功，要么全失败
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var al database.ArticleLike
		var num int

		//判断用户是否收藏
		if errors.Is(tx.Where("user_id=? AND article_id = ?", req.UserID, req.ArticleID).First(&al).Error, gorm.ErrRecordNotFound) {
			//若没有收藏，则创建收藏
			if err := tx.Create(&database.ArticleLike{UserID: req.UserID, ArticleID: req.ArticleID}).Error; err != nil {
				return err
			}
			num = 1
			//否则取消收藏
		} else {
			//
			if err := tx.Delete(&al).Error; err != nil {
				return err
			}
			num = -1
		}

		// 更新文章收藏数 构建脚本语言，实现原子操作
		source := "ctx._source.likes += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).Script(&script).Do(context.TODO())

		//如果 ES 文档不存在（文章已被删除），记录日志但不回滚 MySQL 事务
		if err != nil {
			if strings.Contains(err.Error(), "document_missing_exception") {
				global.Log.Warn("article not found in ES when updating likes, article may have been deleted",
					zap.String("article_id", req.ArticleID),
					zap.Int("num", num),
				)
				return nil // MySQL 操作已成功，不返回错误
			}
			return err // 其他错误正常返回
		}
		return nil
	})
}
func (articleService *ArticleService) ArticleIsLike(req request.ArticleLike) (bool, error) {
	return !errors.Is(global.DB.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&database.ArticleLike{}).Error, gorm.ErrRecordNotFound), nil
}
func (articleService *ArticleService) ArticleLikesList(info request.ArticleLikesList) (interface{}, int64, error) {
	//因为文章表存储在es中，所以先通过MySQL分页查询文章收藏表（article_like），然后再把每条点赞记录关联的文章数据补齐
	db := global.DB.Where("user_id=?", info.UserID)
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	//此时 l里全是当前用户的点赞记录
	l, total, err := utils.MySQLPagination(&database.ArticleLike{}, option)
	if err != nil {
		return nil, 0, err
	}

	//按照es的数据结构转换
	var list []struct {
		Id_     string                `json:"_id"`
		Source_ elasticsearch.Article `json:"_source"`
	}

	//遍历点赞记录，补文章数据
	for _, articleLike := range l {
		//去 Elasticsearch 查文章详情
		article, err := articleService.Get(articleLike.ArticleID)
		if err != nil {
			return nil, 0, err
		}

		//拿到文章后，裁剪不必要字段，优化性能
		article.UpdatedAt = ""
		article.Keyword = ""
		article.Content = ""

		//将数据添加进list
		list = append(list, struct {
			Id_     string                `json:"_id"`
			Source_ elasticsearch.Article `json:"_source"`
		}{
			Id_:     articleLike.ArticleID,
			Source_: article,
		})
	}
	return list, total, nil

}
func (articleService *ArticleService) ArticleInfoById(ctx context.Context, id string) (elasticsearch.Article, error) {
	// 异步更新浏览量
	go func() {
		ctx := context.Background()
		articleView := articleService.NewArticleView()
		_ = articleView.Set(ctx, id)
	}()
	//根据文章id获取文章
	return articleService.Get(id)
}

func (articleService *ArticleService) ArticleSearch(ctx context.Context, info request.ArticleSearch) (interface{}, int64, error) {
	req := &search.Request{
		//一个Query只能是一种类型，这就是为什么后面写req.Query.Bool = boolQuery
		Query: &types.Query{},
	}
	boolQuery := &types.BoolQuery{}

	if info.Query != "" {
		//Should 关键词搜索
		boolQuery.Should = []types.Query{
			//match是全文检索，会计算_score,适合关键词搜索、模糊匹配、分词匹配
			{Match: map[string]types.MatchQuery{"title": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"keyword": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"abstract": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"content": {Query: info.Query}}},
		}
	}

	if info.Tag != "" {
		//Must 标签必须匹配
		boolQuery.Must = []types.Query{
			{Match: map[string]types.MatchQuery{"tags": {Query: info.Tag}}},
		}
	}

	if info.Category != "" {
		//Filter 分类精确过滤
		boolQuery.Filter = []types.Query{
			//term是 精确匹配，不会计算 _score
			{Term: map[string]types.TermQuery{"category": {Value: info.Category}}},
		}
	}
	//match是“像 Google 一样搜”，term是“像数据库 where 一样查”

	//如果 bool查询里什么子句都没有，就不能用 bool，必须改成 match_all返回所有文档。
	if boolQuery.Should != nil || boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
	}

	if info.Sort != "" {
		var sortField string
		switch info.Sort {
		case "time":
			sortField = "created_at"
		case "view":
			sortField = "views"
		case "comment":
			sortField = "comments"
		case "like":
			sortField = "likes"
		default:
			sortField = "created_at"
		}

		var order sortorder.SortOrder
		if info.Order != "asc" {
			order = sortorder.Desc
		} else {
			order = sortorder.Asc
		}

		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					sortField: {Order: &order},
				},
			},
		}
	}
	option := other.EsOption{
		PageInfo:       info.PageInfo,
		Index:          elasticsearch.ArticleIndex(),
		Request:        req,
		SourceIncludes: []string{"cover", "title", "abstract", "created_at", "categroy", "tags", "views", "comments", "likes"},
	}
	//context.TODO()和 context.Background()行为完全一致,但一般作为占位符，此处EsPagination明明是被“别人调用的（api层）”，却自己决定用 TODO()
	// return utils.EsPagination(context.TODO(), option)

	//应该在接收参数中传入ctx,在api层创建并传入ctx
	return utils.EsPagination(ctx, option)
}
func (articleService *ArticleService) ArticleCategory() ([]database.ArticleCategory, error) {
	var category []database.ArticleCategory
	//查询 article_categories表的全部数据
	if err := global.DB.Find(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}
func (articleService *ArticleService) ArticleTags() ([]database.ArticleTag, error) {
	var tags []database.ArticleTag
	if err := global.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
func (articleService *ArticleService) ArticleCreate(req request.ArticleCreate) error {
	b, err := articleService.Exist(req.Title)
	if err != nil {
		return err
	}
	if b {
		return errors.New("this article already exists")
	}

	//文章不存在则构建文章
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToCreate := elasticsearch.Article{
		CreatedAt: now,
		UpdatedAt: now,
		Cover:     req.Cover,
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}

	//启动事务确保数据一致性
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//更新文章类别的计数
		if err := articleService.UpdateCategoryCount(tx, "", articleToCreate.Category); err != nil {
			return err
		}
		//更新文章标签的计数
		if err := articleService.UpdateTagsCount(tx, []string{}, articleToCreate.Tags); err != nil {
			return err
		}

		//修改图片表中的图片类别
		if err := utils.ChangeImagesCategory(tx, []string{articleToCreate.Cover}, appTypes.Cover); err != nil {
			return err
		}
		illustrations, err := utils.FindIllustrations(articleToCreate.Content)
		if err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, illustrations, appTypes.Illustration); err != nil {
			return err
		}

		return articleService.Create(&articleToCreate)
	})
}
func (articleService *ArticleService) ArticleDelete(req request.ArticleDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			//拿到要删除的文章
			articleToDelete, err := articleService.Get(id)
			if err != nil {
				return err
			}
			//更新文章类别表
			if err := articleService.UpdateCategoryCount(tx, articleToDelete.Category, ""); err != nil {
				return err
			}
			//更新文章标签表
			if err := articleService.UpdateTagsCount(tx, articleToDelete.Tags, []string{}); err != nil {
				return err
			}
			//更新图片类别为未使用
			if err := utils.InitImagesCategory(tx, []string{articleToDelete.Cover}); err != nil {
				return err
			}

			//拿到文章内容中的插图,修改图片类别为未使用
			illustrations, err := utils.FindIllustrations(articleToDelete.Content)
			if err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, illustrations); err != nil {
				return err
			}
			//同时删除该文章下的所有评论
			comments, err := ServiceGroupApp.CommentService.CommentInfoByArticleId(request.CommentInfoByArticleId{ArticleID: id})
			if err != nil {
				return err
			}
			for _, comment := range comments {
				if err := ServiceGroupApp.CommentService.DeleteCommentAndChildren(tx, comment.ID); err != nil {
					return err
				}
			}
			// 删除该文章的所有收藏记录
			if err := tx.Where("article_id = ?", id).Delete(&database.ArticleLike{}).Error; err != nil {
				return err
			}
		}
		//删除 Redis 中的浏览量缓存（在事务外执行，不影响主事务）
		//在 ES 删除前清理 Redis（避免事务回滚后 Redis 已删但 ES 没删）
		for _, id := range req.IDs {
			global.Redis.HDel(context.Background(), "article_views", id)
		}
		return articleService.Delete(req.IDs)
	})
}
func (articleService *ArticleService) ArticleUpdate(req request.ArticleUpdate) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToUpdate := struct {
		UpdatedAt string   `json:"updated_at"`
		Cover     string   `json:"cover"`
		Title     string   `json:"title"`
		Keyword   string   `json:"keyword"`
		Category  string   `json:"category"`
		Tags      []string `json:"tags"`
		Abstract  string   `json:"abstract"`
		Content   string   `json:"content"`
	}{
		UpdatedAt: now,
		Cover:     req.Cover,
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		oldArticle, err := articleService.Get(req.ID)
		if err != nil {
			return err
		}
		if err := articleService.UpdateCategoryCount(tx, oldArticle.Category, articleToUpdate.Category); err != nil {
			return err
		}
		if err := articleService.UpdateTagsCount(tx, oldArticle.Tags, articleToUpdate.Tags); err != nil {
			return err
		}

		//更新封面
		if articleToUpdate.Cover != oldArticle.Cover {
			if err := utils.InitImagesCategory(tx, []string{oldArticle.Cover}); err != nil {
				return err
			}
			if err := utils.ChangeImagesCategory(tx, []string{articleToUpdate.Cover}, appTypes.Cover); err != nil {
				return err
			}
		}
		//更新文章中的图片
		oldIllustrations, err := utils.FindIllustrations(oldArticle.Content)
		if err != nil {
			return err
		}
		newIllustrations, err := utils.FindIllustrations(articleToUpdate.Content)
		if err != nil {
			return err
		}
		addedIllustrations, removedIllustrations := utils.DiffArrays(oldIllustrations, newIllustrations)
		//将移除的图片置为未使用，图片类型为未使用才能够将其删除
		if err := utils.InitImagesCategory(tx, removedIllustrations); err != nil {
			return err
		}
		//将新增的图片置为插图类型
		if err := utils.ChangeImagesCategory(tx, addedIllustrations, appTypes.Illustration); err != nil {
			return err
		}

		return articleService.Update(req.ID, articleToUpdate)
	})
}

func (articleService *ArticleService) ArticleList(info request.ArticleList) (list interface{}, total int64, err error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}

	// 根据标题查询
	if info.Title != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"title": {Query: *info.Title}}})
	}

	// 根据简介查询
	if info.Abstract != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"abstract": {Query: *info.Abstract}}})
	}

	// 根据类别筛选
	if info.Category != nil {
		boolQuery.Filter = []types.Query{
			{
				Term: map[string]types.TermQuery{
					"category": {Value: info.Category},
				},
			},
		}
	}

	// 根据条件执行查询
	if boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"created_at": {Order: &sortorder.Desc},
				},
			},
		}
	}

	option := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  req,
	}
	return utils.EsPagination(context.TODO(), option)
}
