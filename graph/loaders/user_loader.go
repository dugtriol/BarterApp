package loaders

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	model2 "github.com/dugtriol/BarterApp/graph/model"
	"github.com/dugtriol/BarterApp/internal/controller"
	"github.com/dugtriol/BarterApp/internal/service"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// userReader reads Users from a database
type userReader struct {
	Log      *slog.Logger
	Services *service.Services
}

func (u *userReader) getUsers(ctx context.Context, userIDs []string) ([]*model2.User, []error) {
	output, errors := u.Services.User.GetUsers(ctx, userIDs)
	result := make([]*model2.User, len(output))

	for i, item := range output {
		var mode model2.UserMode
		err := mode.UnmarshalGQL(item.Mode)
		if err != nil {
			u.Log.Error("userReader - getUsers -  mode.UnmarshalGQL(item.Mode): ", err)
			return nil, []error{controller.ErrNotValid}
		}

		result[i] = &model2.User{
			ID:       item.Id,
			Name:     item.Name,
			Password: item.Password,
			Email:    item.Email,
			Phone:    item.Phone,
			City:     item.City,
			Mode:     mode,
		}
	}
	return result, errors
}

type Loaders struct {
	UserLoader *dataloadgen.Loader[string, *model2.User]
}

func NewLoaders(log *slog.Logger, services *service.Services) *Loaders {
	// define the data loader
	ur := &userReader{
		Log:      log,
		Services: services,
	}
	return &Loaders{
		UserLoader: dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(log *slog.Logger, services *service.Services) func(http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loader := NewLoaders(log, services)
			r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
			next.ServeHTTP(w, r)
		})
	}
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// GetUser returns single user by id efficiently
func GetUser(ctx context.Context, userID string) (*model2.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userID)
}

// GetUsers returns many users by ids efficiently
func GetUsers(ctx context.Context, userIDs []string) ([]*model2.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, userIDs)
}
