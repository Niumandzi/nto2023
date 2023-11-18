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
		`CREATE TABLE IF NOT EXISTS details (
			id INTEGER PRIMARY KEY,
			type_name VARCHAR(255) NOT NULL,
			category TEXT NOT NULL CHECK (category IN ('entertainment', 'enlightenment', 'education')),
		    UNIQUE(type_name, category)
		  );
	
		CREATE TABLE IF NOT EXISTS events (
		  	id INTEGER PRIMARY KEY,
		  	name VARCHAR(255) NOT NULL,
		  	description TEXT,
		  	date TEXT,
		  	details_id INT NOT NULL,
		  	FOREIGN KEY (details_id) REFERENCES details(id))
		  	`)
	if err != nil {
		return err
	}

	return nil
}
