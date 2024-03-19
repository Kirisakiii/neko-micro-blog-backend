/*
Package stores - NekoBlog backend server data access objects.
This file is for factory of storages.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package stores

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewFactory(db *gorm.DB, redisClient *redis.Client) *Factory {
	return &Factory{db: db, rds: redisClient}
}
