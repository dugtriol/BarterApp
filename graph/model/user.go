package model

type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	City     string   `json:"city"`
	Mode     UserMode `json:"mode"`
}
