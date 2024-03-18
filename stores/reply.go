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

// Createreply 存储reply
//
// 参数 ：
//   - uid：用户id,
//   - username: 用户名，
//   - commentID: 评论id，
//   - replyToReplyID: 回复id，
//   - content: 评论内容，
//
// 返回：
//
//	-error 正确返回nil
func (store *ReplyStore) CreateReply(uid uint64, username string, commentID *uint64, replyToReplyID *uint64, content string) error {
	newReply := &models.ReplyInfo{
		CommentID:      commentID,
		ReplyToReplyID: replyToReplyID,
		Username:       username,
		Content:        content,
		UID:            uid,
		Like:           nil,
		Dislike:        nil,
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
//	参数：
//	- replyToReplyID: 评论ID
//
// 返回值：
//   - error：如果回复存在返回true，不存在判断具体的错误类型返回false
func (store *ReplyStore) ValidateReplyExistence(replyToReplyID uint64) (bool, error) {
	var reply models.ReplyInfo
	result := store.db.Where("id = ?", replyToReplyID).First(&reply)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	// 返回错误类型
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}


func (store *ReplyStore) DeleteReply(replyID uint64) error {
	var replyInfo models.ReplyInfo

	// 查询回复信息
	if err := store.db.Where("id = ?", replyID).First(&replyInfo).Error; err != nil {
		return err // 如果没有找到对应的回复，返回错误
	}

	// 根据回复类型删除相应记录
	if *replyInfo.CommentID != 0 {
		// 如果 CommentID 不为 0，则表示这是一个评论回复
		return store.db.Where("id = ?", replyID).Unscoped().Delete(&models.ReplyInfo{}).Error
	} else {
		// 否则，表示这是一个回复的回复
		return store.db.Where("reply_id = ?", replyID).Unscoped().Delete(&models.ReplyInfo{}).Error
	}
}


// UpdateReply 修改回复
//
//    参数：
//    - replyID: 评论ID
//    - content: 修改内容
//
// 返回值：
//   - error：如果回复存在返回true，不存在判断具体的错误类型返回false
func (store *ReplyStore) UpdateReply(replyID uint64, content string) error {
    ReplyInfo := new(models.ReplyInfo)
    result := store.db.Where("id = ?", replyID).First(ReplyInfo)
    if result.Error != nil {
        return result.Error
    }

    ReplyInfo.Content = content
    result = store.db.Save(ReplyInfo)
    if result.Error != nil {
        return result.Error
    }
    return nil
}