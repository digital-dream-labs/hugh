package sql

import (
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewNullString converts a string to sql.NullString
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NewNullTimePB converts a timestamppb.Timestamp to sql.NullString
func NewNullTimePB(s *timestamppb.Timestamp) sql.NullTime {
	ts, err := ptypes.Timestamp(s)
	if err != nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  ts,
		Valid: true,
	}
}

// NewNullTime converts time.Time to sql.NullTime
func NewNullTime(s time.Time) sql.NullTime {
	if s.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  s,
		Valid: true,
	}
}
