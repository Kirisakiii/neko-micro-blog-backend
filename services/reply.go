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
func (service *ReplyService) CreateReply(uid uint64, commentID *uint64, replyToReplyID *uint64, content string, commentStore *stores.CommentStore, userStore *stores.UserStore) error {

	// 校验评论是否存在
	if *commentID == 0 && *replyToReplyID != 0 {
		exist, err := service.replyStore.ValidateReplyExistence(*replyToReplyID)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New("comment does not exist")
		}
	} else if *replyToReplyID == 0 && *commentID != 0 {
		exist, err := commentStore.ValidateCommentExistence(*commentID)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New("reply to comment does not exist")
		}
	}

	// 根据 UID 获取 Username
	user, err := userStore.GetUserByUID(uid)
	if err != nil {
		return err
	}

	// 调用存储层的方法存储评论
	err = service.replyStore.CreateReply(uid, user.UserName, commentID, replyToReplyID, content)
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
//
//	-error 如果评论存在返回修改回复时候的信息
func (service *ReplyService) DeleteReply(replyID uint64) error {
	// TODO: 改为控制器层判断
	exists, err := service.replyStore.ValidateReplyExistence(replyID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}
	// 调用评论存储中的删除回复方法
	err = service.replyStore.DeleteReply(replyID)
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
//    -error 如果评论存在返回修改回复时候的信息
func (service *ReplyService) UpdateReply(replyID uint64, content string) error {
    //检查评论是否存在
    // TODO: 改为控制器层判断
    exists, err := service.replyStore.ValidateReplyExistence(replyID)
    if err != nil {
        return err
    }
    if !exists {
        return errors.New("comment does not exist")
    }

    // 调用数据库或其他存储方法更新评论内容
    err = service.replyStore.UpdateReply(replyID, content)
    if err != nil {
        return err
    }

    // 如果更新成功，返回nil
    return nil
}