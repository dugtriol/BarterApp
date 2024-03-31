package product

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/dugtriol/BarterApp/internal/pkg/lib/api/response"
	"github.com/dugtriol/BarterApp/internal/pkg/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type requestSaveProduct struct {
	Name        string    `json:"name" validate:"required"`
	IdOwner     uuid.UUID `json:"id_owner" validate:"required,uuid"`
	Description string    `json:"description" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	City        string    `json:"city" validate:"required"`
}

type saveResponse struct {
	response.Response
	Id uuid.UUID `json:"id"`
}

type productSaver interface {
	SaveProduct(ctx context.Context, idOwner uuid.UUID, name, description, image, city string) (uuid.UUID, error)
}

func Save(ctx context.Context, log *slog.Logger, saver productSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.product.save.Save"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		//id_owner, err := uuid.Parse(chi.URLParam(r, "id_owner"))
		//if err != nil {
		//	log.Error("failed to get id_owner from url param")
		//	render.JSON(w, r, response.Error("failed to get id_owner from url param"))
		//	return
		//}

		// TODO: проверка что айди пользователя существует

		var req requestSaveProduct

		// decode
		err := render.DecodeJSON(r.Body, &req)
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

		id, err := saver.SaveProduct(ctx, req.IdOwner, req.Name, req.Description, req.Image, req.City)
		if err != nil {
			log.Error("failed to add product", err.Error())
			render.JSON(w, r, response.Error("failed to add product"))
			return
		}

		log.Info("product added", slog.String("id", id.String()))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	render.JSON(
		w, r, saveResponse{
			Response: response.OK(),
			Id:       id,
		},
	)
}
