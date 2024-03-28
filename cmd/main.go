package main

import (
	"context"
	`log/slog`
	`net/http`
	`os`

	`github.com/dugtriol/BarterApp/internal/config`
	`github.com/dugtriol/BarterApp/internal/http-server/handlers/user/url/get`
	`github.com/dugtriol/BarterApp/internal/http-server/handlers/user/url/save`
	mwLogger `github.com/dugtriol/BarterApp/internal/http-server/middleware/logger`
	"github.com/dugtriol/BarterApp/internal/pkg/db"
	`github.com/dugtriol/BarterApp/internal/pkg/lib/logger/handlers/slogpretty`
	`github.com/dugtriol/BarterApp/internal/pkg/lib/logger/sl`
	"github.com/dugtriol/BarterApp/internal/pkg/storage/postgresql"
	`github.com/go-chi/chi/v5`
	`github.com/go-chi/chi/v5/middleware`
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	defer database.GetPool(ctx).Close()

	storage := postgresql.New(database)

	if err != nil {
		log.Error("failed to init storage ", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	router := chi.NewRouter()
	// Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.RequestID)
	// Логирование всех запросов
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	// Парсер URLов поступающих запросов
	router.Use(middleware.URLFormat)
	// Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.Recoverer)

	// TODO: router.POST, router.Get
	router.Post("/user", save.New(ctx, log, storage))
	router.Get("/user/{id}", get.New(ctx, log, storage))

	serv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("start server")
	if err := serv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	// slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
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
