package service

import (
	"context"
	"errors"

	"log/slog"

	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo"
	"github.com/dugtriol/BarterApp/internal/repo/repoerrs"
	log "github.com/sirupsen/logrus"
)

type ProductService struct {
	productRepo repo.Product
}

func NewProductService(
	productRepo repo.Product,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (p ProductService) Create(ctx context.Context, input CreateProductInput) (entity.Product, error) {
	product := entity.Product{
		Name:        input.Name,
		Description: input.Description,
		Image:       input.Image,
		Category:    input.Category,
		UserId:      input.UserId,
	}

	output, err := p.productRepo.Create(ctx, product)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return entity.Product{}, ErrProductAlreadyExists
		}
		log.Error("ProductService.Create - p.productRepo.Create: %v", err)
		return entity.Product{}, ErrCannotCreateProduct
	}
	return output, nil
}

func (p ProductService) GetById(ctx context.Context, log *slog.Logger, input GetByIdProductInput) (
	entity.Product, error,
) {
	product, err := p.productRepo.GetById(ctx, input.Id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return entity.Product{}, ErrProductAlreadyExists
		}
		log.Error("ProductService.GetById - p.productRepo.GetById: %v", err)
		return entity.Product{}, ErrCannotCreateProduct
	}
	return product, nil
}
