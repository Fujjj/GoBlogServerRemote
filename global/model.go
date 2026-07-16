package global

import (
	"time"

	"gorm.io/gorm"
)

type MODEL struct { //使用了gorm.Model代替
	ID        uint           `json:"id" gorm:"primarykey"` // 主键 ID
	CreatedAt time.Time      `json:"created_at"`           // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`           // 更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`       // 删除时间
}
