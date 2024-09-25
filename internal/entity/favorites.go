package entity

type Favorites struct {
	Id        string `db:"id"`
	UserId    string `db:"user_id"`
	ProductId string `db:"product_id"`
}
