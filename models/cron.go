package models

import "gorm.io/gorm"

// DeletedAvatar 头像删除等待列表
type DeletedAvatar struct {
	gorm.Model        // 基本模型
	FileName   string `gorm:"column:file_name"` // 文件名
}

type DeletedCachedImage struct {
	gorm.Model        // 基本模型
	FileName   string `gorm:"column:file_name"` // 文件名
}

//
type CachedPostImage struct {
	gorm.Model        // 基本模型
	FileName   string `gorm:"column:file_name"`   // 文件名
	ExpireTime int64  `gorm:"column:expire_time"` // 过期时间
}
