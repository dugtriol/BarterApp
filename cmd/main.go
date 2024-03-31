package main

import (
	`context`
	"log/slog"
	`net/http`
	"os"

	`github.com/dugtriol/BarterApp/internal/config`
	`github.com/dugtriol/BarterApp/internal/http-server/handlers/authentication`
	`github.com/dugtriol/BarterApp/internal/http-server/handlers/product`
	`github.com/dugtriol/BarterApp/internal/http-server/handlers/user`
	mwLogger `github.com/dugtriol/BarterApp/internal/http-server/middleware/logger`
	`github.com/dugtriol/BarterApp/internal/pkg/db`
	"github.com/dugtriol/BarterApp/internal/pkg/lib/logger/handlers/slogpretty"
	`github.com/dugtriol/BarterApp/internal/pkg/lib/logger/sl`
	`github.com/dugtriol/BarterApp/internal/pkg/storage/postgresql`
	`github.com/go-chi/chi/v5`
	`github.com/go-chi/chi/v5/middleware`
	"github.com/go-chi/jwtauth/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
	secret   = "secret"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

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
	_ = storage

	if err != nil {
		log.Error("failed to init storage ", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.URLFormat)
	router.Use(middleware.Recoverer)

	router.Group(
		func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.Delete("/user/{id}", user.Delete(ctx, log, storage))
			r.Post("/user/update/city/{id}", user.UpdateCity(ctx, log, storage))
			r.Post("/user/update/password/{id}", user.UpdatePassword(ctx, log, storage))
			r.Get("/user/{id}", user.Get(ctx, log, storage))
			r.Post("/product", product.Save(ctx, log, storage))
			r.Delete("/product", product.Delete(ctx, log, storage))
		},
	)

	router.Group(
		func(r chi.Router) {
			r.Post("/register", user.Save(ctx, log, storage))
			r.Post("/login", authentication.LogIn(ctx, log, tokenAuth, storage))
			r.Post("/logout", authentication.LogOut(ctx, log))
		},
	)

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
