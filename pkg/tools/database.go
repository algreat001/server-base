package tools

import (
	"database/sql"
)

type DataBase interface {
	Close() error
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Ping() error
	GetStats() string
}
