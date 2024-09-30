package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/dugtriol/BarterApp/graph/model"
	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo/repoerrs"
	"github.com/dugtriol/BarterApp/pkg/postgres"
	"github.com/google/uuid"
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
		return nil, fmt.Errorf("ProductRepo - Pagination - r.Cluster.Query: %v", err)
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
			return nil, fmt.Errorf("ProductRepo - Pagination - rows.Scan: %v", err)
		}
		output = append(output, t)
	}
	rows.Close()

	return output, nil
}

func (p *ProductRepo) FindLike(ctx context.Context, search, category, sort string) ([]entity.Product, error) {
	query := p.Builder.Select("*").From(productTable)

	query = query.Where("name LIKE ?", "%"+search+"%")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if sort == "created_at" {
		query = query.OrderBy(fmt.Sprintf("%s DESC", sort))
	} else if sort != "" {
		query = query.OrderBy(sort)
	}

	sql, args, _ := query.ToSql()
	log.Info(sql)

	first, err := p.Pagination(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - FindLike - by name and category - Pagination: %v", err)
	}

	_, err = uuid.Parse(search)
	if err != nil {
		return first, nil
	}

	query = p.Builder.Select("*").From(productTable).Where("id LIKE ?", search)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if sort == "created_at" {
		query = query.OrderBy(fmt.Sprintf("%s DESC", sort))
	} else if sort != "" {
		query = query.OrderBy(sort)
	}

	sqlsecond, argssecond, _ := query.ToSql()
	log.Info(sql)

	second, err := p.Pagination(ctx, sqlsecond, argssecond)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - FindLike - by id and category - Pagination: %v", err)
	}

	result := make([]entity.Product, 0)
	for _, product := range first {
		result = append(result, product)
	}
	for _, product := range second {
		result = append(result, product)
	}

	return result, nil
}

func (p *ProductRepo) ChangeStatus(ctx context.Context, product_id, status string) (bool, error) {
	var (
		err error
		tx  pgx.Tx
	)
	tx, err = p.Cluster.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("ProductRepo.ChangeStatus - r.Cluster.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, err := p.
		Builder.
		Update(productTable).
		Set("status", status).
		Where("id = ?", product_id).
		ToSql()
	log.Info(sql)
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("ProductRepo.ChangeStatus - tx.Exec.statusSql: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, fmt.Errorf("ProductRepo.ChangeStatus - tx.Commit: %v", err)
	}
	return true, nil
}

func (p *ProductRepo) GetByCategoryAvailable(ctx context.Context, category string) ([]entity.Product, error) {
	var status = string(model.ProductStatusAvailable)
	sql, args, _ := p.Builder.Select("*").
		From(productTable).
		Where("status = ?", status).
		Where("category = ?", category).
		OrderBy("id").
		ToSql()
	log.Info(sql)
	return p.Pagination(ctx, sql, args)
}

func (p *ProductRepo) GetByUserAvailableProducts(ctx context.Context, userId string) ([]entity.Product, error) {
	var status = string(model.ProductStatusAvailable)
	sql, args, _ := p.Builder.Select("*").
		From(productTable).
		Where("status = ?", status).
		Where("user_id = ?", userId).
		OrderBy("id").
		ToSql()
	log.Info(sql)
	return p.Pagination(ctx, sql, args)
}

func (p *ProductRepo) EditProduct(ctx context.Context, input entity.Product) (bool, error) {
	var (
		err error
		tx  pgx.Tx
	)
	tx, err = p.Cluster.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("ProductRepo.EditProduct - p.Cluster.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, err := p.
		Builder.
		Update(productTable).
		Set("category", input.Category).
		Set("name", input.Name).
		Set("description", input.Description).
		Set("image", input.Image).
		Where("id = ?", input.Id).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("ProductRepo.EditProduct - query build error: %v", err)
	}

	log.Info(sql)

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("ProductRepo.EditProduct - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, fmt.Errorf("ProductRepo.EditProduct - tx.Commit: %v", err)
	}
	return true, nil
}

func (p *ProductRepo) Delete(ctx context.Context, id string) (string, error) {
	var (
		image string
		err   error
		tx    pgx.Tx
	)

	tx, err = p.Cluster.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("ProductRepo.Delete - p.Cluster.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, err := p.
		Builder.
		Delete(productTable).
		Where("id = ?", id).
		Suffix("RETURNING image").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("ProductRepo.Delete - query build error: %v", err)
	}

	log.Info(sql)

	err = tx.QueryRow(ctx, sql, args...).Scan(&image)
	if err != nil {
		return "", fmt.Errorf("ProductRepo.Delete - tx.QueryRow: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("ProductRepo.Delete - tx.Commit: %v", err)
	}

	return image, nil
}

func (p *ProductRepo) GetLikedProductsByUserId(ctx context.Context, userId string) ([]entity.Product, error) {
	sql, args, err := p.Builder.
		Select("p.id, p.name, p.description, p.image, p.status, p.category, p.user_id, p.created_at").
		From("favorites f").
		Join("products p ON f.product_id = p.id").
		Where("f.user_id = ?", userId).
		ToSql()

	log.Info(sql)

	if err != nil {
		return nil, fmt.Errorf("ProductRepo - GetLikedProductsByUserId - query build error: %v", err)
	}

	rows, err := p.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - GetLikedProductsByUserId - r.Cluster.Query: %v", err)
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Image,
			&product.Status,
			&product.Category,
			&product.UserId,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ProductRepo - GetLikedProductsByUserId - rows.Scan: %v", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ProductRepo - GetLikedProductsByUserId - rows.Err: %v", err)
	}

	return products, nil
}

func (p *ProductRepo) GetImage(ctx context.Context, product_id string) (string, error) {
	product, err := p.GetById(ctx, product_id)
	if err != nil {
		log.Error(fmt.Sprintf("ProductRepo - GetImage -  p.GetById(ctx, product_id): %v", err))
		return "", fmt.Errorf("failed to get product: %v", err)
	}
	return product.Image, nil
}
