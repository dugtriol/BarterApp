package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo/repoerrs"
	"github.com/dugtriol/BarterApp/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

const (
	userTable = "users"
)

type UserRepo struct {
	*postgres.Database
}

func (u UserRepo) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	sql, args, err := u.Builder.Insert(userTable).Columns("name", "email", "phone", "password", "city", "mode").Values(
		user.Name,
		user.Email,
		user.Phone,
		user.Password,
		user.City,
		user.Mode,
	).Suffix(
		"RETURNING id, name, email, phone, password, " +
			"city, mode",
	).ToSql()

	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - CreateUser - u.Builder.Insert: %v", err)
	}
	var output entity.User
	err = u.Cluster.QueryRow(ctx, sql, args...).Scan(
		&output.Id,
		&output.Name,
		&output.Email,
		&output.Phone,
		&output.Password,
		&output.City,
		&output.Mode,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo - CreateUser - r.Cluster.QueryRow: %v", err)
	}

	return output, nil
}

func (u UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) GetUserById(ctx context.Context, id int) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepo(db *postgres.Database) *UserRepo {
	return &UserRepo{db}
}
