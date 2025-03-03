package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"scheduler/internal/config"
	"scheduler/internal/http-server/handlers/task/add"
	"scheduler/internal/http-server/handlers/user/login"
	"scheduler/internal/http-server/handlers/user/save"
	"scheduler/internal/service/task"
	"scheduler/internal/service/user"
	"scheduler/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: init logger: slog
	log := setupLogger(cfg.Env)
	log.Info("starting project", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	// TODO: init storage: sqlite
	storage, err := sqlite.New("./SqliteStorage.db")
	if err != nil {
		log.Error("failed to init storage", err)
	}
	//pass := gofakeit.Password(true, true, true, false, false, 10)
	//mail := gofakeit.Email()
	//err = storage.AddUser(mail, pass)
	//if err != nil {
	//	log.Error("failed to add user to storage", err)
	//}
	//err = storage.AddTask(gofakeit.CelebritySport(), gofakeit.URL())
	//if err != nil {
	//	log.Error("failed to add task to storage", err)
	//}
	userService := user.NewUserService(storage)
	taskService := task.NewTaskService(storage)
	//userService.CreateUser("vitalik228@gmail.com", "vitalik228")
	// TODO: init router: chi, chi render
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/user", save.New(log, userService))
	router.Post("/task", add.New(log, taskService))
	router.Get("/login", login.New(log, userService))

	log.Info("starting server", slog.String("address", cfg.Address))
	// TODO: run server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServerConfig.Timeout,
		WriteTimeout: cfg.HttpServerConfig.Timeout,
		IdleTimeout:  cfg.HttpServerConfig.IddleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", err)
	}

	log.Error("failed to start server")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
