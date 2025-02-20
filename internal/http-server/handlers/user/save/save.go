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
	Mail string `json:"mail" validate:"required,email"`
}

type Response struct {
	Status string `json:"status"` // Error, Ok
	Error  string `json:"error,omitempty"`
}

type UserSaver interface {
	AddUser(mail string) error
}

func New(log *slog.Logger, urlSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
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

		log.Info("user added", slog.String("mail", req.Mail))

		render.JSON(w, r, Response{
			Status: "ok",
			Error:  "",
		})
	}
}
