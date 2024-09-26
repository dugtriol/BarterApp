package model

import "github.com/dugtriol/BarterApp/graph/scalar"

type Product struct {
	ID          string          `json:"id"`
	Category    ProductCategory `json:"category"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Status      ProductStatus   `json:"status"`
	CreatedBy   string          `json:"createdBy"`
	CreatedAt   scalar.DateTime `json:"createdAt"`
}
