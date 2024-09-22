package entity

import "time"

type Product struct {
	Id          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Image       string    `db:"image"`
	Status      string    `db:"status"`
	Category    string    `db:"category"`
	UserId      string    `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
}
