package persistence

import (
	"context"

	"gorm.io/gorm"
)

// TxManager 事务管理器
type TxManager struct {
	db *gorm.DB
}

// NewTxManager 创建事务管理器
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

// WithTransaction 在事务中执行操作
func (tm *TxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务放入context中
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

// GetDB 获取数据库实例（支持从context中获取事务）
func GetDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return db.WithContext(ctx)
}
