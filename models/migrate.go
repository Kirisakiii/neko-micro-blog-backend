/*
Package models - NekoBlog backend server models migration.
This file is for database migration function.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import "gorm.io/gorm"

// Migrate 数据库迁移
//
// 参数：
//	- db *gorm.DB 数据库连接
//
// 返回值：
//	- error 错误
func Migrate(db *gorm.DB) error {
	var err error

	// User 相关
	if err = db.AutoMigrate(&UserInfo{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserAuthInfo{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserLoginLog{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserPostStatus{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserCommentStatus{}); err != nil {
		return err
	}

	// Post 相关
	if err = db.AutoMigrate(&PostInfo{}); err != nil {
		return err
	}
	
	// Comment 相关
	if err = db.AutoMigrate(&CommentInfo{}); err != nil {
		return err
	}

	// Reply 相关
	if err = db.AutoMigrate(&ReplyInfo{}); err != nil {
		return err
	}

	return nil
}
