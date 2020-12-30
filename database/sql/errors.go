package sql

import (
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ParseSQLError turns mysql errors into codes
func ParseSQLError(err error) error {
	m, ok := err.(*mysql.MySQLError)
	if !ok {
		return parseText(err)
	}
	return parseMYSQLError(m)
}

func parseText(err error) error {
	switch err.Error() {
	case "sql: no rows in result set":
		return status.Error(codes.NotFound, "NOT_FOUND")
	default:
		return status.Error(codes.Internal, "INTERNAL_ERROR")
	}
}

func parseMYSQLError(m *mysql.MySQLError) error {
	switch m.Number {
	case 1062:
		return status.Error(codes.AlreadyExists, "ALREADY_EXISTS")
	default:
		return status.Error(codes.Internal, "INTERNAL_ERROR")
	}
}
