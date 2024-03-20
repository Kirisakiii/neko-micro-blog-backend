/*
Package stores - NekoBlog backend server data access objects.
This file is for user storage accessing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package stores

import (
	"errors"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ReplyStore 用户信息数据库
type ReplyStore struct {
	db *gorm.DB
}

// NewReplyStore 返回一个新的 ReplyStore 实例。
//
// 返回值：
//   - *ReplyStore：新的 ReplyStore 实例。
func (factory *Factory) NewReplyStore() *ReplyStore {
	return &ReplyStore{factory.db}
}

// CreateReply 创建回复
//
// 参数：
//   - uid：用户ID
//   - commentID: 评论编号
//   - parentReplyID: 回复编号
//   - parentReplyUID: 回复用户ID
//   - content: 回复内容
//
// 返回值：
//   - error：创建失败返回创建失败时候的具体信息
func (store *ReplyStore) CreateReply(uid, commentID uint64, parentReplyID, parentReplyUID *uint64, content string) error {
	newReply := &models.ReplyInfo{
		CommentID:      commentID,
		ParentReplyID:  parentReplyID,
		ParentReplyUID: parentReplyUID,
		Content:        content,
		UID:            uid,
		Like:           pq.Int64Array{},
		Dislike:        pq.Int64Array{},
		IsPublic:       true,
	}

	result := store.db.Create(newReply)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ValidateReplyExistence 判断回复是否存在
//
// 参数：
//   - commentID: 评论ID
//   - parentReplyID: 回复ID
//
// 返回值：
//   - bool：如果回复存在返回true，不存在返回false
//   - error：如果回复存在返回true，不存在判断具体的错误类型返回false
func (store *ReplyStore) ValidateReplyExistence(commentID, parentReplyID uint64) (bool, error) {
	// 查询回复
	var reply models.ReplyInfo
	result := store.db.Where("id = ? AND comment_id = ?", parentReplyID, commentID).First(&reply)

	// 如果没有找到对应的回复，返回否
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	// 其他错误则直接返回
	if result.Error != nil {
		return false, result.Error
	}

	// 找到了回复
	return true, nil
}

// DeleteReply 删除回复
//
// 参数：
//   - replyID：回复ID
//
// 返回值：
//   - error：删除失败返回错误
func (store *ReplyStore) DeleteReply(uid, replyID uint64) error {
	result := store.db.Model(&models.ReplyInfo{}).Where("id = ? AND uid = ?", replyID, uid).Unscoped().Delete(&models.ReplyInfo{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateReply 修改回复
//
// 参数：
//   - replyID：回复ID
//   - content: 回复内容
//
// 返回值：
//   - error：修改失败返回错误
func (store *ReplyStore) UpdateReply(uid, replyID uint64, content string) error {
	result := store.db.Model(&models.ReplyInfo{}).Where("id = ? AND uid = ?", replyID, uid).Update("content", content)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetReply 获取回复
//
// 参数：
//   - replyID：回复ID
//
// 返回值：
//   - models.ReplyInfo：回复信息
//   - error：获取失败返回错误
func (store *ReplyStore) GetReply(replyID uint64) (models.ReplyInfo, error) {
	var reply models.ReplyInfo
	result := store.db.Where("id = ?", replyID).First(&reply)
	if result.Error != nil {
		return models.ReplyInfo{}, result.Error
	}
	return reply, nil
}

// GetReplyList 获取回复列表
//
// 参数：
//   - commentID：评论ID
//
// 返回值：
//   - []models.ReplyInfo：回复列表
//   - error：获取失败返回错误
func (store *ReplyStore) GetReplyList(commentID uint64) ([]models.ReplyInfo, error) {
	var replyList []models.ReplyInfo
	result := store.db.Where("comment_id = ?", commentID).Order("id desc").Find(&replyList)
	if result.Error != nil {
		return nil, result.Error
	}
	return replyList, nil
}
