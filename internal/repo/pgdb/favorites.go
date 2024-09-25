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
	favoritesTable = "favorites"
)

type FavoritesRepo struct {
	*postgres.Database
}

func NewFavoritesRepo(db *postgres.Database) *FavoritesRepo {
	return &FavoritesRepo{db}
}

func (p *FavoritesRepo) Add(ctx context.Context, input entity.Favorites) (string, error) {
	sql, args, err := p.Builder.Insert(favoritesTable).
		Columns("user_id", "product_id").
		Values(
			input.UserId,
			input.ProductId,
		).
		Suffix(
			"RETURNING id",
		).
		ToSql()
	log.Info(sql)
	if err != nil {
		return "", fmt.Errorf("FavoritesRepo - Add - u.Builder.Insert: %v", err)
	}
	var id string
	err = p.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repoerrs.ErrNotFound
		}
		return "", fmt.Errorf("FavoritesRepo - Create - r.Cluster.QueryRow: %v", err)
	}
	return id, nil
}

func (p *FavoritesRepo) Delete(ctx context.Context, input entity.Favorites) (bool, error) {
	sql, args, err := p.Builder.Delete(favoritesTable).
		Where("user_id = ? AND product_id = ?", input.UserId, input.ProductId).
		ToSql()

	log.Info(sql)

	if err != nil {
		return false, fmt.Errorf("FavoritesRepo - Delete - u.Builder.Delete: %v", err)
	}

	result, err := p.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("FavoritesRepo - Delete - r.Cluster.Exec: %v", err)
	}

	if result.RowsAffected() == 0 {
		return false, repoerrs.ErrNotFound
	}

	return true, nil
}
