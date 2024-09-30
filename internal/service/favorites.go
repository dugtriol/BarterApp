package service

import (
	"context"

	"github.com/dugtriol/BarterApp/graph/model"
	"github.com/dugtriol/BarterApp/internal/controller"
	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo"
	log "github.com/sirupsen/logrus"
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

func (p *FavoritesService) GetFavoritesByUserId(ctx context.Context, userId string) ([]*model.Favorites, error) {
	output, err := p.favoritesRepo.GetFavoritesByUserId(ctx, userId)
	if err != nil {
		log.Error("FavoritesService - GetFavoritesByUserId -  p.favoritesRepo.GetFavoritesByUserId: ", err)
		return nil, controller.ErrNotValid
	}
	return p.ParseFavoritesArray(output)
}

func (p *FavoritesService) ParseFavoritesArray(output []entity.Favorites) ([]*model.Favorites, error) {
	result := make([]*model.Favorites, len(output))
	for i, item := range output {
		//var temp *model.Product
		result[i] = &model.Favorites{
			ID:        item.Id,
			UserID:    item.UserId,
			ProductID: item.ProductId,
		}
		//log.Info(temp)
		//result = append(result, temp)
	}
	return result, nil
}

func (p *FavoritesService) DeleteIfProductDeleted(ctx context.Context, product_id string) (bool, error) {
	return p.favoritesRepo.DeleteIfProductDeleted(ctx, product_id)
}
