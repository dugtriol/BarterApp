package graph

import (
	"context"

	"github.com/dugtriol/BarterApp/internal/controller/graph/model"
)

type userResolver struct{ *Resolver }

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

func (u userResolver) PostedProducts(ctx context.Context, obj *model.User) ([]*model.Product, error) {
	// TODO DataLoader
	panic("")
}
