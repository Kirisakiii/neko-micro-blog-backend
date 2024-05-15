/*
Package models - NekoBlog backend server database models
This file is for comment related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- CBofJOU<2023122312@jou.edu.cn>
*/
package models

import (
	"time"
)

// FollowInfo 关注信息模型
type FollowInfo struct {
	UserID      uint64    `bson:"uid"`         // 关注ID
	FollowedID  uint64    `bson:"followed_id"` // 被关注者ID
	FollowedAt  time.Time `bson:"followed_at"` // 关注时间
}