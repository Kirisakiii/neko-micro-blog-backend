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
	"errors"
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

// GetTopicList 获取话题列表
//
// 参数：
//   - from (primitive.ObjectID)：起始 ID。
//   - length (uint64)：长度。
//
// 返回值：
//   - []models.TopicInfo：话题列表。
//   - error：如果获取话题列表时发生错误，则返回一个错误。
func (store *TopicStore) GetTopicList() ([]models.TopicInfo, error) {
	ctx := context.Background()

	// 查询数据库
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
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
func (store *TopicStore) CreateTopic(userID uint64, name string, description string) (primitive.ObjectID, error) {
	ctx := context.Background()
	// 插入topic信息
	topic := models.TopicInfo{
		Name:           name,
		Description:    description,
		CreatorID:      userID,
		Like:           0,
		DisLike:        0,
		CreatedAt:      time.Now(),
	}

	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	result, err := collection.InsertOne(ctx, &topic)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// ValidateTopicExistenceByName 验证目标是否存在
//
// 参数：
//   - name: 目标名
//
// 返回值：
//   - bool: 目标是否存在
//   - error: 错误信息
func (store *TopicStore) ValidateTopicExistenceByName(name string) (bool, error) {
	ctx := context.Background()

	filter := bson.M{
		"name": name,
	}

	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	result := collection.FindOne(ctx, filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return false, nil
	}
	if result.Err() != nil {
		return false, result.Err()
	}

	return true, nil
}

// DeleteTopic 删除目标信息
//
// 参数：
//   - topicID: 目标ID
//
// 返回值：
//   - error: 错误信息
func (store *TopicStore) DeleteTopic(topicID primitive.ObjectID) error {
	ctx := context.Background()

	//删除目标
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": topicID})
	if result.DeletedCount == 0 {
		return errors.New("no such topic")
	}
	if err != nil {
		return err
	}

	return nil
}

// GetTopicDetail 获取目标
//
// 参数：
//   - tagID: 话题ID
//
// 返回值：
//   - models.TagInfo: 目标
//   - error: 错误信息
func (store *TopicStore) GetTopicDetail(topicID primitive.ObjectID) (models.TopicInfo, error) {
	var ctx = context.Background()

	filter := bson.M{
		"_id": topicID,
	}
	ret := models.TopicInfo{}

	//获取目标列表
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	result := collection.FindOne(ctx, filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return ret, errors.New("no such topic")
	}
	if result.Err() != nil {
		return ret, result.Err()
	}
	if err := result.Decode(&ret); err != nil {
		return ret, err
	}

	return ret, nil
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
	ctx := context.Background()

	// 定义查询条件
	filter := bson.M{"_id": topicID}

	// 查询数据库
	result := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics").FindOne(ctx, filter)
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
func (store *TopicStore) LikeTopic(topicID primitive.ObjectID, userID uint64) error {
	ctx := context.Background()

	// 查询是否已经点赞
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_likes")
	count, err := collection.CountDocuments(ctx, bson.M{"topic_id": topicID, "user_id": userID})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("already liked")
	}

	// 点赞目标
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err = collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"like": 1}})
	if err != nil {
		return err
	}

	// 记录点赞记录
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_likes")
	_, err = collection.InsertOne(ctx, bson.M{"topic_id": topicID, "created_at": time.Now(), "user_id": userID})
	if err != nil {
		return err
	}

	// 取消踩
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_dislikes")
	result, err := collection.DeleteOne(ctx, bson.M{"topic_id": topicID, "user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount > 0 {
		collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
		_, err = collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"dislike": -1}})
		if err != nil {
			return err
		}
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
func (store *TopicStore) CancelLikeTopic(topicID primitive.ObjectID, userID uint64) error {
	ctx := context.Background()

	// 删除点赞记录
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_likes")
	result, err := collection.DeleteOne(ctx, bson.M{"topic_id": topicID, "user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no such like record")
	}

	// 取消点赞目标
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err = collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"like": -1}})
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
func (store *TopicStore) DislikeTopic(topicID primitive.ObjectID, userID uint64) error {
	ctx := context.Background()

	// 查询是否已经点踩
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_dislikes")
	count, err := collection.CountDocuments(ctx, bson.M{"topic_id": topicID, "user_id": userID})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("already disliked")
	}

	// 点踩目标
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err = collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"dislike": 1}})
	if err != nil {
		return err
	}

	// 记录点踩记录
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_dislikes")
	_, err = collection.InsertOne(ctx, bson.M{"topic_id": topicID, "created_at": time.Now(), "user_id": userID})
	if err != nil {
		return err
	}

	// 取消点赞
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_likes")
	result, err := collection.DeleteOne(ctx, bson.M{"topic_id": topicID, "user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount > 0 {
		collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
		_, err = collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"like": -1}})
		if err != nil {
			return err
		}
	}

	return nil
}

// CancelDislikeTopic 取消点踩目标
//
// 参数：
//   - topicID: 目标ID
//   - userID: 用户ID
//
// 返回值：
//   - error: 错误信息
func (store *TopicStore) CancelDislikeTopic(topicID primitive.ObjectID, userID uint64) error {
	ctx := context.Background()

	// 删除点踩记录
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_dislikes")
	result, err := collection.DeleteOne(ctx, bson.M{"topic_id": topicID, "user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no such dislike record")
	}

	// 取消点踩目标
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topics")
	_, err = collection.UpdateOne(ctx, bson.M{"_id": topicID}, bson.M{"$inc": bson.M{"dislike": -1}})
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

// GetUserTopicStatus 获取用户话题状态
//
// 参数：
//   - topicID: 话题ID
//   - userID: 用户ID
//
// 返回值：
//   - bool: 用户是否点赞
//   - bool: 用户是否点踩
//   - error: 错误信息
func (store *TopicStore) GetUserTopicStatus(topicID primitive.ObjectID, userID uint64) (bool, bool, error) {
	ctx := context.Background()

	// 查询是否点赞
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_likes")
	count, err := collection.CountDocuments(ctx, bson.M{"topic_id": topicID, "user_id": userID})
	if err != nil {
		return false, false, err
	}
	liked := count > 0

	// 查询是否点踩
	collection = store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_dislikes")
	count, err = collection.CountDocuments(ctx, bson.M{"topic_id": topicID, "user_id": userID})
	if err != nil {
		return false, false, err
	}
	disliked := count > 0

	return liked, disliked, nil
}

// GetBanner 获取话题横幅图
//
// 参数：
//   - topicID: 话题ID
//
// 返回值：
//   - string: 横幅图URL
//   - error: 错误信息
func (store *TopicStore) GetBanner(topicID primitive.ObjectID) (string, error) {
	ctx := context.Background()

	// 查询话题信息
	collection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection("topic_banners")
	result := collection.FindOne(ctx, bson.M{"topic_id": topicID})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return "default.webp", nil
	}
	if result.Err() != nil {
		return "", result.Err()
	}

	// 解码横幅图URL
	var banner struct {
		ResourceName string `bson:"resource_name"`
	}
	if err := result.Decode(&banner); err != nil {
		return "", err
	}

	return banner.ResourceName, nil
}
