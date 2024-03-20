/*
Package services - NekoBlog backend server services.
This file is for user related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package services

import (
	"errors"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
)

// ReplyService 用户服务
type ReplyService struct {
	replyStore *stores.ReplyStore
}

// NewReplayService 返回一个新的评论服务实例。
//
// 返回：
//   - *ReplyService: 返回一个指向新的评论服务实例的指针。
func (factory *Factory) NewReplyService() *ReplyService {
	return &ReplyService{
		replyStore: factory.storeFactory.NewReplyStore(),
	}
}

// NewReplyService 创建评论
//
// 参数：
//   - uid：用户ID
//   - commentID: 评论编号
//   - replyToReplyID: 回复编号
//   - content: 回复内容
//   - commentStore：绑定comment来调用方法
//
// 返回值：
//
//	-error 创建失败返回创建失败时候的具体信息
func (service *ReplyService) CreateReply(uid, commentID, parentReplyID uint64, content string, commentStore *stores.CommentStore, userStore *stores.UserStore) error {
	// 校验评论是否存在
	isExist, err := commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("comment does not exist")
	}

	var parentReplyUIDField *uint64 = nil
	// 校验回复是否存在
	if parentReplyID != 0 {
		isExist, err := service.replyStore.ValidateReplyExistence(commentID, parentReplyID)
		if err != nil {
			return err
		}
		if !isExist {
			return errors.New("reply does not exist")
		}
		parentReplyInfo, err := service.replyStore.GetReply(parentReplyID)
		if err != nil {
			return err
		}
		parentReplyUIDField = &parentReplyInfo.UID
	}

	var parentReplyIDField *uint64 = nil
	if parentReplyID != 0 {
		parentReplyIDField = &parentReplyID
	}

	// 调用存储层的方法存储评论
	err = service.replyStore.CreateReply(uid, commentID, parentReplyIDField, parentReplyUIDField, content)
	if err != nil {
		return err
	}
	return nil
}

// DelteeReply 修改回复
//
// 参数：
//   - replyID：评论ID
//
// 返回值：
//   - error 如果评论存在返回修改回复时候的信息
func (service *ReplyService) DeleteReply(uid, replyID uint64) error {
	// 调用评论存储中的删除回复方法
	err := service.replyStore.DeleteReply(uid, replyID)
	if err != nil {
		// 如果发生错误，则返回错误
		return err
	}

	// 如果没有发生错误，则返回 nil
	return nil
}

// UpdateReply 修改回复
//
// 参数：
//   - replyID：评论ID
//   - content: 博文内容
//
// 返回值：
//
//	-error 如果评论存在返回修改回复时候的信息
func (service *ReplyService) UpdateReply(uid, replyID uint64, content string) error {
	// 调用数据库或其他存储方法更新评论内容
	err := service.replyStore.UpdateReply(uid, replyID, content)
	if err != nil {
		return err
	}

	// 如果更新成功，返回nil
	return nil
}

// GetReplyList 获取回复列表
//
// 参数：
//   - commentID：评论ID
//
// 返回值：
//   - []uint64：回复列表
//   - error：获取失败返回错误
func (service *ReplyService) GetReplyList(commentID uint64) ([]uint64, error) {
	// 调用数据库或其他存储方法获取评论列表
	replyList, err := service.replyStore.GetReplyList(commentID)
	if err != nil {
		return nil, err
	}
	replyListUint64 := make([]uint64, len(replyList))
	for index, reply := range replyList {
		replyListUint64[index] = uint64(reply.ID)
	}

	// 如果获取成功，返回评论列表
	return replyListUint64, nil
}

// GetReplyDetail 获取回复
//
// 参数：
//   - replyID：回复ID
//
// 返回值：
//   - models.ReplyInfo：回复信息
//   - error：获取失败返回错误
func (service *ReplyService) GetReplyDetail(replyID uint64) (models.ReplyInfo, error) {
	// 调用数据库或其他存储方法获取评论
	reply, err := service.replyStore.GetReply(replyID)
	if err != nil {
		return models.ReplyInfo{}, err
	}

	// 如果获取成功，返回评论
	return reply, nil
}