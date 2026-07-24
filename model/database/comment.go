package database

import (
	"context"
	"server/global"
	"server/model/elasticsearch"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Comment 评论表
type Comment struct {
	gorm.Model
	ArticleID string    `json:"article_id"` // 文章 ID，评论与文章关联，但文章存储在es中
	PID       *uint     `json:"p_id"`       // 父评论 ID，为null则代表是根评论
	PComment  *Comment  `json:"-" gorm:"foreignKey:PID"`
	Children  []Comment `json:"children" gorm:"foreignKey:PID"` // 子评论
	//以下为评论主体
	UserUUID uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`                  // 用户 uuid
	User     User      `json:"user" gorm:"foreignKey:UserUUID;references:UUID"` // 关联的用户，references:UUID：去 User 表中查找 uuid 列，看它是否等于当前评论的 user_uuid。
	Content  string    `json:"content"`                                         // 内容
}

// 创建和删除评论时需要更新文章评论数
func (c *Comment) AfterCreate(_ *gorm.DB) error {
	source := "ctx._source.comments += 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), c.ArticleID).Script(&script).Do(context.TODO())
	return err
}
func (c *Comment) BeforeDelete(tx *gorm.DB) error {
	// DeleteCommentAndChildren 在删除前已加载完整的 Comment，c.ArticleID 可直接使用
	if c.ArticleID == "" {
		return nil
	}
	source := "ctx._source.comments -= 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), c.ArticleID).Script(&script).Do(context.TODO())
	return err
}
