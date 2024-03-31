package storage

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	IdOwner   uuid.UUID `db:"id_owner"`
	Name      string    `db:"name"`
	Lastname  string    `db:"lastname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	City      string    `db:"city"`
	CreatedAt time.Time `db:"created_at"`
}
