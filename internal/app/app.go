package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ql "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dugtriol/BarterApp/config"
	"github.com/dugtriol/BarterApp/internal/controller/graph"
	"github.com/dugtriol/BarterApp/internal/repo"
	"github.com/dugtriol/BarterApp/internal/service"
	"github.com/dugtriol/BarterApp/pkg/httpserver"
	"github.com/dugtriol/BarterApp/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func Run(configPath string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(cfg)

	//logger
	log := setLogger(cfg.Level)
	log.Info("Init logger")

	//postgres
	database, err := postgres.New(ctx, cfg.Conn, postgres.MaxPoolSize(cfg.MaxPoolSize))
	if err != nil {
		fmt.Println(err.Error())
	}

	//repositories
	repos := repo.NewRepositories(database)
	dependencies := service.ServicesDependencies{Repos: repos}

	//services
	services := service.NewServices(dependencies)
	_ = services

	//handlers
	log.Info("Initializing graphql api endpoint")

	graphConfig := graph.Config{Resolvers: &graph.Resolver{}}
	gserver := ql.NewDefaultServer(graph.NewExecutableSchema(graphConfig))
	router := chi.NewRouter()
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", gserver)

	// HTTP server
	log.Info("Starting http server...")
	log.Debug(fmt.Sprintf("Server port: %s", cfg.Port))
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}
