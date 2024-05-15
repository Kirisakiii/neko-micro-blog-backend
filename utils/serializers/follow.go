package serializers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
)

type FollowListResponse struct {
	IDs []uint64 `json:"ids"`
}

// NewFollowListResponse 创建关注列表的响应
//
// 参数：
//   - 关注信息列表
//
// 返回值：
//   - 关注列表的响应
func NewFollowListResponse(followInfos []models.FollowInfo) FollowListResponse {
	var ids []uint64
	for _, followInfos := range followInfos {
		ids = append(ids, followInfos.FollowedID)
	}
	return FollowListResponse{IDs: ids}
}

// NewFollowerListResponse 创建关注列表的响应
//
// 参数：
//   - 粉丝信息列表
//
// 返回值：
//   - 粉丝列表的响应
func NewFollowerListResponse(followInfos []models.FollowInfo) FollowListResponse {
	var ids []uint64
	for _, followInfos := range followInfos {
		ids = append(ids, followInfos.UserID)
	}
	return FollowListResponse{IDs: ids}
}