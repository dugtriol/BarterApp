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
	productTable = "products"
)

type ProductRepo struct {
	*postgres.Database
}

func NewProductRepo(db *postgres.Database) *ProductRepo {
	return &ProductRepo{db}
}

func (p ProductRepo) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	sql, args, err := p.Builder.Insert(productTable).
		Columns("name", "description", "image", "category", "user_id").
		Values(
			product.Name,
			product.Description,
			product.Image,
			product.Category,
			product.UserId).
		Suffix(
			"RETURNING id, name, description, image, status, " +
				"category, user_id, created_at").
		ToSql()
	log.Info(sql)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductRepo - Create - u.Builder.Insert: %v", err)
	}
	var output entity.Product
	err = p.Cluster.QueryRow(ctx, sql, args...).Scan(
		&output.Id,
		&output.Name,
		&output.Description,
		&output.Image,
		&output.Status,
		&output.Category,
		&output.UserId,
		&output.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Product{}, repoerrs.ErrNotFound
		}
		return entity.Product{}, fmt.Errorf("ProductRepo - Create - r.Cluster.QueryRow: %v", err)
	}

	return output, nil
}

func (p ProductRepo) GetByField(ctx context.Context, field, value string) (entity.Product, error) {
	sql, args, _ := p.Builder.
		Select("*").
		From(productTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		ToSql()
	log.Info(sql)
	var output entity.Product
	err := p.Cluster.QueryRow(ctx, sql, args...).Scan(
		&output.Id,
		&output.Name,
		&output.Description,
		&output.Image,
		&output.Status,
		&output.Category,
		&output.UserId,
		&output.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Product{}, repoerrs.ErrNotFound
		}
		return entity.Product{}, fmt.Errorf("ProductRepo - GetByField %s - r.Cluster.QueryRow: %v", field, err)
	}
	return output, nil
}

func (p ProductRepo) GetById(ctx context.Context, id string) (entity.Product, error) {
	return p.GetByField(ctx, "id", id)
}
