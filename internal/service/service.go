package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo"
)

type AuthRegisterInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
	City     string
	Mode     string
}

type AuthGenerateTokenInput struct {
	Email    string
	Password string
}

type UserGetByIdInput struct {
	Id string
}

type Auth interface {
	Register(ctx context.Context, input AuthRegisterInput) (entity.User, error)
	GetById(ctx context.Context, log *slog.Logger, input UserGetByIdInput) (entity.User, error)
	//GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	//ParseToken(token string) (int, error)
}

type Product interface {
}

type Services struct {
	Auth    Auth
	Product Product
}

type ServicesDependencies struct {
	Repos *repo.Repositories
	//Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:    NewAuthService(deps.Repos.User, deps.SignKey, deps.TokenTTL),
		Product: NewProductService(deps.Repos.Product),
	}
}
