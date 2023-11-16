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
		`CREATE TABLE IF NOT EXISTS event_types (
			id SERIAL PRIMARY KEY,
			type_name VARCHAR(255) NOT NULL
		  );
	
		CREATE TABLE IF NOT EXISTS categories_types (
			id SERIAL PRIMARY KEY,
			type_id INT NOT NULL,
			category TEXT CHECK (category IN ('entertainment', 'enlightenment', 'education')),
			FOREIGN KEY (type_id) REFERENCES event_types(id)
		  );
	
		CREATE TABLE IF NOT EXISTS events (
		  id SERIAL PRIMARY KEY,
		  name VARCHAR(255) NOT NULL,
		  description TEXT,
		  details_id INT NOT NULL,
		  FOREIGN KEY (details_id) REFERENCES categories_types(id)
		  )`)
	if err != nil {
		return err
	}

	return nil
}
