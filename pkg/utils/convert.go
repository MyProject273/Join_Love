package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PgDateToProtoTimestamp(pgDate pgtype.Date) *timestamppb.Timestamp {
	t, err := pgDate.Value()
	if err != nil {
		return nil
	}

	// Nếu null thì t sẽ là nil
	tm, ok := t.(time.Time)
	if !ok {
		return nil
	}

	return timestamppb.New(tm)
}

func PgTimestampToProto(ts pgtype.Timestamp) *timestamppb.Timestamp {
	t, err := ts.Value()
	if err != nil {
		return nil
	}

	tm, ok := t.(time.Time)
	if !ok {
		return nil
	}

	return timestamppb.New(tm)
}

func PgUUIDToString(uuid pgtype.UUID) (string, error) {
	val, err := uuid.Value()
	if err != nil {
		return "", err
	}
	str, ok := val.(string)
	if !ok {
		if b, ok := val.([]byte); ok {
			return string(b), nil
		}
		return "", fmt.Errorf("unexpected UUID type: %T", val)
	}
	return str, nil
}
