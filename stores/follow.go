/*
Package stores - NekoBlog backend server data access objects.
This file is for follow storage accessing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
- CBofJOU<2023122312@jou.edu.cn>
*/
package stores

import (
	"context"
	"errors"
	"time"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// Comment 评论信息数据库
type FollowStore struct {
	db    *gorm.DB
	mongo *mongo.Client
}

// NewFollowStore 返回一个新的用户存储实例。
// 返回：
//   - *FollowStore: 返回一个指向新的用户存储实例的指针。
func (factory *Factory) NewFollowStore() *FollowStore {
	return &FollowStore{
		factory.db,
		factory.mongo,
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
func (store *FollowStore) FollowUser(uid, followedID uint64) error {
	// 构造查询条件
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "followed_id", Value: followedID},
	}
	// 构造更新内容
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "followed_at", Value: time.Now()},
		}},
	}
	
	// 更新用户关注记录
	followRecordCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION)
	_, err := followRecordCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	// 重复关注
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("user has liked this followed")
	}
	return err
}

// CancelFollowUser 取消关注用户
//
// 参数：
//   - uid：用户ID
//   - followedID：被关注用户ID
//
// 返回值：
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *FollowStore) CancelFollowUser(uid, followedID uint64) error {
	// 构造查询条件
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "followed_id", Value: followedID},
	}
	
	followRecordCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION)
	_, err := followRecordCollection.DeleteOne(context.Background(), filter)
	if err == mongo.ErrNoDocuments {
		return errors.New("user has not liked this followed")
	}
	return err
}

// GetFollowStatus 获取关注状态
//
// 参数：
//   - uid：用户ID
//   - followedID：被关注用户ID
//
// 返回值：
//   - bool：关注状态
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *FollowStore) GetFollowStatus(uid, followedID uint64) (bool, error) {
	filter := bson.M{
		"uid":        uid,
		"followed_id": followedID,
	}
	count, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetFollowList 获取关注列表
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - []models.FollowInfo：关注列表
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *FollowStore) GetFollowList(userID uint64) ([]models.FollowInfo, error) {
	var followInfos []models.FollowInfo
	filter := bson.M{
		"uid": userID,
	}
	cur, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(context.Background(), &followInfos); err != nil {
		return nil, err
	}
	return followInfos, nil
}

// GetFollowersByUID 获取关注人数
//
// 参数：
//   - uid：用户ID
// 
// 返回值：
//   - int64：关注人数
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *FollowStore) GetFollowedsByUID(uid uint64) (int64, error) {
    filter := bson.M{
		"uid": uid,
	}
	return store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).CountDocuments(context.Background(), filter)
}

// GetFollowerList 获取粉丝列表
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - []models.FollowInfo：关注列表
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *FollowStore) GetFollowerList(userID uint64) ([]models.FollowInfo, error) {
	var followInfos []models.FollowInfo
	filter := bson.M{
		"followed_id": userID,
	}
	cur, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).Find(context.Background(), filter)
	if err != nil {
	    return nil, err
	}
	if err := cur.All(context.Background(), &followInfos); err != nil {
	    return nil, err
	}
	return followInfos, nil
}

// GetFollowersByUID 获取粉丝人数
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - int64：粉丝人数
//   - error：如果发生错误，返回相应错误信息；否则返回 nil
func (store *FollowStore) GetFollowersByUID(uid uint64) (int64, error) {
    filter := bson.M{
		"followed_id": uid,
	}
	return store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).CountDocuments(context.Background(), filter)
}