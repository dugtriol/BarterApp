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
	GetUsers(ctx context.Context, userIDs []string) ([]*entity.User, []error)
	UpdateProfile(ctx context.Context, input entity.User) (bool, error)
}

type Product interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	GetById(ctx context.Context, id string) (entity.Product, error)
	All(ctx context.Context, limit, offset int) ([]entity.Product, error)
	GetByUserId(ctx context.Context, limit, offset int, userId string) ([]entity.Product, error)
	FindLike(ctx context.Context, search, category, sort string) ([]entity.Product, error)
	ChangeStatus(ctx context.Context, product_id, status string) (bool, error)
	GetByUserAvailableProducts(ctx context.Context, userId string) ([]entity.Product, error)
	GetByCategoryAvailable(ctx context.Context, category string) ([]entity.Product, error)
	EditProduct(ctx context.Context, input entity.Product) (bool, error)
	Delete(ctx context.Context, id string) (string, error)
	GetLikedProductsByUserId(ctx context.Context, userId string) ([]entity.Product, error)
	GetImage(ctx context.Context, product_id string) (string, error)
}

type Favorites interface {
	Add(ctx context.Context, input entity.Favorites) (string, error)
	Delete(ctx context.Context, input entity.Favorites) (bool, error)
	GetByProductId(ctx context.Context, productId string) (entity.Favorites, error)
	GetByUserId(ctx context.Context, userId string) (entity.Favorites, error)
	GetFavoritesByUserId(ctx context.Context, userId string) ([]entity.Favorites, error)
	DeleteIfProductDeleted(ctx context.Context, product_id string) (bool, error)
}

type Transaction interface {
	Create(ctx context.Context, input entity.Transaction) (entity.Transaction, error)
	GetByOwner(ctx context.Context, value string) ([]entity.Transaction, error)
	GetByBuyer(ctx context.Context, value string) ([]entity.Transaction, error)
	ChangeStatus(ctx context.Context, transactionID, status string) (entity.Transaction, error)
	CheckIsOwner(ctx context.Context, userId string, transactionId string) (bool, error)
	CheckIsBuyer(ctx context.Context, userId string, transactionId string) (bool, error)
	UpdateTime(ctx context.Context, transaction_id string) (bool, error)
	GetOngoing(ctx context.Context, buyer_id string) ([]entity.Transaction, error)
	GetCreated(ctx context.Context, owner_id string) ([]entity.Transaction, error)
	GetArchive(ctx context.Context, id string) ([]entity.Transaction, error)
}

type Repositories struct {
	User
	Product
	Favorites
	Transaction
}

func NewRepositories(pg *postgres.Database) *Repositories {
	return &Repositories{
		User:        pgdb.NewUserRepo(pg),
		Product:     pgdb.NewProductRepo(pg),
		Favorites:   pgdb.NewFavoritesRepo(pg),
		Transaction: pgdb.NewTransactionRepo(pg),
	}
}
