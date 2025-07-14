package db

import (
	"errors"

	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: consts.UniqueViolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func ErrorConstraint(err error) string {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return ""
	}
	return pgErr.ConstraintName
}
