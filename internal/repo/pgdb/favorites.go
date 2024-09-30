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

func (p *FavoritesRepo) GetByField(ctx context.Context, field, value string) (entity.Favorites, error) {
	sql, args, _ := p.Builder.
		Select("*").
		From(favoritesTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		ToSql()
	log.Info(sql)
	var output entity.Favorites
	err := p.Cluster.QueryRow(ctx, sql, args...).Scan(
		&output.Id,
		&output.ProductId,
		&output.ProductId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Favorites{}, repoerrs.ErrNotFound
		}
		return entity.Favorites{}, fmt.Errorf("FavoritesRepo - GetByField %s - r.Cluster.QueryRow: %v", field, err)
	}
	return output, nil
}

func (p *FavoritesRepo) GetByProductId(ctx context.Context, productId string) (entity.Favorites, error) {
	return p.GetByField(ctx, "product_id", productId)
}

func (p *FavoritesRepo) GetByUserId(ctx context.Context, userId string) (entity.Favorites, error) {
	return p.GetByField(ctx, "user_id", userId)
}

func (p *FavoritesRepo) GetFavoritesByUserId(ctx context.Context, userId string) ([]entity.Favorites, error) {
	sql, args, err := p.Builder.
		Select("*").
		From(favoritesTable).
		Where("user_id = ?", userId).
		ToSql()

	log.Info(sql)

	if err != nil {
		return nil, fmt.Errorf("FavoritesRepo - GetFavoritesByUserId - query build error: %v", err)
	}

	rows, err := p.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("FavoritesRepo - GetFavoritesByUserId - r.Cluster.Query: %v", err)
	}
	defer rows.Close()

	var favorites []entity.Favorites

	for rows.Next() {
		var fav entity.Favorites
		err := rows.Scan(
			&fav.Id,
			&fav.UserId,
			&fav.ProductId,
		)
		if err != nil {
			return nil, fmt.Errorf("FavoritesRepo - GetFavoritesByUserId - rows.Scan: %v", err)
		}
		favorites = append(favorites, fav)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("FavoritesRepo - GetFavoritesByUserId - rows.Err: %v", err)
	}

	return favorites, nil
}

func (p *FavoritesRepo) DeleteIfProductDeleted(ctx context.Context, product_id string) (bool, error) {
	sql, args, err := p.Builder.Delete(favoritesTable).
		Where("product_id = ?", product_id).
		ToSql()

	log.Info(sql)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("FavoritesRepo - Delete - u.Builder.Delete: %v", err)
	}

	result, err := p.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("FavoritesRepo - Delete - r.Cluster.Exec: %v", err)
	}

	if result.RowsAffected() == 0 {
		return true, nil
	}

	return true, nil
}
