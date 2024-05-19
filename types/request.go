/*
Package type - NekoBlog backend server types.
This file is for user related types.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserAuthBody 认证请求体
type UserAuthBody struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// UserRegisterBody 注册请求体
type UserUpdatePasswordBody struct {
	UserAuthBody        // 认证请求体
	NewPassword  string `json:"new_password"` // 新密码
}

// UserUpdateProfileBody 更新用户资料请求体
type UserUpdateProfileBody struct {
	NickName *string `json:"nickname"` // 昵称
	Birth    *uint64 `json:"birth"`    // 出生日期
	Gender   *string `json:"gender"`   // 性别
}

// CommentCreatebody 创建评论请求体
type UserCommentCreateBody struct {
	PostID  *uint64 `json:"post_id" form:"post_id"` // 博文ID
	Content string  `json:"content" form:"content"` // 内容
}

// UserCommentUpdateBody 更新评论请求体
type UserCommentUpdateBody struct {
	CommentID *uint64 `json:"comment_id" form:"comment_id"` // 评论ID
	Content   string  `json:"content" form:"content"`       // 内容
}

// UserPostInfo 创建博文请求体
type UserPostInfo struct {
	UID   uint   `json:"id"`    // 用户ID
	Title string `json:"title"` // 标题

}

// PostURLImageBody 创建博文请求体
type PostURLImageBody struct {
	URL string `json:"url" form:"url"` // 图片URL
}

// PostCreateBody 创建博文请求体
type PostCreateBody struct {
	Title   string             `json:"title" form:"title"`       //标题
	Content string             `json:"content" form:"content"`   //内容
	Images  []string           `json:"images" form:"images"`     // 上传图片的UUID
	TopicID primitive.ObjectID `json:"topic_id" form:"topic_id"` // 话题ID
}

// UserCommentDeleteBody 创建博文请求体
type UserCommentDeleteBody struct {
	CommentID *uint64 `json:"comment_id" form:"comment_id"` // 评论ID
}

// ReplyCreateBody 创建回复评论请求体
type ReplyCreateBody struct {
	CommentID     uint64 `json:"comment_id" form:"comment_id"`           // 博文ID
	ParentReplyID uint64 `json:"parent_reply_id" form:"parent_reply_id"` // 回复ID
	Content       string `json:"content" form:"content"`                 // 内容
}

// UserCommentUpdateBody 更新评论请求体
type UserReplyUpdateBody struct {
	ReplyID uint64 `json:"reply_id" form:"reply_id"` // 回复ID
	Content string `json:"content" form:"content"`   // 内容
}

// UserReplyDeleteBody 删除博文请求体
type UserReplyDeleteBody struct {
	ReplyID uint64 `json:"reply_id" form:"reply_id"` // 回复ID
}

// TopicCreateBody 创建话题请求体
type TopicCreateBody struct {
	Name           string `json:"name" form:"name"`                         // 话题名
	Description    string `json:"description" form:"description"`           // 话题描述
	BundledGroupID string `json:"bundled_group_id" form:"bundled_group_id"` // 绑定的群组ID
}

// TopicDeleteBody 删除话题请求体
type TopicDeleteBody struct {
	TopicID string `json:"topic_id" form:"topic_id"` // 话题ID
}

// TopicLikeBody 点赞话题请求体
type TopicLikeBody struct {
	TopicID primitive.ObjectID `json:"topic_id" form:"topic_id"` // 话题ID
}

// TopicUnLikeBody 取消点赞话题请求体
type TopicDisLikeBody struct {
	TopicID primitive.ObjectID `json:"topic_id" form:"topic_id"` // 话题ID
}
