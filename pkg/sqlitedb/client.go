package sqlitedb

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewClient(driverName string, filePath string) (*sql.DB, error) {
	db, err := sql.Open(driverName, filePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	return nil
}
