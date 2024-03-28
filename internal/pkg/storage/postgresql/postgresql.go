package postgresql

import (
	"context"
	"database/sql"

	"github.com/dugtriol/BarterApp/internal/pkg/db"
	"github.com/dugtriol/BarterApp/internal/pkg/storage"
)

type Storage struct {
	db *db.Database
}

func New(database *db.Database) *Storage {
	return &Storage{db: database}
}

func (r *Storage) SaveUser(ctx context.Context, name, email, city, password string) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO users(username, email, city, password) VALUES($1, $2, $3, $4) RETURNING id;`,
		name,
		email,
		city,
		password,
	).Scan(&id)
	return id, err
}

func (r *Storage) GetByID(ctx context.Context, id int64) (*storage.User, error) {
	var a storage.User
	err := r.db.Get(ctx, &a, "SELECT username,email,city,created_at FROM users WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, storage.ErrObjectNoFound
	}

	return &a, nil
}
