package main

import (
	"log/slog"
	"main/internal/config"
	"main/internal/lib/logger/sl"
	"main/internal/storage/sqlite"
	"os"
)

const (
	envLocal = "local"
	envProf  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("Hello world", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StaragePath)
	if err != nil {
		log.Error("Error, cant create new storage", sl.Err(err)) // Возможно работает без sl.Err()
		os.Exit(1)
	}
	err = storage.DeleteURL("google1")
	if err != nil {
		log.Error("cant delete url", sl.Err(err))
		os.Exit(1)
	}
	log.Info("Success delete")
	_ = storage
}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProf:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return log

}
