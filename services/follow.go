/*
Package services - NekoBlog backend server services.
This file is for follow related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
- CBofJOU<2023122312@jou.edu.cn>
*/
package services

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
)

// CommentService 评论服务
type FollowService struct {
	followStore *stores.FollowStore
}

// NewFollowService 返回一个新的评论服务实例。
//
// 返回：
//   - *FollowService: 返回一个指向新的评论服务实例的指针。
func (factory *Factory) NewFollowService() *FollowService {
	return &FollowService{
		followStore: factory.storeFactory.NewFollowStore(),
	}
}

// FollowUser 关注用户
//
// 参数：
//   - uid：用户ID
//   - followedID：被关注用户ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *FollowService) FollowUser(uid, followedID uint64) error {
	// 调用follow存储中的点赞方法
	return service.followStore.FollowUser(uid, followedID)
}

// CancelFollowUser 取消关注用户
//
// 参数：
//   - uid：用户ID
//   - followedID：被关注用户ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *FollowService) CancelFollowUser(uid, followedID uint64) error {
	// 调follow存储中的关注方法
	return service.followStore.CancelFollowUser(uid, followedID)
}

// GetFOllowList 获取关注列表
//
// 返回值：
//   - 成功则返回关注列表
//   - 失败返回nil
func (service *FollowService) GetFollowList(userID uint64) ([]models.FollowInfo, error) {
	return service.followStore.GetFollowList(userID)
}

// GetFollowCountByUID 获取用户的关注数量
// 
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - 成功则返回关注数量
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *FollowService) GetFollowCountByUID(uid uint64) (int, error) {
    count, err := service.followStore.GetFollowedsByUID(uid)
    if err != nil {
        return 0, err
    }
    return count, nil
}

// GetFOllowerList 获取关注列表
//
// 返回值：
//   - 成功则返回粉丝列表
//   - 失败返回nil
func (service *FollowService) GetFollowerList(userID uint64) ([]models.FollowInfo, error) {
	return service.followStore.GetFollowerList(userID)
}

// GetFollowCountByUID 获取用户的粉丝数量
//
//	参数：
//	  - uid：用户ID
//
// 返回值：
//	  - 成功则返回粉丝数量
//	  - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *FollowService) GetFollowerCountByUID(uid uint64) (int, error) {
    count, err := service.followStore.GetFollowersByUID(uid)
    if err != nil {
        return 0, err
    }
    return count, nil
}