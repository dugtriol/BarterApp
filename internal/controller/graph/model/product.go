package model

type Product struct {
	ID          string          `json:"id"`
	Category    ProductCategory `json:"category"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Status      ProductStatus   `json:"status"`
	CreatedBy   string          `json:"createdBy"`
	CreatedAt   string          `json:"createdAt"`
}
