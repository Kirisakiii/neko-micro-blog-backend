/*
Package models - NekoBlog backend server database models
This file is for post related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// PostInfo 博文信息模型
type PostInfo struct {
	gorm.Model                      // 基本模型
	ParentPostID *uint64            `gorm:"column:parent_post_id"`          // 转发自文章ID
	UID          uint64             `gorm:"column:uid"`                     // 用户ID
	IpAddrress   *string            `gorm:"column:ip_address"`              // IP地址
	Title        string             `gorm:"column:title"`                   // 标题
	Content      string             `gorm:"column:content"`                 // 内容
	TopicID      primitive.ObjectID `bson:"topic_id,omitempty"`             // 所属话题ID
	Images       pq.StringArray     `gorm:"column:images;type:text[]"`      // 图片
	Like         pq.Int64Array      `gorm:"column:like;type:bigint[]"`      // 点赞数 记录UID
	Favourite    pq.Int64Array      `gorm:"column:favourite;type:bigint[]"` // 收藏数 记录UID
	Farward      pq.Int64Array      `gorm:"column:farward;type:bigint[]"`   // 转发数 记录UID
	IsPublic     bool               `gorm:"column:is_public;default:true"`  // 是否公开
	// Share     uint64 `gorm:"column:share"`                           // 分享数 暂时不实现
}
