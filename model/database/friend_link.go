package database

import "gorm.io/gorm"

// FriendLink 友链表
type FriendLink struct {
	gorm.Model
	Logo        string `json:"logo" gorm:"size:255"` // Logo
	Image       Image  `json:"-" gorm:"foreignKey:Logo;references:URL"`
	Link        string `json:"link"`        // 链接
	Name        string `json:"name"`        // 名称
	Description string `json:"description"` // 描述
}
