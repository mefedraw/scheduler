package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	_ "modernc.org/sqlite"
	"scheduler/internal/dto"
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

	//  users
	stmt, err := db.Prepare(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            mail TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            registration_date DATE NOT NULL
        )`)
	if err != nil {
		return nil, fmt.Errorf("error preparing users table: %v", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// tasks
	stmt, err = db.Prepare(`
        CREATE TABLE IF NOT EXISTS tasks (
            task_id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            url TEXT,
            creation_timestamp DATE NOT NULL
        )`)
	if err != nil {
		return nil, fmt.Errorf("error preparing tasks table: %v", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// users_tasks
	stmt, err = db.Prepare(`
        CREATE TABLE IF NOT EXISTS users_tasks (
            task_id TEXT PRIMARY KEY REFERENCES tasks(task_id),
            user_id INTEGER REFERENCES users(id)
        )`)
	if err != nil {
		return nil, fmt.Errorf("error preparing users_tasks table: %v", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddUser(mail, password string, registrationDate time.Time) (int, error) {
	const op = "storage.sqlite.addUser"
	stmt, err := s.db.Prepare("INSERT INTO users(mail, password, registration_date) VALUES(?, ?, ?)")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.Exec(mail, password, registrationDate)
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: failed to get last insert ID: %w", op, err)
	}

	return int(id), nil
}

func (s *Storage) GetPassword(mail string) (string, error) {
	const op = "storage.sqlite.checkPassword"
	stmt, err := s.db.Prepare("SELECT password FROM users WHERE mail=?")

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRow(mail)
	var pass string
	err = row.Scan(&pass)

	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return pass, nil
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

func (s *Storage) AddTask(taskId uuid.UUID, creatorId uint, title, url string, creationTimestamp time.Time) error {
	const op = "storage.sqlite.addTask"
	stmt, err := s.db.Prepare("INSERT INTO tasks(task_id, title, url, creation_timestamp) VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(taskId, title, url, creationTimestamp)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	const op2 = "storage.sqlite.addUserTaskRelation"

	stmt, err = s.db.Prepare("INSERT INTO users_tasks(user_id, task_id) VALUES(?, ?)")
	if err != nil {
		slog.Error("Add userTask relation aql query error", "err", err)
		return fmt.Errorf("%s: %w", op2, err)
	}

	_, err = stmt.Exec(creatorId, taskId)
	if err != nil {
		slog.Error("Add userTask relation execute error", "err", err)
		return fmt.Errorf("%s: %w", op2, err)
	}

	return nil
}

func (s *Storage) GetUserTasks(id uint) ([]dto.Task, error) {
	tasks := []dto.Task{}
	const op = "storage.sqlite.GetUserTasks"
	stmt, err := s.db.Prepare("SELECT task_id, title, url, creation_timestamp FROM tasks WHERE user_id = ?")
	if err != nil {
		slog.Error("Get userTasks query error", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_ = stmt
	_ = tasks
	return tasks, nil
}
