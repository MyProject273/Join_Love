package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store is defines all functions to execute db queries and transaction
type Store interface {
	Querier
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	CreateUserByAdminTx(ctx context.Context, arg CreateUserByAdminTxParams) (CreateUserTxResult, error)
	UpdateUserTx(ctx context.Context, arg UpdateUserTxParams) (UpdateUserTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// Function NewStore create a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
