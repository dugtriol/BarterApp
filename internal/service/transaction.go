package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dugtriol/BarterApp/graph/model"
	"github.com/dugtriol/BarterApp/graph/scalar"
	"github.com/dugtriol/BarterApp/internal/controller"
	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo"
	"github.com/dugtriol/BarterApp/internal/repo/repoerrs"
	log "github.com/sirupsen/logrus"
)

type TransactionService struct {
	transactionRepo repo.Transaction
	productRepo     repo.Product
}

func NewTransactionService(
	transactionRepo repo.Transaction,
	productRepo repo.Product,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

type CreateTransactionInput struct {
	Owner          string `yaml:"owner"`
	Buyer          string `yaml:"buyer"`
	ProductIdBuyer string `yaml:"product_id_buyer"`
	ProductIdOwner string `yaml:"product_id_owner"`
	Shipping       string `yaml:"shipping"`
	Address        string `yaml:"address"`
}

func (p *TransactionService) Create(ctx context.Context, input CreateTransactionInput) (string, error) {
	transaction := entity.Transaction{
		Owner:          input.Owner,
		Buyer:          input.Buyer,
		ProductIdBuyer: input.ProductIdBuyer,
		ProductIdOwner: input.ProductIdOwner,
		Shipping:       input.Shipping,
		Address:        input.Address,
	}

	output, err := p.transactionRepo.Create(ctx, transaction)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return "", ErrAlreadyExists
		}
		log.Error("TransactionService.Create - p.transactionRepo.Create: %v", err)
		return "", ErrCannotCreate
	}

	var status = string(model.ProductStatusExchanging)
	ok, err := p.productRepo.ChangeStatus(ctx, transaction.ProductIdOwner, status)

	if !ok {
		log.Error("Не получилось обновить статус товара ProductIdOwner", err)
		return "", err
	}

	ok, err = p.productRepo.ChangeStatus(ctx, transaction.ProductIdBuyer, status)

	if !ok {
		log.Error("Не получилось обновить статус товара ProductIdBuyer", err)
		return "", err
	}

	return output.Id, nil
}

type UpdateStatusInput struct {
	TransactionId string
	UserId        string
	Status        string
}

func (p *TransactionService) UpdateOngoingOrDeclined(ctx context.Context, input UpdateStatusInput) (bool, error) {
	isOwner, err := p.transactionRepo.CheckIsOwner(ctx, input.UserId, input.TransactionId)
	if err != nil {
		log.Error("TransactionService.Approve - p.transactionRepo.CheckIsOwner: %v", err)
		return false, err
	}

	if !isOwner {
		log.Error("Пользователь не является владельцем товара", err)
		return false, err
	}

	transaction, err := p.transactionRepo.ChangeStatus(ctx, input.TransactionId, input.Status)
	if err != nil {
		log.Error("TransactionService.Approve - p.transactionRepo.ChangeStatus: %v", err)
		return false, err
	}

	productStatus := p.GetStatus(input.Status)
	ok, err := p.productRepo.ChangeStatus(ctx, transaction.ProductIdOwner, string(productStatus))

	if !ok {
		log.Error("Не получилось обновить статус товара ProductIdOwner", err)
		return false, err
	}

	ok, err = p.productRepo.ChangeStatus(ctx, transaction.ProductIdBuyer, string(productStatus))

	if !ok {
		log.Error("Не получилось обновить статус товара ProductIdBuyer", err)
		return false, err
	}
	return ok, nil
}

func (p *TransactionService) GetByBuyer(ctx context.Context, buyer_id string) ([]*model.Transaction, error) {
	output, err := p.transactionRepo.GetByBuyer(ctx, buyer_id)
	if err != nil {
		log.Error("TransactionService - p.transactionRepo.GetByBuyer: ", err)
		return nil, controller.ErrNotFound
	}
	return p.ParseTransactionArray(output)
}

func (p *TransactionService) GetByOwner(ctx context.Context, owner_id string) ([]*model.Transaction, error) {
	output, err := p.transactionRepo.GetByOwner(ctx, owner_id)
	if err != nil {
		log.Error("TransactionService - p.transactionRepo.GetByOwner: ", err)
		return nil, controller.ErrNotFound
	}
	return p.ParseTransactionArray(output)
}

func (p *TransactionService) ParseTransactionArray(output []entity.Transaction) ([]*model.Transaction, error) {
	var err error
	result := make([]*model.Transaction, len(output))
	for i, item := range output {
		var shipping model.TransactionShipping
		err = shipping.UnmarshalGQL(item.Shipping)
		if err != nil {
			log.Error("TransactionService - shipping.UnmarshalGQL: ", err)
			return nil, controller.ErrNotValid
		}

		var status model.TransactionStatus
		err = status.UnmarshalGQL(item.Status)
		if err != nil {
			log.Error("TransactionService - status.UnmarshalGQL(item.Status): ", err)
			return nil, controller.ErrNotValid
		}

		//var temp *model.Product
		result[i] = &model.Transaction{
			ID:             item.Id,
			Owner:          item.Owner,
			Buyer:          item.Buyer,
			ProductIDOwner: item.ProductIdOwner,
			ProductIDBuyer: item.ProductIdBuyer,
			CreatedAt:      scalar.DateTime(item.CreatedAt),
			Shipping:       shipping,
			Address:        item.Address,
			Status:         status,
		}
		//log.Info(result[i])
		//result = append(result, temp)
	}
	//log.Info(result)
	return result, nil
}

func (p *TransactionService) UpdateDone(ctx context.Context, input UpdateStatusInput) (bool, error) {
	isBuyer, err := p.transactionRepo.CheckIsBuyer(ctx, input.UserId, input.TransactionId)
	if err != nil {
		log.Error("TransactionService.Approve - p.transactionRepo.CheckIsOwner: %v", err)
		return false, err
	}

	if !isBuyer {
		log.Error("Пользователь не является покупателем товара", err)
		return false, err
	}

	transaction, err := p.transactionRepo.ChangeStatus(ctx, input.TransactionId, input.Status)
	if err != nil {
		log.Error("TransactionService.Approve - p.transactionRepo.ChangeStatus: %v", err)
		return false, err
	}

	status := p.GetStatus(input.Status)
	ok, err := p.productRepo.ChangeStatus(ctx, transaction.ProductIdOwner, string(status))

	if !ok {
		log.Error("Не получилось обновить статус товара ProductIdOwner", err)
		return false, err
	}

	ok, err = p.productRepo.ChangeStatus(ctx, transaction.ProductIdBuyer, string(status))

	if !ok {
		log.Error("Не получилось обновить статус товара ProductIdBuyer", err)
		return false, err
	}
	return ok, nil
}

func (p *TransactionService) GetStatus(status string) model.ProductStatus {
	var res model.ProductStatus
	switch status {
	case string(model.TransactionStatusCreated), string(model.TransactionStatusOngoing):
		res = model.ProductStatusExchanging
	case string(model.TransactionStatusDone):
		res = model.ProductStatusExchanged
	default:
		res = model.ProductStatusAvailable
	}
	return res
}

func (p *TransactionService) GetOngoing(ctx context.Context, buyer_id string) ([]*model.Transaction, error) {
	output, err := p.transactionRepo.GetOngoing(ctx, buyer_id)
	log.Info(fmt.Sprintf("GetOngoing - %v", output))
	if err != nil {
		log.Error("TransactionService - p.transactionRepo.GetOngoing: ", err)
		return nil, controller.ErrNotFound
	}
	return p.ParseTransactionArray(output)
}

func (p *TransactionService) GetCreated(ctx context.Context, owner_id string) ([]*model.Transaction, error) {
	output, err := p.transactionRepo.GetCreated(ctx, owner_id)
	log.Info(fmt.Sprintf("GetCreated - %v", output))
	if err != nil {
		log.Error("TransactionService - p.transactionRepo.GetCreated: ", err)
		return nil, controller.ErrNotFound
	}
	return p.ParseTransactionArray(output)
}

func (p *TransactionService) GetArchive(ctx context.Context, id string) ([]*model.Transaction, error) {
	output, err := p.transactionRepo.GetArchive(ctx, id)
	log.Info(fmt.Sprintf("GetArchive - %v", output))
	if err != nil {
		log.Error("TransactionService - p.transactionRepo.GetArchive: ", err)
		return nil, controller.ErrNotFound
	}
	return p.ParseTransactionArray(output)
}
