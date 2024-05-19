/*
Package stores - NekoBlog backend server data access objects.
This file is for tag storage accessing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package stores

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TopicStore 话题信息数据库
type TopicStore struct {
	db    *gorm.DB
	rds   *redis.Client
	mongo *mongo.Client
}

// NewTopicStore 返回一个新的 TopicStore 实例。
//
// 返回值：
//   - *TopicStore：新的 TopicStore 实例。
func (factory *Factory) NewTopicStore() *TopicStore {
	return &TopicStore{
		factory.db,
		factory.rds,
		factory.mongo,
	}
}

// GetTopicByID 根据ID获取目标信息
//
// 参数：
//   - userID: 用户ID
//   - postID: 目标ID
//   - description: 目标描述
//
// 返回值：
//   - topic.ID: 话题ID
//   - error: 错误信息
func (store *TopicStore) CreateTopic(userID uint64, postID uint64, description string) (primitive.ObjectID, error) {

	// 开启事务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 插入topic信息
	topic := &models.TopicInfo{
		UserID:      userID,
		PostID:      postID,
		Description: description,
		ID:          primitive.NewObjectID(), // 初始化 ObjectID
	}

	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err := collection.InsertOne(ctx, topic)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return topic.ID, nil
}

// DeleteTopic 删除目标信息
//
// 参数：
//   - topicID: 目标ID
//
// 返回值：
//   - error: 错误信息
func (store *TopicStore) DeleteTopic(topic primitive.ObjectID) error {

	//开启事务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//删除目标
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": topic})
	if err != nil {
		return err
	}

	return nil
}

// GetTopicList 获取目标列表
//
// 参数：
//   - tagID: 话题ID
//
// 返回值：
//   - []models.TagInfo: 目标列表
//   - error: 错误信息
func (store *TopicStore) GetTopicList(topicID primitive.ObjectID) ([]models.TopicInfo, error) {

	//开启事务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//获取目标列表
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	cursor, err := collection.Find(ctx, bson.M{"_id": topicID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var topics []models.TopicInfo
	err = cursor.All(ctx, &topics)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

// ValidateTopicExistence 验证目标是否存在
//
// 参数：
//   - topicID: 目标ID
//
// 返回值：
//   - bool: 目标是否存在
//   - error: 错误信息
func (store *TopicStore) ValidateTopicExistence(topicID primitive.ObjectID) (bool, error) {
	// 开启上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 定义查询条件
	filter := bson.M{"_id": topicID}

	// 查询数据库
	result := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics").FindOne(ctx, filter)
	fmt.Println(result)
	if result.Err() == mongo.ErrNoDocuments {
		// 如果找不到对应的记录，则说明topic不存在
		return false, nil
	} else if result.Err() != nil {
		// 如果查询过程中出现其他错误，则返回错误信息
		return false, result.Err()
	}

	// 找到了记录，topic存在
	return true, nil
}

// LikeTopic 点赞目标
//
// 参数：
//   - topicID: 目标ID
//
// 返回值：
//   - error: 错误信息
func (store *TopicStore) LikeTopic(topicID primitive.ObjectID) error {
	// 开启事务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 点赞目标
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"like": 1}})
	if err != nil {
		return err
	}

	return nil
}

// CancelLikeTopic 取消点赞目标
//
// 参数：
//   - topicID: 目标ID
//
// 返回值：
//   - error: 错误信息
func (store *TopicStore) CancelLikeTopic(topicID primitive.ObjectID) error {
	// 开启事务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 取消点赞目标
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"like": -1}})
	if err != nil {
		return err
	}

	return nil
}

// DislikeTopic 点踩目标
//
// 参数：
//   - topicID: 目标ID
//
// 返回值：
//   - error: 错误信息
func (store *TopicStore) DislikeTopic(topicID primitive.ObjectID) error {
	// 开启事务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 点踩目标
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"dislike": 1}})
	if err != nil {
		return err
	}

	return nil
}

// GetHotTopics 获取热门话题
//
// 参数：
//   - limit: 限制返回的话题数量
//
// 返回值：
//   - []models.TopicInfo: 热门话题列表
//   - error: 错误信息
func (store *TopicStore) GetHotTopics(limit int) ([]models.TopicInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	opts := options.Find().SetSort(bson.D{{Key: "like", Value: -1}}).SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var topics []models.TopicInfo
	if err := cursor.All(ctx, &topics); err != nil {
		return nil, err
	}

	return topics, nil
}
