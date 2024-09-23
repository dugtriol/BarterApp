package app

import (
	"context"
	"fmt"
	"log/slog"

	ql "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	graph2 "github.com/dugtriol/BarterApp/graph"
	"github.com/dugtriol/BarterApp/graph/loaders"
	"github.com/dugtriol/BarterApp/internal/service"
	middleware2 "github.com/dugtriol/BarterApp/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func NewRouter(ctx context.Context, log *slog.Logger, router *chi.Mux, services *service.Services, port string) {
	graphConfig := graph2.Config{Resolvers: &graph2.Resolver{Log: log, Services: services}}
	gserver := ql.NewDefaultServer(graph2.NewExecutableSchema(graphConfig))

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("http://localhost:%s", port)},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware2.AuthMiddleware(ctx, log, services.User))
	router.Use(loaders.Middleware(log, services))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", gserver)
}
