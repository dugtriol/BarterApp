package postgresql

import (
	"context"
	"database/sql"

	"github.com/dugtriol/BarterApp/internal/pkg/db"
	"github.com/dugtriol/BarterApp/internal/pkg/repository"
)

type UserRepo struct {
	db *db.Database
}

func NewArticles(database *db.Database) *UserRepo {
	return &UserRepo{db: database}
}

func (r *UserRepo) Add(ctx context.Context, article *repository.User) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO articles(name, rating) VALUES($1, $2) RETURNING id;`,
		article.Name,
		article.Rating,
	).Scan(&id)
	return id, err
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*repository.User, error) {
	var a repository.User
	err := r.db.Get(ctx, &a, "SELECT id,name,rating,created_at FROM articles WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, repository.ErrObjectNoFound
	}

	return &a, nil
}
