package service

import (
	"context"
	"errors"
	"fmt"

	"log/slog"

	"github.com/dugtriol/BarterApp/graph/model"
	"github.com/dugtriol/BarterApp/internal/controller"
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

func (p *ProductService) Create(ctx context.Context, input CreateProductInput) (entity.Product, error) {
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

func (p *ProductService) GetById(ctx context.Context, log *slog.Logger, input GetByIdProductInput) (
	entity.Product, error,
) {
	product, err := p.productRepo.GetById(ctx, input.Id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return entity.Product{}, ErrProductAlreadyExists
		}
		log.Error("ProductService.GetById - p.productRepo.GetById: %v", err)
		return entity.Product{}, ErrCannotGetProduct
	}
	return product, nil
}

func (p *ProductService) All(ctx context.Context, limit, offset int) ([]entity.Product, error) {
	output, err := p.productRepo.All(ctx, limit, offset)
	if err != nil {
		log.Error(fmt.Sprintf("Service - ProductService - All: %v", err))
		return nil, ErrCannotGetProduct
	}
	return output, nil
}

func (p *ProductService) GetByUserId(ctx context.Context, limit, offset int, userId string) ([]*model.Product, error) {
	output, err := p.productRepo.GetByUserId(ctx, limit, offset, userId)
	if err != nil {
		log.Error(fmt.Sprintf("Service - ProductService - All: %v", err))
		return nil, ErrCannotGetProduct
	}

	result := make([]*model.Product, len(output))
	for i, item := range output {
		var category model.ProductCategory
		err = category.UnmarshalGQL(item.Category)
		if err != nil {
			log.Error("Resolvers.Product -  category.UnmarshalGQL(product.Category): ", err)
			return nil, controller.ErrNotValid
		}

		var status model.ProductStatus
		err = status.UnmarshalGQL(item.Status)
		if err != nil {
			log.Error("Resolvers.Product -  category.UnmarshalGQL(product.Status): ", err)
			return nil, controller.ErrNotValid
		}

		//var temp *model.Product
		result[i] = &model.Product{
			ID:          item.Id,
			Category:    category,
			Name:        item.Name,
			Description: item.Description,
			Image:       item.Image,
			Status:      status,
			CreatedBy:   item.UserId,
			CreatedAt:   item.CreatedAt.String(),
		}
		//log.Info(temp)
		//result = append(result, temp)
	}
	//log.Info(result)
	//log.Info(len(result))
	return result, nil
}
