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

// GetTopicList 获取话题列表
//
// 参数：
//   - from (primitive.ObjectID)：起始 ID。
//   - length (uint64)：长度。
//
// 返回值：
//   - []models.TopicInfo：话题列表。
//   - error：如果获取话题列表时发生错误，则返回一个错误。
func (service *TopicService) GetTopicList(from primitive.ObjectID, length uint64) ([]models.TopicInfo, error) {
	// 调用存储层的方法
	return service.topicStore.GetTopicList(from, length)
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
func (service *TopicService) CreateTopic(userID uint64, reqBody types.TopicCreateBody, bundledGroupID primitive.ObjectID) (primitive.ObjectID, error) {
	// 检查是否存在重复的标签
	exists, err := service.topicStore.ValidateTopicExistenceByName(reqBody.Name)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if exists {
		return primitive.NilObjectID, errors.New("topic already exists")
	}

	// TODO: valid if bundledGroupID exists

	// 调用存储层的方法
	return service.topicStore.CreateTopic(userID, reqBody.Name, reqBody.Description, bundledGroupID)
}

// DeleteTopic 删除标签
//
// 参数：
//   - tagID (primitive.ObjectID)：标签的 ID。
//
// 返回值：
//   - error：如果删除标签时发生错误，则返回一个错误。
func (service *TopicService) DeleteTopic(topicID primitive.ObjectID) error {
	// 调用存储层的方法
	return service.topicStore.DeleteTopic(topicID)
}

// NewTopicListResponse 返回一个新的 TopicListResponse 实例
//
// 参数：
//   - topciID: 话题ID
//
// 返回值：
//   - []models.TopicInfo: 目标列表
//   - error: 错误信息
func (service *TopicService) GetTopicDetail(topicID primitive.ObjectID) (models.TopicInfo, error) {
	// 调用存储层的方法
	return service.topicStore.GetTopicDetail(topicID)
}

// NewLikeTopic 点赞目标
//
// 参数：
//   - topicID (primitive.ObjectID)：目标 ID。
//
// 返回值：
//   - error：如果点赞目标时发生错误，则返回一个错误。
func (service *TopicService) NewLikeTopic(topicID primitive.ObjectID, userID uint64) error {

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
	return service.topicStore.LikeTopic(topicID, userID)
}

// NewCancelLikeTopic 取消点赞目标
//
// 参数：
//   - topicID (primitive.ObjectID)：目标 ID
func (service *TopicService) NewCancelLikeTopic(topicID primitive.ObjectID, userID uint64) error {
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
	return service.topicStore.CancelLikeTopic(topicID, userID)
}

// NewDislikeTopic 点踩数目标
//
// 参数：
//   - topicID (primitive.ObjectID)：目标 ID
func (service *TopicService) NewDislikeTopic(topicID primitive.ObjectID, userID uint64) error {

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
	err = service.topicStore.DislikeTopic(topicID, userID)

	// 检查存储层是否返回错误
	if err != nil {
		return err
	}

	return nil
}

// CancelDislikeTopic 取消点踩数目标
//
// 参数：
//   - topicID (primitive.ObjectID)：目标 ID
//   - userID (uint64)：用户 ID
//
// 返回值：
//   - error：如果取消点踩数目标时发生错误，则返回一个错误。
func (service *TopicService) CancelDislikeTopic(topicID primitive.ObjectID, userID uint64) error {
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
	return service.topicStore.CancelDislikeTopic(topicID, userID)
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
