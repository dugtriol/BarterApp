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
	productTable           = "products"
	maxPaginationLimit     = 20
	defaultPaginationLimit = 5
)

type ProductRepo struct {
	*postgres.Database
}

func NewProductRepo(db *postgres.Database) *ProductRepo {
	return &ProductRepo{db}
}

func (p *ProductRepo) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	sql, args, err := p.Builder.Insert(productTable).
		Columns("name", "description", "image", "category", "user_id").
		Values(
			product.Name,
			product.Description,
			product.Image,
			product.Category,
			product.UserId,
		).
		Suffix(
			"RETURNING id, name, description, image, status, " +
				"category, user_id, created_at",
		).
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

func (p *ProductRepo) GetByField(ctx context.Context, field, value string) (entity.Product, error) {
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

func (p *ProductRepo) GetById(ctx context.Context, id string) (entity.Product, error) {
	return p.GetByField(ctx, "id", id)
}

func (p *ProductRepo) All(ctx context.Context, limit, offset int) ([]entity.Product, error) {
	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}
	if limit == 0 {
		limit = defaultPaginationLimit
	}

	sql, args, _ := p.Builder.Select("*").
		From(productTable).
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	log.Info(sql)

	return p.Pagination(ctx, sql, args)
}

func (p *ProductRepo) GetByUserId(ctx context.Context, limit, offset int, userId string) ([]entity.Product, error) {
	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}
	if limit == 0 {
		limit = defaultPaginationLimit
	}

	sql, args, _ := p.Builder.Select("*").
		From(productTable).
		Where("user_id = ?", userId).
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	log.Info(sql)
	return p.Pagination(ctx, sql, args)
}

func (p *ProductRepo) Pagination(ctx context.Context, sql string, args []interface{}) ([]entity.Product, error) {
	var output []entity.Product
	rows, err := p.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - All - r.Cluster.Query: %v", err)
	}
	for rows.Next() {
		var t entity.Product
		if err = rows.Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&t.Image,
			&t.Status,
			&t.Category,
			&t.UserId,
			&t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("ProductRepo - All - rows.Scan: %v", err)
		}
		output = append(output, t)
	}
	rows.Close()

	return output, nil
}
