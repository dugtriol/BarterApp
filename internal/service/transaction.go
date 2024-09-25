package service

import (
	"github.com/dugtriol/BarterApp/internal/repo"
)

type TransactionService struct {
	transactionRepo repo.Transaction
}

func NewTransactionService(
	transactionRepo repo.Transaction,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}
