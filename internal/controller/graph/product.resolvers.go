package graph

import (
	"context"

	"github.com/dugtriol/BarterApp/internal/controller"
	"github.com/dugtriol/BarterApp/internal/controller/graph/model"
	"github.com/dugtriol/BarterApp/internal/service"
)

type productResolver struct{ *Resolver }

func (r *Resolver) Product() ProductResolver {
	return &productResolver{r}
}

func (p *productResolver) CreatedBy(ctx context.Context, obj *model.Product) (*model.User, error) {
	output, err := p.Services.Auth.GetById(ctx, p.Log, service.UserGetByIdInput{Id: obj.CreatedBy})
	if err != nil {
		p.Log.Error("productResolver - CreatedBy - ", err)
		return &model.User{}, err
	}
	var mode model.UserMode
	err = mode.UnmarshalGQL(output.Mode)
	if err != nil {
		p.Log.Error("Resolvers.User -  mode.UnmarshalGQL(output.Mode): ", err)
		return nil, controller.ErrNotValid
	}

	result := model.User{
		ID:       output.Id,
		Name:     output.Name,
		Password: output.Password,
		Email:    output.Email,
		Phone:    output.Phone,
		City:     output.City,
		Mode:     mode,
	}

	return &result, nil
}
