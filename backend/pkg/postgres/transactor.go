package postgres

import (
	"context"
)

type contextKey struct{}

var txKey = contextKey{}

type Transactor struct {
	db *DB
}

func NewTransactor(db *DB) *Transactor {
	return &Transactor{db: db}
}

func (t *Transactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.beginTx(ctx)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey, tx)
	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
