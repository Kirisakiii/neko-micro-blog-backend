/*
Package models - NekoBlog backend server database models
This file is for user related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// UserInfo 用户信息模型
type UserInfo struct {
	gorm.Model            // 基本模型
	UserName   string     `gorm:"unique;column:username"`             // 用户名
	NickName   *string    `gorm:"column:nickname"`                    // 昵称
	Avatar     string     `gorm:"default:vanilla.webp;column:avatar"` // 头像
	Birth      *time.Time `gorm:"column:birth"`                       // 生日
	Gender     *string    `gorm:"column:gender"`                      // 性别
	Authority  uint64     `gorm:"default:0;column:authority"`         // 权限等级
	Level      uint64     `gorm:"default:1;column:level"`             // 等级
}

// UserAuthInfo 用户认证信息模型
type UserAuthInfo struct {
	gorm.Model          // 基本模型
	UID          uint64 `gorm:"unique;column:uid"`      // 用户ID
	UserName     string `gorm:"unique;column:username"` // 用户名
	Salt         string `gorm:"column:salt"`            // 盐
	PasswordHash string `gorm:"column:psw_hash"`        // 密码哈希值
}

// UserLoginLog 用户登录日志模型
type UserLoginLog struct {
	gorm.Model            // 基本模型
	UID         uint64    `gorm:"column:uid"`                         // 用户ID
	LoginTime   time.Time `gorm:"column:login_time"`                  // 登录时间
	LoginIP     string    `gorm:"column:login_ip"`                    // 登录IP
	IsSucceed   bool      `gorm:"column:is_succeed"`                  // 登录是否成功
	IfChecked   bool      `gorm:"default:false;column:if_checked"`    // 是否已经以被发送到客户端
	Reason      string    `gorm:"column:reason"`                      // 原因
	Device      string    `gorm:"default:unknown;column:device"`      // 登录时登陆的设备 如：Windows iOS Android
	Application string    `gorm:"default:unknown;column:application"` // 登录时使用的应用 如 Chrome 236.12
	BearerToken string    `gorm:"column:bearer_token"`                // 此次登录获取到的令牌
}

// UserAvaliableToken 用户可用Token模型
type UserAvaliableToken struct {
	gorm.Model           // 基本模型
	UID        uint64    `gorm:"column:uid"`          // 用户ID
	Username   string    `gorm:"column:username"`     // 用户名
	Token      string    `gorm:"unique;column:token"` // Token
	ExpireTime time.Time `gorm:"column:expire_time"`  // 过期时间
}

// UserLikedRecord 用户点赞记录
type UserLikedRecord struct {
	gorm.Model
	UID          uint64        `gorm:"column:uid"`                         // 用户ID
	LikedPost    pq.Int64Array `gorm:"column:liked_post;type:bigint[]"`    // 点赞过的博文
	LikedComment pq.Int64Array `gorm:"column:liked_comment;type:bigint[]"` // 点赞过的评论
	LikedReply   pq.Int64Array `gorm:"column:liked_reply;type:bigint[]"`   // 点赞过的回复
}

// UserDislikeRecord 用户点踩记录
type UserDislikeRecord struct {
	gorm.Model
	UID            uint64        `gorm:"column:uid"`                           // 用户ID
	DislikeComment pq.Int64Array `gorm:"column:dislike_comment;type:bigint[]"` // 点踩过的评论
	DislikeReply   pq.Int64Array `gorm:"column:dislike_reply;type:bigint[]"`   // 点踩过的回复
}

// UserFavouriteRecord 用户收藏记录
type UserFavouriteRecord struct {
	gorm.Model
	UID       uint64        `gorm:"column:uid"`                     // 用户ID
	Favourite pq.Int64Array `gorm:"column:favourite;type:bigint[]"` // 收藏过的博文
}
