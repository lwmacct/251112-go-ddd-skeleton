package model

import "gorm.io/gorm"

// AllModels 返回所有需要迁移的模型
func AllModels() []interface{} {
	return []interface{}{
		// User相关
		&User{},
		&TwoFactor{},
		&PersonalAccessToken{},
		&Session{},
		
		// Order相关
		&Order{},
		&OrderItem{},
		&Payment{},
		&Shipment{},
		&Invoice{},
	}
}

// AutoMigrate 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}

