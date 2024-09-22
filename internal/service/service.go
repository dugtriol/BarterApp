package service

import (
	"context"
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

type Auth interface {
	Register(ctx context.Context, input AuthRegisterInput) (entity.User, error)
	//GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	//ParseToken(token string) (int, error)
}

type Services struct {
	Auth Auth
}

type ServicesDependencies struct {
	Repos *repo.Repositories
	//Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos.User, deps.SignKey, deps.TokenTTL),
	}
}
