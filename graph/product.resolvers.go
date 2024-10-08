package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.54

import (
	"context"

	"github.com/dugtriol/BarterApp/graph/loaders"
	"github.com/dugtriol/BarterApp/graph/model"
)

// CreatedBy is the resolver for the createdBy field.
func (r *productResolver) CreatedBy(ctx context.Context, obj *model.Product) (*model.User, error) {
	//output, err := p.Services.User.GetById(ctx, p.Log, service.UserGetByIdInput{Id: obj.CreatedBy})
	//if err != nil {
	//	p.Log.Error("productResolver - CreatedBy - ", err)
	//	return &model.User{}, err
	//}
	//var mode model.UserMode
	//err = mode.UnmarshalGQL(output.Mode)
	//if err != nil {
	//	p.Log.Error("Resolvers.User -  mode.UnmarshalGQL(output.Mode): ", err)
	//	return nil, controller.ErrNotValid
	//}
	//
	//result := model.User{
	//	ID:       output.Id,
	//	Name:     output.Name,
	//	Password: output.Password,
	//	Email:    output.Email,
	//	Phone:    output.Phone,
	//	City:     output.City,
	//	Mode:     mode,
	//}
	//
	//return &result, nil
	return loaders.GetUser(ctx, obj.CreatedBy)
}

// Product returns ProductResolver implementation.
func (r *Resolver) Product() ProductResolver { return &productResolver{r} }

type productResolver struct{ *Resolver }
