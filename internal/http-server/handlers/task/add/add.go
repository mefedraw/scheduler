package add

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Title string `json:"title" validate:"required"`
	URl   string `json:"url" validate:"required,url"`
}

type Response struct {
	Status string `json:"status"` // Error, Ok
	Error  string `json:"error,omitempty"`
}

type TaskAdder interface {
	AddTask(title, URL string) error
}

func New(log *slog.Logger, taskAdder TaskAdder) http.HandlerFunc {
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

		err = taskAdder.AddTask(req.Title, req.URl)
		if err != nil {
			log.Error("failed to add task", slog.String("error", err.Error()))
			render.JSON(w, r, fmt.Errorf("failed to add task: %w", err))
			return
		}

		log.Info("task added", slog.String("title", req.Title))

		render.JSON(w, r, Response{
			Status: "ok",
			Error:  "",
		})
	}
}
