package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 能够执行数据库的所有功能 以及混合功能
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	//CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	//VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// SQLStore provides all functions to execute SQL queries and transactions
