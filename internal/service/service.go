package service

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
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

type UserEditProfile struct {
	Id    string
	Name  string
	Email string
	Phone string
	City  string
}

type User interface {
	Register(ctx context.Context, input AuthRegisterInput) (entity.User, error)
	GetById(ctx context.Context, log *slog.Logger, input UserGetByIdInput) (entity.User, error)
	GetByEmail(ctx context.Context, log *slog.Logger, input UserGetByEmailInput) (entity.User, error)
	GenToken(id string) (*model.AuthToken, error)
	ParseTokenFromRequest(r *http.Request) (*jwt.Token, error)
	GetUsers(ctx context.Context, userIDs []string) ([]*entity.User, []error)
	ParseToken(token string) (*jwt.Token, error)
	UpdateProfile(ctx context.Context, input UserEditProfile) (bool, error)
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
	FindLike(ctx context.Context, category model.ProductCategory, search string, sort model.ProductSort) ([]*model.Product, error)
	GetByUserAvailableProducts(ctx context.Context, userId string) ([]*model.Product, error)
	GetByCategoryAvailable(ctx context.Context, category string) ([]*model.Product, error)
	EditProduct(ctx context.Context, input EditProductInput) (bool, error)
	Delete(ctx context.Context, id string) (string, error)
	GetLikedProductsByUserId(ctx context.Context, userId string) ([]*model.Product, error)
	IsImageChanged(ctx context.Context, product_id string) (string, error)
}

type Favorites interface {
	Add(ctx context.Context, input entity.Favorites) (string, error)
	Delete(ctx context.Context, input entity.Favorites) (bool, error)
	GetFavoritesByUserId(ctx context.Context, userId string) ([]*model.Favorites, error)
	DeleteIfProductDeleted(ctx context.Context, product_id string) (bool, error)
}

type Transaction interface {
	Create(ctx context.Context, input CreateTransactionInput) (string, error)
	UpdateOngoingOrDeclined(ctx context.Context, input UpdateStatusInput) (bool, error)
	GetByBuyer(ctx context.Context, buyer_id string) ([]*model.Transaction, error)
	GetByOwner(ctx context.Context, owner_id string) ([]*model.Transaction, error)
	UpdateDone(ctx context.Context, input UpdateStatusInput) (bool, error)
	GetOngoing(ctx context.Context, buyer_id string) ([]*model.Transaction, error)
	GetCreated(ctx context.Context, owner_id string) ([]*model.Transaction, error)
	GetArchive(ctx context.Context, id string) ([]*model.Transaction, error)
}

type File interface {
	Upload(ctx context.Context, file graphql.Upload) (string, error)
	Delete(ctx context.Context, path string) (bool, error)
	BuildImageURL(pathName string) string
}

type Services struct {
	User        User
	Product     Product
	Favorites   Favorites
	Transaction Transaction
	File        File
}

type ServicesDependencies struct {
	Repos *repo.Repositories
	//Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration

	BucketName       string
	Region           string
	EndpointResolver string
}

func NewServices(ctx context.Context, deps ServicesDependencies) *Services {
	return &Services{
		User:        NewUserService(deps.Repos.User, deps.SignKey, deps.TokenTTL),
		Product:     NewProductService(deps.Repos.Product),
		Favorites:   NewFavoritesService(deps.Repos.Favorites),
		Transaction: NewTransactionService(deps.Repos.Transaction, deps.Repos.Product),
		File:        NewFileService(ctx, deps.BucketName, deps.Region, deps.EndpointResolver),
	}
}
