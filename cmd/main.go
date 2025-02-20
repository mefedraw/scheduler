package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"scheduler/internal/config"
	"scheduler/internal/http-server/handlers/user/save"
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
	// TODO: init storage: postgreSQL
	storage, err := sqlite.New("./storage.db")
	if err != nil {
		log.Error("failed to init storage", err)
	}
	randomNum := rand.Int() % 1000
	mail := fmt.Sprintf("vitalik%d@gmail.com", randomNum)
	err = storage.AddUser(mail)
	if err != nil {
		log.Error("failed to add user to storage", err)
	}
	err = storage.AddTask(mail, "university", "OS lab1")
	if err != nil {
		log.Error("failed to add task to storage", err)
	}
	// TODO: init router: chi, chi render
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/user", save.New(log, storage))

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
