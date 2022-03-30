package database

import "gorm.io/gorm"

type IDatabase interface {
	// NewDatabase 获取 DB 实例
	NewDatabase(name ...string) *gorm.DB
	// AutoMigrate 注册数据模型
	AutoMigrate(tables ...interface{}) error
}
