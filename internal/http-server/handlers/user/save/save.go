package save

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Mail     string `json:"mail" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Status string `json:"status"` // Error, Ok
	Error  string `json:"error,omitempty"`
}

type UserService interface {
	CheckPassword(mail, password string) error
	CreateUser(mail, password string) error
}

func New(log *slog.Logger, userService UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to deserialize request", slog.String("error", err.Error()))

			render.JSON(w, r, fmt.Errorf("failed to deserialize request: %w", err))

			return
		}
		log.Info("request body deserialized", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, fmt.Errorf("invalid request: %w", err))

			return
		}

		err = userService.CreateUser(req.Mail, req.Password)
		if err != nil {
			log.Error("failed to add user", slog.String("error", err.Error()))
			render.JSON(w, r, fmt.Errorf("failed to add user: %w", err))
			return
		}

		log.Info("user added", slog.String("mail", req.Mail))

		render.JSON(w, r, Response{
			Status: "ok",
			Error:  "",
		})
	}
}
