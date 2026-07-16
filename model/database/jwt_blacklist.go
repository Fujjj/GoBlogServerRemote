package database

import (
	"gorm.io/gorm"
)

// JwtBlacklist JWT 黑名单表
type JwtBlacklist struct {
	gorm.Model
	Jwt string `json:"jwt" gorm:"type:text"` // Jwt
}
