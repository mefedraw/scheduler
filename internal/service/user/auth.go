package user

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type storage interface {
	GetPassword(mail string) (string, error)
	UserExists(mail string) (string, error)
	AddUser(mail, password string, registrationDate time.Time) error
	AddTask(taskId uuid.UUID, title, url string, taskAddDate time.Time) error
}

type Service struct {
	storage storage
}

func NewUserService(storage storage) *Service {
	return &Service{storage}
}

func (u *Service) CheckPassword(mail, password string) error {
	pass, err := u.storage.GetPassword(mail)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrStorageFail, err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
	if err != nil {
		slog.Error("password is incorrect", "err", err, "mail", mail, "password", password)
		return fmt.Errorf("%w:%w", ErrWrongPassword, err)
	}
	return nil
}

func (u *Service) CreateUser(mail, password string) error {
	userExists, _ := u.storage.UserExists(mail)
	if userExists == "" {
		regDate := time.Now()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		err = u.storage.AddUser(mail, string(hashedPassword), regDate)
		if err != nil {
			slog.Error("error creating new user", "err", err, "mail", mail, "password", password)
			return err
		}
		slog.Info("created new user", "mail", mail, "password", password)
		return nil
	}
	slog.Error("user already exists", "mail", mail)
	return fmt.Errorf("user %s already exists", mail)
}
