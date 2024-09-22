package entity

type User struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
	City     string `db:"city"`
	Mode     string `db:"mode"`
}
