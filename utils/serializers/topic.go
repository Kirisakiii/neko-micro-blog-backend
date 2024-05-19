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
	ID string `bson:"_id,omitempty" json:"id"`
}

// TopicDetailResponse 用于将 TopicInfo 转换为 JSON 格式的结构体
type TopicDetailResponse struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name,omitempty" json:"name"`
	Description    string             `bson:"description,omitempty" json:"description"`
	BundledGroupID primitive.ObjectID `bson:"bundled_group_id,omitempty" json:"bundled_group_id"`
	CreatorID      uint64             `bson:"creator_id,omitempty" json:"creator_id"`
	Like           uint64             `bson:"like,omitempty" json:"like"`
	DisLike        uint64             `bson:"dislike,omitempty" json:"dislike"`
}

// TopicListResponse 用于将 TopicInfo 转换为 JSON 格式的结构体
type TopicListResponse struct {
	IDs []primitive.ObjectID `bson:"_id,omitempty" json:"ids"`
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
		ID: TopicID.Hex(),
	}
	return resp
}

// NewTopicDetailResponse 用于创建 TopicDetailResponse 实例
//
// 参数：
// - TagInfo: 目标标签的信息
//
// 返回值：
// - TagDetailResponse: 创建的 TagDetailResponse 实例
func NewTopicDetailResponse(TopicInfo models.TopicInfo) TopicDetailResponse {
	return TopicDetailResponse{
		ID:             TopicInfo.ID,
		Name:           TopicInfo.Name,
		Description:    TopicInfo.Description,
		BundledGroupID: TopicInfo.BundledGroupID,
		CreatorID:      TopicInfo.CreatorID,
		Like:           TopicInfo.Like,
		DisLike:        TopicInfo.DisLike,
	}
}

// NewTopicListResponse 用于创建 TopicListResponse 实例
//
// 参数：
// - TopicList: 话题列表
//
// 返回值：
// - []TopicDetailResponse: 创建的 TopicListResponse 实例
func NewTopicListResponse(TopicList []models.TopicInfo) TopicListResponse {
	var resp TopicListResponse
	for _, TopicInfo := range TopicList {
		resp.IDs = append(resp.IDs, TopicInfo.ID)
	}
	return resp
}
