package login

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"scheduler/internal/service/user"
)

type Request struct {
	Mail     string `json:"mail" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Status string `json:"status"` // Error, Ok
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, passwordChecker *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.add.New"
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

		err = passwordChecker.CheckPassword(req.Mail, req.Password)
		if err != nil {
			if errors.Is(err, user.ErrWrongPassword) {
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, Response{
					Status: "error",
					Error:  "wrong password",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "Internal Server Error",
			})
			return
		}

		log.Info("authorized", slog.String("title", req.Mail))

		render.JSON(w, r, Response{
			Status: "ok",
			Error:  "",
		})
	}
}
