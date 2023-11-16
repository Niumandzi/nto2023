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
	_, err := db.Exec(
		`CREATE TABLE contact (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			contact_type TEXT CHECK (contact_type IN ('worker', 'private_client', 'legal_client')),
			name TEXT,
			number TEXT,
			email TEXT
	   );

		CREATE TABLE IF NOT EXISTS hotel (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name_test VARCHAR(200),
			location_id INTEGER,
			number VARCHAR,
			worker_id INTEGER,
			description VARCHAR(500),
			FOREIGN KEY (worker_id) REFERENCES contact (id)
		)`)

	if err != nil {
		return err
	}

	return nil
}
