package service

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dugtriol/BarterApp/graph/model"
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
	Id       string
	Email    string
	Password string
}

type UserGetByIdInput struct {
	Id string
}
type UserGetByEmailInput struct {
	Email string
}

type User interface {
	Register(ctx context.Context, input AuthRegisterInput) (entity.User, error)
	GetById(ctx context.Context, log *slog.Logger, input UserGetByIdInput) (entity.User, error)
	GetByEmail(ctx context.Context, log *slog.Logger, input UserGetByEmailInput) (entity.User, error)
	GenToken(id string) (*model.AuthToken, error)
	ParseToken(r *http.Request) (*jwt.Token, error)
	GetUsers(ctx context.Context, userIDs []string) ([]*entity.User, []error)
}

type CreateProductInput struct {
	Category    string
	Name        string
	Description string
	Image       string
	UserId      string
}

type GetByIdProductInput struct {
	Id string
}

type Product interface {
	Create(ctx context.Context, input CreateProductInput) (entity.Product, error)
	GetById(ctx context.Context, log *slog.Logger, input GetByIdProductInput) (entity.Product, error)
	All(ctx context.Context, limit, offset int) ([]entity.Product, error)
	GetByUserId(ctx context.Context, limit, offset int, userId string) ([]*model.Product, error)
}

type Services struct {
	User    User
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
		User:    NewUserService(deps.Repos.User, deps.SignKey, deps.TokenTTL),
		Product: NewProductService(deps.Repos.Product),
	}
}
