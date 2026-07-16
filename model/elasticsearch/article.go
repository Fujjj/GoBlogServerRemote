package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Article 文章表
type Article struct {
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间

	Cover    string   `json:"cover"`    // 文章封面
	Title    string   `json:"title"`    // 文章标题
	Keyword  string   `json:"keyword"`  // 文章标题-关键字
	Category string   `json:"category"` // 文章类别
	Tags     []string `json:"tags"`     // 文章标签
	Abstract string   `json:"abstract"` // 文章简介
	Content  string   `json:"content"`  // 文章内容

	Views    int `json:"views"`    // 浏览量
	Comments int `json:"comments"` // 评论量
	Likes    int `json:"likes"`    // 收藏量
}

// ArticleIndex 文章 ES 索引
func ArticleIndex() string { //es的索引相当于MySql的表，是存放数据的容器
	return "article_index"
}

// ArticleMapping 文章 Mapping 映射，类似MySql的表结构，用来规定每个字段的类型是否需要分词
func ArticleMapping() *types.TypeMapping {
	return &types.TypeMapping{
		Properties: map[string]types.Property{
			"created_at": types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")},
			"updated_at": types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")},
			"cover":      types.TextProperty{},
			//TextProperty用于模糊匹配
			"title": types.TextProperty{},
			//KeywordProperty用于精确匹配
			"keyword":  types.KeywordProperty{},
			"category": types.KeywordProperty{},
			"tags":     []types.KeywordProperty{},
			"abstract": types.TextProperty{},
			"content":  types.TextProperty{},
			"views":    types.IntegerNumberProperty{},
			"comments": types.IntegerNumberProperty{},
			"likes":    types.IntegerNumberProperty{},
		},
	}
}
