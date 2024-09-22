package app

import (
	"context"
	"fmt"
	"log/slog"

	ql "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dugtriol/BarterApp/internal/controller/graph"
	"github.com/dugtriol/BarterApp/internal/service"
	middleware2 "github.com/dugtriol/BarterApp/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func NewRouter(ctx context.Context, log *slog.Logger, router *chi.Mux, services *service.Services, port string) {
	graphConfig := graph.Config{Resolvers: &graph.Resolver{Log: log, Services: services}}
	gserver := ql.NewDefaultServer(graph.NewExecutableSchema(graphConfig))

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("http://localhost:%s", port)},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware2.AuthMiddleware(ctx, log, services.Auth))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", gserver)
}
