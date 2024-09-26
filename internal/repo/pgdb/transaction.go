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
	transactionsTable = "transactions"
)

type TransactionRepo struct {
	*postgres.Database
}

func NewTransactionRepo(db *postgres.Database) *TransactionRepo {
	return &TransactionRepo{db}
}

func (p *TransactionRepo) Create(ctx context.Context, input entity.Transaction) (entity.Transaction, error) {
	sql, args, err := p.Builder.Insert(transactionsTable).
		Columns("owner", "buyer", "product_id_owner", "product_id_buyer", "shipping", "address").
		Values(
			input.Owner,
			input.Buyer,
			input.ProductIdOwner,
			input.ProductIdBuyer,
			input.Shipping,
			input.Address,
		).
		Suffix(
			"RETURNING id, product_id_owner, product_id_buyer",
		).
		ToSql()
	log.Info(sql)

	if err != nil {
		return entity.Transaction{}, fmt.Errorf("TransactionRepo - Create - u.Builder.Insert: %v", err)
	}
	var id, productIDOwner, productIDBuyer string
	err = p.Cluster.QueryRow(ctx, sql, args...).Scan(
		&id,
		&productIDOwner,
		&productIDBuyer,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Transaction{}, repoerrs.ErrNotFound
		}
		return entity.Transaction{}, fmt.Errorf("TransactionRepo - Create - r.Cluster.QueryRow: %v", err)
	}
	output := entity.Transaction{
		Id:             id,
		ProductIdBuyer: productIDBuyer,
		ProductIdOwner: productIDOwner,
	}
	return output, nil
}

func (p *TransactionRepo) GetByField(ctx context.Context, field, value string) ([]entity.Transaction, error) {
	sql, args, _ := p.Builder.Select("*").
		From(transactionsTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		OrderBy("id").
		ToSql()
	log.Info(sql)
	return p.Pagination(ctx, sql, args)
}

func (p *TransactionRepo) Pagination(ctx context.Context, sql string, args []interface{}) (
	[]entity.Transaction, error,
) {
	var output []entity.Transaction
	rows, err := p.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - Pagination - r.Cluster.Query: %v", err)
	}
	for rows.Next() {
		var t entity.Transaction
		if err = rows.Scan(
			&t.Id,
			&t.Owner,
			&t.Buyer,
			&t.ProductIdOwner,
			&t.ProductIdBuyer,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.Shipping,
			&t.Address,
			&t.Status,
		); err != nil {
			return nil, fmt.Errorf("TransactionRepo - Pagination - rows.Scan: %v", err)
		}
		output = append(output, t)
	}
	rows.Close()

	return output, nil
}

func (p *TransactionRepo) GetByOwner(ctx context.Context, value string) ([]entity.Transaction, error) {
	return p.GetByField(ctx, "owner", value)
}

func (p *TransactionRepo) GetByBuyer(ctx context.Context, value string) ([]entity.Transaction, error) {
	return p.GetByField(ctx, "buyer", value)
}

func (p *TransactionRepo) ChangeStatus(ctx context.Context, transactionID, status string) (entity.Transaction, error) {
	var (
		err error
		tx  pgx.Tx
	)
	tx, err = p.Cluster.Begin(ctx)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("TransactionRepo.ChangeStatus - r.Cluster.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Prepare SQL query with RETURNING clause
	sql, args, err := p.
		Builder.
		Update(transactionsTable).
		Set("status", status).
		Where("id = ?", transactionID).
		Suffix(
			"RETURNING product_id_owner, product_id_buyer",
		).
		ToSql()
	log.Info(sql)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("TransactionRepo.ChangeStatus - p.Builder.Update: %v", err)
	}

	// Declare variables to hold the returned values
	var productIDOwner, productIDBuyer string

	// Execute the query and scan the returned values into variables
	err = tx.QueryRow(ctx, sql, args...).Scan(&productIDOwner, &productIDBuyer)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("TransactionRepo.ChangeStatus - tx.QueryRow: %v", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("TransactionRepo.ChangeStatus - tx.Commit: %v", err)
	}

	// Return the updated transaction entity
	return entity.Transaction{
		Id:             transactionID,
		ProductIdOwner: productIDOwner,
		ProductIdBuyer: productIDBuyer,
		Status:         status,
	}, nil
}

func (p *TransactionRepo) CheckIs(ctx context.Context, field, id string, transactionId string) (bool, error) {
	sql, args, err := p.Builder.
		Select("COUNT(*)").
		From(transactionsTable).
		Where("id = ?", transactionId).
		Where(fmt.Sprintf("%s = ?", field), id).
		ToSql()
	log.Info(sql)
	if err != nil {
		return false, fmt.Errorf("TransactionRepo - CheckIsOwner - p.Builder: %v", err)
	}

	var count int
	err = p.Cluster.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("TransactionRepo - CheckIsOwner - QueryRow: %v", err)
	}

	return count > 0, nil
}

func (p *TransactionRepo) CheckIsOwner(ctx context.Context, userId string, transactionId string) (bool, error) {
	return p.CheckIs(ctx, "owner", userId, transactionId)
}

func (p *TransactionRepo) CheckIsBuyer(ctx context.Context, userId string, transactionId string) (bool, error) {
	return p.CheckIs(ctx, "buyer", userId, transactionId)
}
