package pgdb

import "github.com/dugtriol/BarterApp/pkg/postgres"

type TransactionRepo struct {
	*postgres.Database
}

func NewTransactionRepo(db *postgres.Database) *TransactionRepo {
	return &TransactionRepo{db}
}
