package sqlite

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS users (
    				id SERIAL PRIMARY KEY,
    				mail TEXT NOT NULL,
   					registration_date DATE NOT NULL);
				CREATE INDEX IF NOT EXISTS idx_mail ON users(mail);
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating the users table: %v", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("error creating the users table: %v", err)
	}
	fmt.Println("Successfully connected!")

	return &Storage{db: db}, nil
}

func (s *Storage) AddUser(mail string) error {
	const op = "storage.sqlite.addUser"
	fmt.Println("AddUser", op, mail)
	registrationDate := time.Now()
	stmt, err := s.db.Prepare("INSERT INTO users(mail, registration_date) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(mail, registrationDate)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
