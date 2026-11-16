package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New(SQLitePath string) (*sql.DB, error) {
	return sql.Open("sqlite3", SQLitePath+"?_fk=on")
}
