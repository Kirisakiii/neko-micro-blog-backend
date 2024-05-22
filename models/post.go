/*
Package models - NekoBlog backend server database models
This file is for post related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// PostInfo 博文信息模型
type PostInfo struct {
	gorm.Model                  // 基本模型
	ParentPostID *uint64        `gorm:"column:parent_post_id"`         // 转发自文章ID
	UID          uint64         `gorm:"column:uid"`                    // 用户ID
	IpAddrress   *string        `gorm:"column:ip_address"`             // IP地址
	Title        string         `gorm:"column:title"`                  // 标题
	Content      string         `gorm:"column:content"`                // 内容
	TopicID      *string        `gorm:"column:topic_id"`               // 所属话题ID
	Images       pq.StringArray `gorm:"column:images;type:text[]"`     // 图片
	IsPublic     bool           `gorm:"column:is_public;default:true"` // 是否公开
}
