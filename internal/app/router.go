package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	ql "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	graph2 "github.com/dugtriol/BarterApp/graph"
	"github.com/dugtriol/BarterApp/graph/loaders"
	"github.com/dugtriol/BarterApp/graph/model"
	"github.com/dugtriol/BarterApp/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

func NewRouter(ctx context.Context, log *slog.Logger, router *chi.Mux, services *service.Services, port string) {
	graphConfig := graph2.Config{
		Resolvers: &graph2.Resolver{
			Log: log, Services: services, ChatMessages: []*model.Message{},
			ChatObservers: map[string]chan []*model.Message{},
		},
	}
	gserver := ql.NewDefaultServer(graph2.NewExecutableSchema(graphConfig))
	gserver.AddTransport(transport.POST{})
	gserver.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	gserver.Use(extension.Introspection{})

	router.Use(
		cors.New(
			cors.Options{
				AllowedOrigins:   []string{fmt.Sprintf("http://localhost:%s", port)},
				AllowCredentials: true,
				Debug:            true,
			},
		).Handler,
	)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	//router.Use(middleware2.AuthMiddleware(ctx, log, services.User))
	router.Use(loaders.Middleware(log, services))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", gserver)
}
