/*
Package serializers - NekoBlog backend server data serialization.
This file is for response serialization.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package serializers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTopicResponse 用于将 TopicInfo 转换为 JSON 格式的结构体
type CreateTopicResponse struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

// TopicListResponse 用于将 TopicInfo 转换为 JSON 格式的结构体
type TopicListResponse struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Description string             `bson:"description,omitempty" json:"description"`
	PostID      uint64             `bson:"post_id,omitempty" json:"post_id"`
	UserID      uint64             `bson:"user_id,omitempty" json:"user_id"`
}

// NewCreateTopicResponse 用于创建 TopicResponse 实例
//
// 参数：
// - TopicID: 话题标签的 ID
//
// 返回值：
// - CreateTopicResponse: 创建的 CreateTopicResponse 实例
func NewCreateTopicResponse(TopicID primitive.ObjectID) CreateTopicResponse {
	var resp = CreateTopicResponse{
		ID: TopicID,
	}
	return resp
}

// NewTopicListResponse 用于创建 TagListResponse 实例
//
// 参数：
// - TagInfo: 目标标签的信息
//
// 返回值：
// - TagListResponse: 创建的 TagListResponse 实例
func NewTopicListResponse(TopicInfo []models.TopicInfo) TopicListResponse {
	var resp = TopicListResponse{
		ID:          TopicInfo[0].ID,
		Description: TopicInfo[0].Description,
		PostID:      TopicInfo[0].PostID,
		UserID:      TopicInfo[0].UserID,
	}
	return resp
}
