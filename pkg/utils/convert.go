package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Strict version with error
func PgDateToProtoTimestamp(pgDate pgtype.Date) (*timestamppb.Timestamp, error) {
	t, err := pgDate.Value()
	if err != nil {
		return nil, fmt.Errorf("PgDateToProtoTimestamp: cannot get value: %w", err)
	}
	tm, ok := t.(time.Time)
	if !ok {
		return nil, fmt.Errorf("PgDateToProtoTimestamp: unexpected type %T", t)
	}
	return timestamppb.New(tm), nil
}

// Safe version: ignore error, return nil on fail
func PgDateToProtoTimestampSafe(pgDate pgtype.Date) *timestamppb.Timestamp {
	ts, err := PgDateToProtoTimestamp(pgDate)
	if err != nil {
		return nil
	}
	return ts
}

func PgTimestampToProto(ts pgtype.Timestamp) (*timestamppb.Timestamp, error) {
	t, err := ts.Value()
	if err != nil {
		return nil, fmt.Errorf("PgTimestampToProto: cannot get value: %w", err)
	}
	tm, ok := t.(time.Time)
	if !ok {
		return nil, fmt.Errorf("PgTimestampToProto: unexpected type %T", t)
	}
	return timestamppb.New(tm), nil
}

func PgTimestampToProtoSafe(ts pgtype.Timestamp) *timestamppb.Timestamp {
	tm, err := PgTimestampToProto(ts)
	if err != nil {
		return nil
	}
	return tm
}

func PgUUIDToString(uuid pgtype.UUID) (string, error) {
	val, err := uuid.Value()
	if err != nil {
		return "", fmt.Errorf("PgUUIDToString: cannot get value: %w", err)
	}
	switch v := val.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("PgUUIDToString: unexpected type %T", val)
	}
}

func PgUUIDToStringSafe(uuid pgtype.UUID) string {
	s, err := PgUUIDToString(uuid)
	if err != nil {
		return ""
	}
	return s
}
