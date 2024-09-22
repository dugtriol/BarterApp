package repo

import (
	"context"

	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo/pgdb"
	"github.com/dugtriol/BarterApp/pkg/postgres"
)

type User interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetById(ctx context.Context, id string) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type Product interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	GetById(ctx context.Context, id string) (entity.Product, error)
}
type Repositories struct {
	User
	Product
}

func NewRepositories(pg *postgres.Database) *Repositories {
	return &Repositories{
		User:    pgdb.NewUserRepo(pg),
		Product: pgdb.NewProductRepo(pg),
	}
}
