package user

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/dugtriol/BarterApp/internal/pkg/lib/api/response"
	"github.com/dugtriol/BarterApp/internal/pkg/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type userUpdateCity interface {
	UpdateUserCity(ctx context.Context, id uuid.UUID, newcity string) error
}

type requestUpdateCity struct {
	City string `json:"city" validate:"required"`
}

func UpdateCity(ctx context.Context, log *slog.Logger, update userUpdateCity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update_city.UpdateCity"
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

		var req requestUpdateCity

		// decode
		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, response.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded")

		// validator
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Invalid request")

			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		err = update.UpdateUserCity(ctx, id, req.City)
		if err != nil {
			log.Error("failed to update cityl", "id", id)
			render.JSON(w, r, response.Error("failed to update city"))
			return
		}

		log.Info("city has been updated", "id", id, "city", req.City)
	}
}
