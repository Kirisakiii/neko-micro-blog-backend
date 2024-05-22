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
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	CreatorID   uint64             `json:"creator_id"`
	Like        uint64             `json:"like"`
	DisLike     uint64             `json:"dislike"`
	RelatedPost uint64             `json:"related_post"`
}

// TopicListResponse 用于将 TopicInfo 转换为 JSON 格式的结构体
type TopicListResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Like        uint64 `json:"like"`
	RelatedPost uint64 `json:"related_post"`
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
func NewTopicDetailResponse(TopicInfo models.TopicInfo, relatedPostCount uint64) TopicDetailResponse {
	return TopicDetailResponse{
		ID:          TopicInfo.ID,
		Name:        TopicInfo.Name,
		Description: TopicInfo.Description,
		CreatorID:   TopicInfo.CreatorID,
		Like:        TopicInfo.Like,
		DisLike:     TopicInfo.DisLike,
		RelatedPost: relatedPostCount,
	}
}

// NewTopicListResponse 用于创建 TopicListResponse 实例
//
// 参数：
// - TopicList: 话题列表
//
// 返回值：
// - []TopicDetailResponse: 创建的 TopicListResponse 实例
func NewTopicListResponse(TopicList []TopicInfo) []TopicListResponse {
	var resp []TopicListResponse

	for _, topicInfo := range TopicList {
		resp = append(resp, TopicListResponse{
			ID:          topicInfo.TopicInfo.ID.Hex(),
			Name:        topicInfo.TopicInfo.Name,
			Like:        topicInfo.TopicInfo.Like,
			Description: topicInfo.TopicInfo.Description,
			RelatedPost: uint64(topicInfo.RelatedPostCount),
		})
	}
	return resp
}

// TopicInfo 包含热门话题信息和相关帖子数量
type TopicInfo struct {
	TopicInfo        models.TopicInfo
	RelatedPostCount int64
}

// HotTopicListResponse 用于将 TopicInfo 转换为 JSON 格式的结构体
type HotTopicListResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Like        uint64 `json:"like"`
	RelatedPost uint64 `json:"related_post"`
}

// NewHotTopicListResponse 用于创建 GetHotTopicListResponse 实例
//
// 参数：
// - TopicList: 话题列表
//
// 返回值：
// - []HotTopicListResponse: 创建的 HotTopicListResponse 实例
func NewHotTopicListResponse(TopicList []TopicInfo) []HotTopicListResponse {
	var resp []HotTopicListResponse

	for _, topicInfo := range TopicList {
		resp = append(resp, HotTopicListResponse{
			ID:          topicInfo.TopicInfo.ID.Hex(),
			Name:        topicInfo.TopicInfo.Name,
			Description: topicInfo.TopicInfo.Description,
			Like:        topicInfo.TopicInfo.Like,
			RelatedPost: uint64(topicInfo.RelatedPostCount),
		})
	}
	return resp
}

// TopicStatusResponse 用于将 TopicStatus 转换为 JSON 格式的结构体
type TopicStatusResponse struct {
	ID       string `json:"id"`
	Liked    bool   `json:"liked"`
	Disliked bool   `json:"disliked"`
}

// NewTopicStatusResponse 用于创建 TopicStatusResponse 实例
//
// 参数：
// - TopicID: 话题标签的 ID
// - Liked: 用户是否点赞了目标
// - Disliked: 用户是否点踩了目标
//
// 返回值：
// - TopicStatusResponse: 创建的 TopicStatusResponse 实例
func NewTopicStatusResponse(TopicID primitive.ObjectID, Liked bool, Disliked bool) TopicStatusResponse {
	return TopicStatusResponse{
		ID:       TopicID.Hex(),
		Liked:    Liked,
		Disliked: Disliked,
	}
}
