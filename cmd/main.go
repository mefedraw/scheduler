﻿package main

import (
	"fmt"
	"log/slog"
	"os"
	"scheduler/internal/config"
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
	err = storage.AddUser("vitalik@gmail.com")
	if err != nil {
		log.Error("failed to add user to storage", err)
	}
	// TODO: init router: chi, chi render
	// TODO: run server
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
