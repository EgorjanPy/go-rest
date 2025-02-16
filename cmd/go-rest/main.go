package main

import (
	"log/slog"
	"main/internal/config"
	"main/internal/http-server/middleware/mwLogger"
	"main/internal/http-server/middleware/mwLogger/handlers/url/save"
	"main/internal/lib/logger/handlers/slogpretty"
	"main/internal/lib/logger/sl"
	"main/internal/storage/sqlite"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	log.Debug("Debug logger enabled")
	log.Error("Error logger enabled")

	storage, err := sqlite.New(cfg.StaragePath)
	if err != nil {
		log.Error("Error, cant create new storage", sl.Err(err)) // Возможно работает без sl.Err()
		os.Exit(1)
	}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/url", save.New(log, storage))

	log.Info("server starn", slog.String("addres", cfg.HttpServer.Addres))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.Idle_timeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	_ = storage
}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProf:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
