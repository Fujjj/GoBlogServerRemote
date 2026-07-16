package database

import (
	"gorm.io/gorm"
)

// ArticleLike 文章收藏表
type ArticleLike struct {
	gorm.Model
	ArticleID string `json:"article_id"` // 文章 ID
	UserID    uint   `json:"user_id"`    // 用户 ID
	User      User   `json:"-" gorm:"foreignKey:UserID"`
}
