package entity

type Transaction struct {
	Id             string `yaml:"id"`
	Owner          string `yaml:"owner"`
	Buyer          string `yaml:"buyer"`
	ProductIdBuyer string `yaml:"product_id_buyer"`
	ProductIdOwner string `yaml:"product_id_owner"`
	Status         string `yaml:"status"`
	CreatedAt      string `yaml:"created_at"`
	UpdatedAt      string `yaml:"updated_at"`
	Shipping       string `yaml:"shipping"`
	Address        string `yaml:"address"`
}
