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
	userLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FOLLOW_RECORD_COLLECTION)
	_, err := userLikeCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
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
	
	userFollowCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FOLLOW_RECORD_COLLECTION)
	_, err := userFollowCollection.DeleteOne(context.Background(), filter)
	if mongo.ErrNoDocuments == err {
		return errors.New("user has not liked this followed")
	}
	return err
}

// GetFollowList 获取关注列表
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - 成功则返回关注列表
//   - 失败返回nil
func (store *FollowStore) GetFollowList(userID uint64) ([]models.FollowInfo, error) {
	var followInfos []models.FollowInfo
	filter := bson.M{
		"uid": userID,
	}
	cur, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FOLLOW_RECORD_COLLECTION).Find(context.Background(), filter)
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
//   - 成功则返回关注人数
func (store *FollowStore) GetFollowedsByUID(uid uint64) (int, error) {
    filter := bson.M{
		"uid": uid,
	}
	cur, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FOLLOW_RECORD_COLLECTION).Find(context.Background(), filter)
    if err != nil {
        return 0, err
    }

	count := 0
    for cur.Next(context.Background()) {
        var followInfo models.FollowInfo
        if err := cur.Decode(&followInfo); err != nil {
            return 0, err
        }
        count++
    }
    if err := cur.Err(); err != nil {
        return 0, err
    }

    return count, nil
}

// GetFollowerList 获取粉丝列表
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - 成功则返回关注列表
//   - 失败返回nil
func (store *FollowStore) GetFollowerList(userID uint64) ([]models.FollowInfo, error) {
	var followInfos []models.FollowInfo
	filter := bson.M{
		"followed_id": userID,
	}
	cur, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FOLLOW_RECORD_COLLECTION).Find(context.Background(), filter)
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
//   - 成功则返回粉丝人数
//   - 失败返回nil
func (store *FollowStore) GetFollowersByUID(uid uint64) (int, error) {
    filter := bson.M{
		"followed_id": uid,
	}
	cur, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FOLLOW_RECORD_COLLECTION).Find(context.Background(), filter)
    if err != nil {
        return 0, err
    }

	count := 0
    for cur.Next(context.Background()) {
        var followInfo models.FollowInfo
        if err := cur.Decode(&followInfo); err != nil {
            return 0, err
        }
        count++
    }
    if err := cur.Err(); err != nil {
        return 0, err
    }

    return count, nil
}