package sqlite

import (
	"database/sql"
	"time"
)

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func toNullTime(t time.Time, layout string) sql.NullString {
	return sql.NullString{String: t.Format(layout), Valid: !t.IsZero()}
}

func toNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: i != 0}
}

func toNullInt64FromPointer(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i}
}

func fromNullableInt64ToPointer(i sql.NullInt64) *int64 {
	if !i.Valid {
		return nil
	}
	v := i.Int64
	return &v
}

func fromNullableTime(t sql.NullString, layout string) (time.Time, error) {
	if !t.Valid {
		return time.Time{}, nil
	}
	return time.Parse(layout, t.String)
}
