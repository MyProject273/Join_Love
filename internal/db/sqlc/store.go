package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store is defines all functions to execute db queries and transaction
type Store interface {
	Querier
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
