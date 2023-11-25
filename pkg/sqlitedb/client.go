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
		  	FOREIGN KEY (details_id) REFERENCES details(id));

		CREATE TABLE IF NOT EXISTS work_type (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		  );

		CREATE TABLE IF NOT EXISTS facility (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		  );

		CREATE TABLE IF NOT EXISTS application (
		  	id INTEGER PRIMARY KEY,
		  	description TEXT,
		  	created_at TEXT,
		  	due TEXT,
			status TEXT NOT NULL CHECK (status IN ('created', 'todo', 'done')),
		    work_type_id INT,
			event_id INT,
			facility_id INT,
		    FOREIGN KEY (work_type_id) REFERENCES work_type(id),
			FOREIGN KEY (event_id) REFERENCES events(id),
			FOREIGN KEY (facility_id) REFERENCES facility(id));
		  	`)
	if err != nil {
		return err
	}

	return nil
}
