package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dugtriol/BarterApp/internal/pkg/db"
	"github.com/dugtriol/BarterApp/internal/pkg/storage"
	"github.com/google/uuid"
)

type Storage struct {
	db *db.Database
}

func New(database *db.Database) *Storage {
	return &Storage{db: database}
}

func (r *Storage) SaveUser(ctx context.Context, name, lastname, email, city, password string) (uuid.UUID, error) {
	id := uuid.New()

	_, err := r.db.Exec(
		ctx,
		`INSERT INTO users(id, name, lastname, email, city, password) VALUES($1, $2, $3, $4,$5,$6);`,
		id,
		name,
		lastname,
		email,
		city,
		password,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return id, err
}

func (r *Storage) GetByID(ctx context.Context, id uuid.UUID) (*storage.User, error) {
	var a storage.User
	err := r.db.Get(ctx, &a, "SELECT name, lastname,email,city,created_at FROM users WHERE id=$1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrObjectNoFound
	}

	return &a, nil
}

func (r *Storage) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrObjectNoFound
	}

	return nil
}

func (r *Storage) UpdateUserPassword(ctx context.Context, id uuid.UUID, oldpassword, newpassword string) error {
	_, err := r.db.Exec(
		ctx,
		"UPDATE users SET password = $1 WHERE id = $2 AND password = $3",
		newpassword,
		id,
		oldpassword,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Storage) UpdateUserCity(ctx context.Context, id uuid.UUID, newcity string) error {
	_, err := r.db.Exec(
		ctx, "UPDATE users SET city = $1 WHERE id = $2",
		newcity,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Storage) SaveProduct(
	ctx context.Context, id_owner uuid.UUID, name, description, image, city string,
) (uuid.UUID, error) {
	id := uuid.New()
	status := "available"

	_, err := r.db.Exec(
		ctx,
		`INSERT INTO products(id, id_owner, name, description, image, city, status) VALUES($1, $2, $3, $4,$5, $6,$7);`,
		id,
		id_owner,
		name,
		description,
		image,
		city,
		status,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return id, err
}

func (r *Storage) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrObjectNoFound
	}
	return nil
}

func (r *Storage) GetUserByEmail(ctx context.Context, email string) (*storage.User, error) {
	var a storage.User
	err := r.db.Get(ctx, &a, "SELECT id,name, lastname,email,password,city,created_at FROM users WHERE email=$1", email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrObjectNoFound
	}

	return &a, nil
}


