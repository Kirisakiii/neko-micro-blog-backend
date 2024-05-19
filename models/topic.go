/*
Package models - NekoBlog backend server database models
This file is for tag related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TopicInfo 话题信息模型
type TopicInfo struct {
	ID          primitive.ObjectID `bson:"_id"`         // 话题ID
	Description string             `bson:"description"` // 话题描述
	PostID      uint64             `bson:"post_id"`     // 帖子ID
	UserID      uint64             `bson:"user_id"`     // 用户ID
	Like        uint64             `bson:"like"`        // 点赞数
	DisLike     uint64             `bson:"dis_like"`    // 点踩数
}

const TOPIC_INFO_COLLECTION = "topic_info"
