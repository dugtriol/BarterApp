package graph

import (
	"context"

	"github.com/dugtriol/BarterApp/internal/controller"
	"github.com/dugtriol/BarterApp/internal/controller/graph/model"
	"github.com/dugtriol/BarterApp/internal/service"
	"github.com/dugtriol/BarterApp/pkg/hasher"
	"github.com/dugtriol/BarterApp/pkg/middleware"
)

type mutationResolver struct{ *Resolver }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Create is the resolver for the Create field.
func (r *mutationResolver) Register(ctx context.Context, input *model.CreateUserInput) (*model.AuthResponse, error) {
	isValid := validation(ctx, input)
	if !isValid {
		return nil, controller.ErrInput
	}

	_, err := r.Services.Auth.GetByEmail(ctx, r.Log, service.UserGetByEmailInput{Email: input.Email})
	if err == nil {
		r.Log.Error("Resolvers.Create - r.Services.Auth.GetByEmail: email already exists")
		return nil, controller.ErrAlreadyExists
	}
	output, err := r.Services.Auth.Register(
		ctx, service.AuthRegisterInput{
			Name:     input.Name,
			Email:    input.Email,
			Phone:    input.Phone,
			Password: input.Password,
			City:     input.City,
			Mode:     input.Mode.String(),
		},
	)
	if err != nil {
		r.Log.Error("Resolvers.Create - r.Services.Auth.Register: ", err)
		return nil, err
	}

	authToken, err := r.Services.Auth.GenToken(output.Id)
	if err != nil {
		r.Log.Error("Resolvers.Create - r.Services.Auth.GenerateToken: ", err)
		return nil, err
	}

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

	response := model.AuthResponse{
		AuthToken: authToken,
		User:      &result,
	}
	return &response, nil
}

func (r *mutationResolver) Login(ctx context.Context, input *model.LoginInput) (*model.AuthResponse, error) {
	isValid := validation(ctx, input)
	if !isValid {
		return nil, controller.ErrInput
	}

	output, err := r.Services.Auth.GetByEmail(ctx, r.Log, service.UserGetByEmailInput{Email: input.Email})
	if err != nil {
		r.Log.Error("Resolvers.Login - r.Services.Auth.GetByEmail: ", err)
		return nil, controller.ErrAlreadyExists
	}
	err = hasher.CheckPassword(input.Password, output.Password)
	if err != nil {
		r.Log.Error("Resolvers.Login - hasher.CheckPassword: ", err)
		return nil, controller.ErrInvalidPassword
	}

	authToken, err := r.Services.Auth.GenToken(output.Id)
	if err != nil {
		r.Log.Error("Resolvers.Create - r.Services.Auth.GenerateToken: ", err)
		return nil, controller.ErrNotValid
	}

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

	response := model.AuthResponse{
		AuthToken: authToken,
		User:      &result,
	}
	return &response, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input *model.CreateProductInput) (*model.Product, error) {
	currentUser, err := middleware.GetCurrentUserFromCTX(ctx)
	if err != nil {
		r.Log.Error("Resolvers.Product -  middleware.GetCurrentUserFromCTX: no user in context")
		return nil, controller.ErrNotAuthenticated
	}
	product, err := r.Services.Product.Create(
		ctx, service.CreateProductInput{
			Category:    input.Category.String(),
			Name:        input.Name,
			Description: input.Description,
			Image:       input.Image,
			UserId:      currentUser.Id,
		},
	)
	if err != nil {
		r.Log.Error("Resolvers.Product -  r.Services.Product.Create: ", err)
		return nil, controller.ErrNotCreated
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
		CreatedBy:   product.UserId,
	}
	return &result, nil
}
