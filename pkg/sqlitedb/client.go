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
		`
		PRAGMA foreign_keys = ON;
        CREATE TABLE IF NOT EXISTS details (
			id INTEGER PRIMARY KEY,
			type_name VARCHAR(255) NOT NULL UNIQUE,
			category TEXT NOT NULL CHECK (category IN ('entertainment', 'enlightenment', 'education')),
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
		    UNIQUE(type_name, category)
		);
	
		CREATE TABLE IF NOT EXISTS events (
		  	id INTEGER PRIMARY KEY,
		  	name VARCHAR(255) NOT NULL UNIQUE,
		  	description TEXT NOT NULL ,
		  	date TEXT NOT NULL,
		  	details_id INT NOT NULL,
		  	is_active BOOLEAN NOT NULL DEFAULT TRUE,
		  	FOREIGN KEY (details_id) REFERENCES details(id)
		);

		CREATE TABLE IF NOT EXISTS work_type (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			is_active BOOLEAN NOT NULL DEFAULT TRUE        
		);

		CREATE TABLE IF NOT EXISTS facility (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
		    have_parts BOOLEAN NOT NULL,
		    is_active BOOLEAN NOT NULL DEFAULT TRUE
		);

		CREATE TABLE IF NOT EXISTS part (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			facility_id INT NOT NULL,
		    is_active BOOLEAN NOT NULL DEFAULT TRUE,
		    FOREIGN KEY (facility_id) REFERENCES facility(id) ON DELETE CASCADE
		);
    	
		CREATE TABLE IF NOT EXISTS application (
		  	id INTEGER PRIMARY KEY,
		  	description TEXT NOT NULL,
		  	created_at TEXT NOT NULL,
		  	due TEXT NOT NULL,
			status TEXT NOT NULL CHECK (status IN ('created', 'todo', 'done')),
		    work_type_id INT NOT NULL,
			event_id INT NOT NULL,
			facility_id INT NOT NULL,
		    FOREIGN KEY (work_type_id) REFERENCES work_type(id),
			FOREIGN KEY (event_id) REFERENCES events(id),
			FOREIGN KEY (facility_id) REFERENCES facility(id)
		);

		CREATE TABLE IF NOT EXISTS booking (
			id INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			create_date TEXT NOT NULL,
			start_date TEXT NOT NULL,
			start_time TEXT NOT NULL,
			end_date TEXT NOT NULL,
			end_time TEXT NOT NULL,
			event_id INT NOT NULL,
			facility_id INT NOT NULL, 
			FOREIGN KEY (facility_id) REFERENCES facility(id)
        );
		
		CREATE TABLE IF NOT EXISTS booking_part (
			booking_id INT NOT NULL,
			part_id INT NOT NULL,
			FOREIGN KEY (booking_id) REFERENCES booking(id) ON DELETE CASCADE,
			FOREIGN KEY (part_id) REFERENCES part(id)
		);
		
		CREATE TABLE IF NOT EXISTS mug_type (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			is_active BOOLEAN NOT NULL DEFAULT TRUE        
		);
				
		CREATE TABLE IF NOT EXISTS teacher (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			is_active BOOLEAN NOT NULL DEFAULT TRUE        
		);

		CREATE TABLE IF NOT EXISTS schedule (
		    id INTEGER PRIMARY KEY,
			day TEXT NOT NULL CHECK (day IN ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday')),
			start_time TEXT NOT NULL,
			end_time TEXT NOT NULL,
			registration_id INT NOT NULL,
			FOREIGN KEY (registration_id) REFERENCES registration(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS registration_part (
			registration_id INT NOT NULL,
			part_id INT NOT NULL,
			FOREIGN KEY (registration_id) REFERENCES registration(id) ON DELETE CASCADE,
			FOREIGN KEY (part_id) REFERENCES part(id)
		);

		CREATE TABLE IF NOT EXISTS registration (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE, 
			start_date TEXT NOT NULL,
			number_of_days INT NOT NULL CHECK (number_of_days IN (1, 2, 3)),
			facility_id INT NOT NULL,
			mug_type_id INT NOT NULL,
			teacher_id INT NOT NULL,
			FOREIGN KEY (facility_id) REFERENCES facility(id),
			FOREIGN KEY (mug_type_id) REFERENCES mug_type(id),
			FOREIGN KEY (teacher_id) REFERENCES teacher(id)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
