package utils

import (
	"database/sql"
	"time"
)

func TimePtrToSQLNullInt64(t *time.Time) sql.NullInt64 {
	if t == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{
		Int64: t.UnixMilli(),
		Valid: true,
	}
}

func TimePtrToInt64Ptr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	var i int64
	i = t.UnixMilli()
	return &i
}

func SQLNullInt64ToTimePtr(ni64 sql.NullInt64) *time.Time {
	if !ni64.Valid {
		return nil
	} else {
		t := time.UnixMilli(ni64.Int64)
		return &t
	}
}

func NowToTimePtr() *time.Time {
	t := time.Now()
	return &t
}

func SQLInt64Ptr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

func SQLStrPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func NullableInt64ToSQL(i *int64) interface{} {
	if i != nil {
		return i
	}
	return nil
}

func NullableStringToSQL(s *string) interface{} {
	if s != nil {
		return s
	}
	return nil
}
