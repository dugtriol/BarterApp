package user

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dugtriol/BarterApp/internal/pkg/lib/api/response"
	"github.com/dugtriol/BarterApp/internal/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type userGetByID interface {
	GetByID(ctx context.Context, id uuid.UUID) (*storage.User, error)
}

func Get(ctx context.Context, log *slog.Logger, get userGetByID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get.New"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		// strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Error("failed to get id from url param", err.Error())
			render.JSON(w, r, response.Error("failed to get id from url param"))
			return
		}

		res, err := get.GetByID(ctx, id)
		if err != nil {
			log.Error("failed to redirect url", "id", id)
			render.JSON(w, r, response.Error("failed to redirect url"))
			return
		}

		log.Info("got user", "user", res)
		w.Write([]byte(fmt.Sprintf("user: %s, %s, %s, %s", res.Name, res.Lastname, res.Email, res.City, res.CreatedAt)))
	}
}
