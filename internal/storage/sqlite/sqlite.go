package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
    				mail TEXT PRIMARY KEY NOT NULL,
    				password TEXT NOT NULL,
   					registration_date DATE NOT NULL);
				CREATE INDEX IF NOT EXISTS idx_mail ON users(mail);
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating the users table: %v", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(
		`CREATE TABLE IF NOT EXISTS tasks (
    				task_id uid PRIMARY KEY,
    				title TEXT NOT NULL,
    				url TEXT,
   					task_add_date DATE NOT NULL);
				CREATE INDEX IF NOT EXISTS idx_task_id ON tasks(task_id);
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating the tasks table: %v", err)
	}
	_, err = stmt.Exec()

	return &Storage{db: db}, nil
}

func (s *Storage) AddUser(mail, password string) error {
	const op = "storage.sqlite.addUser"
	registrationDate := time.Now()
	stmt, err := s.db.Prepare("INSERT INTO users(mail, password, registration_date) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err = stmt.Exec(mail, hashedPassword, registrationDate)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CheckPassword(mail, password string) error {
	const op = "storage.sqlite.checkPassword"
	stmt, err := s.db.Prepare("SELECT password FROM users WHERE mail=?")

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRow(mail)
	var pass string
	err = row.Scan(&pass)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password! %s: %w", op, err)
	}
	return nil
}

func (s *Storage) UserExists(mail string) (string, error) {
	const op = "storage.sqlite.userExists"
	stmt, err := s.db.Prepare("SELECT mail FROM users WHERE mail = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRow(mail)
	var email string
	err = row.Scan(&email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	return email, nil
}

func (s *Storage) AddTask(title, url string) error {
	const op = "storage.sqlite.addTask"
	stmt, err := s.db.Prepare("INSERT INTO tasks(task_id, title, url, task_add_date) VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	taskDate := time.Now()
	taskId := uuid.New()
	_, err = stmt.Exec(taskId, title, url, taskDate)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
