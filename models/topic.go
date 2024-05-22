/*
Package models - NekoBlog backend server database models
This file is for tag related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TopicInfo 话题信息模型
type TopicInfo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"` // 话题ID
	Name        string             `bson:"name"`          // 话题名
	Description string             `bson:"description"`   // 话题描述
	CreatorID   uint64             `bson:"user_id"`       // 用户ID
	CreatedAt   time.Time          `bson:"created_at"`    // 创建时间
	Like        uint64             `bson:"like"`          // 点赞数
	DisLike     uint64             `bson:"dislike"`       // 点踩数
}

const TOPIC_INFO_COLLECTION = "topic_info"
