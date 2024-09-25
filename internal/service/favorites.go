package service

import (
	"context"

	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo"
)

type FavoritesService struct {
	favoritesRepo repo.Favorites
}

func NewFavoritesService(
	favoritesRepo repo.Favorites,
) *FavoritesService {
	return &FavoritesService{
		favoritesRepo: favoritesRepo,
	}
}

func (p *FavoritesService) Add(ctx context.Context, input entity.Favorites) (string, error) {
	return p.favoritesRepo.Add(ctx, input)
}

func (p *FavoritesService) Delete(ctx context.Context, input entity.Favorites) (bool, error) {
	return p.favoritesRepo.Delete(ctx, input)
}
