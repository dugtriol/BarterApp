package product

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dugtriol/BarterApp/internal/pkg/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type productDeleter interface {
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

func Delete(ctx context.Context, log *slog.Logger, deleter productDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.product.deleter.Delete"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to get id from url param", err.Error())
			render.JSON(w, r, response.Error("failed to get id from url param"))
			return
		}

		err = deleter.DeleteProduct(ctx, id)
		if err != nil {
			log.Error("failed to delete product", "id", id)
			render.JSON(w, r, response.Error("failed to delete product"))
			return
		}

		log.Info("product deleted")
	}
}
