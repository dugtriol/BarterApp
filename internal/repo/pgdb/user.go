package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo/repoerrs"
	"github.com/dugtriol/BarterApp/pkg/postgres"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

const (
	userTable = "users"
)

type UserRepo struct {
	*postgres.Database
}

func NewUserRepo(db *postgres.Database) *UserRepo {
	return &UserRepo{db}
}

func (u UserRepo) Create(ctx context.Context, user entity.User) (entity.User, error) {
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
	log.Info(sql)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - Create - u.Builder.Insert: %v", err)
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
		return entity.User{}, fmt.Errorf("UserRepo - Create - r.Cluster.QueryRow: %v", err)
	}

	return output, nil
}

func (u UserRepo) GetByField(ctx context.Context, field, value string) (entity.User, error) {
	sql, args, _ := u.Builder.
		Select("*").
		From(userTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		ToSql()
	log.Info(sql)
	var output entity.User
	err := u.Cluster.QueryRow(ctx, sql, args...).Scan(
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
		return entity.User{}, fmt.Errorf("UserRepo - GetByField %s - r.Cluster.QueryRow: %v", field, err)
	}
	return output, nil
}

func (u UserRepo) GetById(ctx context.Context, id string) (entity.User, error) {
	return u.GetByField(ctx, "id", id)
}

func (u UserRepo) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	return u.GetByField(ctx, "name", username)
}

func (u UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	return u.GetByField(ctx, "email", email)
}
