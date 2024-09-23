package graph

import (
	"context"
	"fmt"

	"github.com/dugtriol/BarterApp/internal/controller"
	"github.com/dugtriol/BarterApp/internal/controller/graph/model"
	"github.com/dugtriol/BarterApp/internal/service"
	"github.com/dugtriol/BarterApp/pkg/middleware"
)

type queryResolver struct{ *Resolver }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// TotalUsers is the resolver for the totalUsers field.
func (r *queryResolver) TotalUsers(ctx context.Context) (int, error) {
	return 1, nil
}

// AllUsers is the resolver for the allUsers field.
func (r *queryResolver) AllUsers(ctx context.Context) ([]*model.User, error) {
	return nil, nil
}

// TotalProducts is the resolver for the totalProducts field.
func (r *queryResolver) TotalProducts(ctx context.Context) (int, error) {
	panic(fmt.Errorf("not implemented: TotalProducts - totalProducts"))
}

// AllProducts is the resolver for the allProducts field.
func (r *queryResolver) AllProducts(
	ctx context.Context, category *model.ProductCategory, first *int, start *int,
) ([]*model.Product, error) {
	panic(fmt.Errorf("not implemented: AllProducts - allProducts"))
}

// User is the resolver for the User field.
func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	output, err := middleware.GetCurrentUserFromCTX(ctx)
	if err != nil {
		r.Log.Error("Resolvers.Product -  middleware.GetCurrentUserFromCTX: no user in context")
		return nil, controller.ErrNotAuthenticated
	}

	//output, err := r.Services.Auth.GetById(ctx, r.Log, service.UserGetByIdInput{Id: id})
	//if err != nil {
	//	r.Log.Error("Resolvers.User -  r.Services.Auth.GetById: ", err)
	//	return nil, controller.ErrNotFound
	//}

	var mode model.UserMode
	err = mode.UnmarshalGQL(output.Mode)
	if err != nil {
		r.Log.Error("Resolvers.User -  mode.UnmarshalGQL(output.Mode): ", err)
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

// Product is the resolver for the Product field.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	product, err := r.Services.Product.GetById(ctx, r.Log, service.GetByIdProductInput{Id: id})
	if err != nil {
		r.Log.Error("Resolvers.Product -  r.Services.Product.GetById: ", err)
		return nil, controller.ErrNotFound
	}

	var category model.ProductCategory
	err = category.UnmarshalGQL(product.Category)
	if err != nil {
		r.Log.Error("Resolvers.Product -  category.UnmarshalGQL(product.Category): ", err)
		return nil, controller.ErrNotValid
	}

	var status model.ProductStatus
	err = status.UnmarshalGQL(product.Status)
	if err != nil {
		r.Log.Error("Resolvers.Product -  status.UnmarshalGQL(product.Status): ", err)
		return nil, controller.ErrNotValid
	}

	result := model.Product{
		ID:          product.Id,
		Category:    category,
		Name:        product.Name,
		Description: product.Description,
		Image:       product.Image,
		Status:      status,
		CreatedAt:   product.CreatedAt.String(),
	}
	return &result, nil

}
