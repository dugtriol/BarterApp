package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dgrijalva/jwt-go"
	"github.com/dugtriol/BarterApp/config"
	graph2 "github.com/dugtriol/BarterApp/graph"
	"github.com/dugtriol/BarterApp/graph/loaders"
	"github.com/dugtriol/BarterApp/graph/model"
	"github.com/dugtriol/BarterApp/internal/repo"
	"github.com/dugtriol/BarterApp/internal/service"
	middleware2 "github.com/dugtriol/BarterApp/pkg/middleware"
	"github.com/dugtriol/BarterApp/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

const CurrentUserKey = "currentUser"

func webSocketInit(
	ctx context.Context, log *slog.Logger, initPayload transport.InitPayload, services *service.Services,
) (context.Context, *transport.InitPayload, error) {
	token := initPayload.Authorization()

	parseToken, err := services.User.ParseToken(token)
	if err != nil {
		log.Info("webSocketInit - services.User.ParseToken(token) - ", err)
		return nil, nil, errors.New("failed to parse token services.User.ParseToken(token)")
	}
	claims, ok := parseToken.Claims.(jwt.MapClaims)

	if !ok || !parseToken.Valid {
		log.Info(fmt.Sprintf("webSocketInit - !ok = %v || !parseToken.Valid = %v", !ok, !parseToken.Valid))
		log.Info(fmt.Sprintf("parseToken: %v claims: %v", parseToken, parseToken.Claims.(jwt.MapClaims)))
		return nil, nil, errors.New("failed to get claims from token")
	}

	user, err := services.User.GetById(ctx, log, service.UserGetByIdInput{Id: claims["jti"].(string)})
	if err != nil {
		log.Info("webSocketInit - services.User.GetById", ok)
		return nil, nil, errors.New("failed to get user")
	}

	// put it in context
	ctxNew := context.WithValue(ctx, CurrentUserKey, user)
	log.Info("webSocket init")
	return ctxNew, &initPayload, nil
}

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
	tokenTTL := time.Hour * 24
	repos := repo.NewRepositories(database)
	dependencies := service.ServicesDependencies{
		Repos: repos, SignKey: os.Getenv("JWT_SECRET"), TokenTTL: tokenTTL, BucketName: cfg.BucketName,
		Region: cfg.Region, EndpointResolver: cfg.EndpointResolver,
	}

	//services
	services := service.NewServices(ctx, dependencies)
	_ = services

	//handlers
	log.Info("Initializing graphql api endpoint")

	router := chi.NewRouter()
	c := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			Debug:            false,
		},
	)
	graphConfig := graph2.Config{
		Resolvers: &graph2.Resolver{
			Log: log, Services: services, ChatMessages: []*model.Message{},
			ChatObservers: map[string]chan []*model.Message{},
		},
	}
	graphConfig.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver) (
		res interface{}, err error,
	) {
		return middleware2.AuthMiddleware(ctx, log, services.User, obj, next)
	}
	var mb int64 = 1 << 20
	srv := handler.New(graph2.NewExecutableSchema(graphConfig))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(
		transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (
				context.Context, *transport.InitPayload, error,
			) {
				return webSocketInit(ctx, log, initPayload, services)
			},
		},
	)
	srv.Use(extension.Introspection{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(
		transport.MultipartForm{
			MaxMemory:     32 * mb,
			MaxUploadSize: 50 * mb,
		},
	)
	//srv.SetQueryCache(lru.New(1000))
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	//router.Use(middleware2.AuthMiddleware(ctx, log, services.User))
	router.Use(loaders.Middleware(log, services))
	router.Handle("/", playground.Handler("My GraphQL App", "/query"))
	router.Handle("/query", c.Handler(srv))

	log.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port))
	if err = http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
		log.Error("", err)
	}
	// -----
	//NewRouter(ctx, log, router, services, cfg.Port)
	// HTTP server
	//log.Info("Starting http server...")
	//log.Debug(fmt.Sprintf("Server port: %s", cfg.Port))
	//httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	//log.Info("Configuring graceful shutdown...")
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//
	//select {
	//case s := <-interrupt:
	//	log.Info("app - Run - signal: " + s.String())
	//case err = <-httpServer.Notify():
	//	log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	//}
	//
	//// Graceful shutdown
	//log.Info("Shutting down...")
	//err = httpServer.Shutdown()
	//if err != nil {
	//	log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	//}
}
