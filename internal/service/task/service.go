package task

import (
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type storage interface {
	AddTask(taskId uuid.UUID, creatorId uint, title, url string, creationTimestamp time.Time) error
}

type TaskService struct {
	storage storage
}

func NewTaskService(storage storage) *TaskService {
	return &TaskService{storage}
}

func (s *TaskService) AddTask(creatorId uint, title, url string) error {
	taskId := uuid.New()
	creationTimestamp := time.Now()
	err := s.storage.AddTask(taskId, creatorId, title, url, creationTimestamp)
	if err != nil {
		slog.Error("Failed to add task to storage", "err", err, "title", title, "url", url)
		return err
	}
	return nil
}
