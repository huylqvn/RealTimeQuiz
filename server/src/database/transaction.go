package database

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type ContextKey string

const ContextKeyTransaction ContextKey = "transaction"

func Transaction(
	ctx context.Context,
	db *gorm.DB,
	execute func(context.Context, *gorm.DB) error,
	options ...*sql.TxOptions,
) (err error) {
	// Check current ctx have transaction or not
	// If not, create new transaction
	// If yes, use current transaction, gorm Transaction() handle nested transaction
	var (
		tx *gorm.DB
	)
	if ctx.Value(ContextKeyTransaction) == nil {
		tx = db
	} else {
		tx = ctx.Value(ContextKeyTransaction).(*gorm.DB)
	}

	// Pass executeCtx to transaction execution
	return tx.Transaction(func(tx *gorm.DB) error {
		return execute(context.WithValue(ctx, ContextKeyTransaction, tx), tx)
	}, options...)
}
