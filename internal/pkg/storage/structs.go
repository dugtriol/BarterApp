package storage

import "time"

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	City      string    `db:"city"`
	CreatedAt time.Time `db:"created_at"`
}
