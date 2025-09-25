package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
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

func StringToPgUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("StringToPgUUID: invalid UUID format: %w", err)
	}

	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(parsed.String()); err != nil {
		return pgtype.UUID{}, fmt.Errorf("StringToPgUUID: cannot scan to pgtype.UUID: %w", err)
	}
	return pgUUID, nil
}

func StringToPgUUIDSafe(s string) pgtype.UUID {
	uuid, err := StringToPgUUID(s)
	if err != nil {
		return pgtype.UUID{}
	}
	return uuid
}

func UUIDToPgUUID(u uuid.UUID) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(u.String()); err != nil {
		return pgtype.UUID{}, fmt.Errorf("UUIDToPgUUID: cannot scan to pgtype.UUID: %w", err)
	}
	return pgUUID, nil
}

// Safe version: ignore error, return zero value on fail
func UUIDToPgUUIDSafe(u uuid.UUID) pgtype.UUID {
	pgUUID, err := UUIDToPgUUID(u)
	if err != nil {
		return pgtype.UUID{}
	}
	return pgUUID
}

// Strict version: string -> int64
func StringToInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("StringToInt64: cannot parse %q: %w", s, err)
	}
	return val, nil
}

// Safe version: string -> int64, return 0 on error
func StringToInt64Safe(s string) int64 {
	val, err := StringToInt64(s)
	if err != nil {
		return 0
	}
	return val
}

// Strict version: pgtype.Date -> "YYYY-MM-DD"
func PgDateToString(pgDate pgtype.Date) (string, error) {
	t, err := pgDate.Value()
	if err != nil {
		return "", fmt.Errorf("PgDateToString: cannot get value: %w", err)
	}
	tm, ok := t.(time.Time)
	if !ok {
		return "", fmt.Errorf("PgDateToString: unexpected type %T", t)
	}
	return tm.Format("2006-01-02"), nil
}

// Safe version: ignore error, return "" on fail
func PgDateToStringSafe(pgDate pgtype.Date) string {
	s, err := PgDateToString(pgDate)
	if err != nil {
		return ""
	}
	return s
}

// Strict version: string -> pgtype.Date
func StringToPgDate(s string) (pgtype.Date, error) {
	if s == "" {
		return pgtype.Date{}, fmt.Errorf("StringToPgDate: empty string")
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("StringToPgDate: cannot parse %q: %w", s, err)
	}

	return pgtype.Date{
		Time: t,
	}, nil
}

// Safe version: ignore error, return zero value on fail
func StringToPgDateSafe(s string) pgtype.Date {
	d, err := StringToPgDate(s)
	if err != nil {
		return pgtype.Date{} // Status == Null
	}
	return d
}
