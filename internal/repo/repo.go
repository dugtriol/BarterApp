package repo

import (
	"context"

	"github.com/dugtriol/BarterApp/internal/entity"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
	GetUserById(ctx context.Context, id int) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

type Repositories struct {
	User
}

//func NewRepositories(pg *postgres.Postgres) *Repositories {
//	return &Repositories{
//		User:        pgdb.NewUserRepo(pg),
//	}
//}
