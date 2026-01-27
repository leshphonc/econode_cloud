package txm

import (
	"context"

	"gorm.io/gorm"
)

type TxManager struct{ db *gorm.DB }

func NewTxManager(db *gorm.DB) *TxManager { return &TxManager{db: db} }

func (m *TxManager) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

//type ctxKey string

//const ctxKeyTx ctxKey = "gorm_tx"

//func (m *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
//	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//		txCtx := WithTx(ctx, tx)
//		return fn(txCtx)
//	})
//}
//func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
//	return context.WithValue(ctx, ctxKeyTx, tx)
//}

//
//func DBFrom(ctx context.Context, fallback *gorm.DB) *gorm.DB {
//	if v := ctx.Value(ctxKeyTx); v != nil {
//		if tx, ok := v.(*gorm.DB); ok && tx != nil {
//			return tx
//		}
//	}
//	return fallback
//}
