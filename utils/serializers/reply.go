package serializers

import "github.com/Kirisakiii/neko-micro-blog-backend/models"

// ReplyListResponse 回复列表响应结构
type ReplyListResponse struct {
	IDs []uint64 `json:"ids"`
}

// NewReplyListResponse 创建新的回复列表响应
//
// 参数：
//   - replies：回复ID列表
//
// 返回值：
//   - ReplyListResponse：新的回复列表响应结构
func NewReplyListResponse(replies []uint64) ReplyListResponse {
	return ReplyListResponse{IDs: replies}
}

// ReplyDetailResponse 回复信息响应结构
type ReplyDetailResponse struct {
	CreateTime     int64   `json:"create_time"`      // 创建时间
	CommentID      uint64  `json:"comment_id"`       // 评论ID
	UID            uint64  `json:"uid"`              // 用户ID
	ParentReplyID  *uint64 `json:"parent_reply_id"`  // 父回复ID
	ParentReplyUID *uint64 `json:"parent_reply_uid"` // 父回复UID
	Content        string  `json:"content"`          // 内容
	// Like           int     `json:"like"`             // 点赞数
	// Dislike        int     `json:"dislike"`          // 踩数
}

// NewReplyDetailResponse 创建新的回复信息响应
//
// 参数：
//   - model：回复信息模型
//
// 返回值：
//   - *ReplyDetailResponse：新的回复信息响应结构
func NewReplyDetailResponse(reply models.ReplyInfo) ReplyDetailResponse {
	// 创建一个新的 ReplyDetailResponse 实例
	profileData := ReplyDetailResponse{
		CreateTime:     reply.CreatedAt.Unix(),
		CommentID:      reply.CommentID,
		UID:            reply.UID,
		ParentReplyID:  reply.ParentReplyID,
		ParentReplyUID: reply.ParentReplyUID,
		Content:        reply.Content,
		// Like:           len(reply.Like),
		// Dislike:        len(reply.Dislike),
	}

	return profileData
}
