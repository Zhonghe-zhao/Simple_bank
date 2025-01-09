package db

import (
	"context"
	"fmt"
)

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	//tx是sqlx的Tx对象，用来与数据库交互
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	//用来于数据库交互
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
