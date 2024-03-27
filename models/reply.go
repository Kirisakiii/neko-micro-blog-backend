/*
Package models - NekoBlog backend server database models
This file is for comment related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ReplyInfo 评论信息模型
type ReplyInfo struct {
	gorm.Model                   // 基本模型
	CommentID      uint64        `gorm:"column:comment_id"`             // 博文ID
	ParentReplyID  *uint64       `gorm:"column:reply_to_reply_id"`      // 父回复ID
	UID            uint64        `gorm:"column:uid"`                    // 用户ID
	ParentReplyUID *uint64       `gorm:"column:parent_reply_uid"`       // 父回复UID
	Content        string        `gorm:"column:content"`                // 内容
	Like           pq.Int64Array `gorm:"column:like;type:bigint[]"`     // 点赞数 记录UID
	Dislike        pq.Int64Array `gorm:"column:dislike;type:bigint[]"`  // 踩数 记录UID
	IsPublic       bool          `gorm:"column:is_public;default:true"` // 是否公开
	// Share   uint64 `gorm:"column:share"`                             // 分享数 暂时不实现
}
