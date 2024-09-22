package service

import (
	"github.com/dugtriol/BarterApp/internal/repo"
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
