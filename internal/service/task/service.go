package task

import (
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type storage interface {
	AddTask(taskId uuid.UUID, title, url string, taskAddDate time.Time) error
}

type TaskService struct {
	storage storage
}

func NewTaskService(storage storage) *TaskService {
	return &TaskService{storage}
}

func (s *TaskService) AddTask(title, url string) error {
	taskId := uuid.New()
	taskAddDate := time.Now()
	err := s.storage.AddTask(taskId, title, url, taskAddDate)
	if err != nil {
		slog.Error("Failed to add task to storage", "err", err, "title", title, "url", url)
		return err
	}
	return nil
}
