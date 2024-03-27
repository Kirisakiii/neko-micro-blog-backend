/*
Package stores - NekoBlog backend server data access objects.
This file is for factory of storages.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package stores

import (
	search "github.com/Kirisakiii/neko-micro-blog-backend/proto"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Factory struct {
	db            *gorm.DB
	rds           *redis.Client
	mongo         *mongo.Client
	searchService search.SearchEngineClient
}

func NewFactory(db *gorm.DB, redisClient *redis.Client, mongoClient *mongo.Client, searchService search.SearchEngineClient) *Factory {
	return &Factory{
		db:            db,
		rds:           redisClient,
		mongo:         mongoClient,
		searchService: searchService,
	}
}
