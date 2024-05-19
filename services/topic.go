/*
Package services - NekoBlog backend server services.
This file is for topic related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/

package services

import (
	"errors"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TopicService 用户服务
type TopicService struct {
	topicStore *stores.TopicStore
}

// NewTopicService 返回一个新的 TopicService 实例。
//
// 返回值：
//   - *TopicService：新的 TopicService 实例。
func (factory *Factory) NewTopicService() *TopicService {
	return &TopicService{
		topicStore: factory.storeFactory.NewTopicStore(),
	}
}

// CreateTopic 返回一个新的 TagService 实例。
//
// 参数：
//   - userID (uint64)：用户的 ID。
//   - reBody (types.TagCreateBody)：创建标签的请求体。
//   - postStore (stores.PostStore)：一个 PostStore 实例。
//
// 返回值：
//   - *TagService：新的 TagService 实例。
func (service *TopicService) CreateTopic(userID uint64, reBody *types.TopicCreateBody) (primitive.ObjectID, error) {

	// 调用存储层的方法
	tagID, err := service.topicStore.CreateTopic(userID, reBody.PostID, reBody.Description)

	// 检查存储层是否返回错误
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return tagID, nil

}

// DeleteTarget 删除标签
//
// 参数：
//   - tagID (primitive.ObjectID)：标签的 ID。
//
// 返回值：
//   - error：如果删除标签时发生错误，则返回一个错误。
func (service *TopicService) DeleteTarget(topic primitive.ObjectID) error {
	// 调用存储层的方法
	err := service.topicStore.DeleteTopic(topic)

	// 检查存储层是否返回错误
	if err != nil {
		return err
	}
	return nil
}

// NewTopicListResponse 返回一个新的 TopicListResponse 实例
//
// 参数：
//   - topciID: 话题ID
//
// 返回值：
//   - []models.TopicInfo: 目标列表
//   - error: 错误信息
func (service *TopicService) GetTopicList(topicID primitive.ObjectID) ([]models.TopicInfo, error) {
	// 调用存储层的方法
	tagList, err := service.topicStore.GetTopicList(topicID)

	// 检查存储层是否返回错误
	if err != nil {
		return nil, err
	}
	return tagList, nil
}

// NewLikeTopicHandler 点赞目标
//
// 参数：
//   - topicID (primitive.ObjectID)：目标 ID。
//
// 返回值：
//   - error：如果点赞目标时发生错误，则返回一个错误。
func (service *TopicService) NewLikeTopicHandler(topicID primitive.ObjectID) error {

	// 检查topic是否存在
	exists, err := service.topicStore.ValidateTopicExistence(topicID)

	// 如果topic不存在，则返回错误
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("topic does not exist")
	}

	// 调用存储层的方法
	err = service.topicStore.LikeTopic(topicID)

	// 检查存储层是否返回错误
	if err != nil {
		return err
	}
	return nil
}

// NewCancelLikeTopicHandler 取消点赞目标
//
// 参数：
//   - topicID (primitive.ObjectID)：目标 ID
func (service *TopicService) NewCancelLikeTopicHandler(topicID primitive.ObjectID) error {

	// 检查topic是否存在
	exists, err := service.topicStore.ValidateTopicExistence(topicID)

	// 如果topic不存在，则返回错误
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("topic does not exist")
	}

	// 调用存储层的方法
	err = service.topicStore.CancelLikeTopic(topicID)

	// 检查存储层是否返回错误
	if err != nil {
		return err
	}

	return nil
}

// NewDislikeTopicHandler 点踩数目标
//
// 参数：
//   - topicID (primitive.ObjectID)：目标 ID
func (service *TopicService) NewDislikeTopicHandler(topicID primitive.ObjectID) error {

	// 检查topic是否存在
	exists, err := service.topicStore.ValidateTopicExistence(topicID)

	// 如果topic不存在，则返回错误
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("topic does not exist")
	}

	// 调用存储层的方法
	err = service.topicStore.DislikeTopic(topicID)

	// 检查存储层是否返回错误
	if err != nil {
		return err
	}

	return nil
}

// NewCancelDislikeTopicHandler 取消点踩数目标
//
// 参数：
//   - topicID (primitive.ObjectID)：目标 ID
//
// 返回值：
//   - error：如果取消点踩数目标时发生错误，则返回一个错误。
func (service *TopicService) GetHotTopics(limit int) ([]models.TopicInfo, error) {
	return service.topicStore.GetHotTopics(limit)
}
