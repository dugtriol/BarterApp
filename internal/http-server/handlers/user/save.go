package user

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/dugtriol/BarterApp/internal/pkg/lib/api/response"
	"github.com/dugtriol/BarterApp/internal/pkg/lib/logger/sl"
	"github.com/dugtriol/BarterApp/internal/pkg/lib/password"
	`github.com/dugtriol/BarterApp/internal/pkg/mail`
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RequestSaveUser struct {
	Name     string `json:"name" validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	City     string `json:"city" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type saveResponse struct {
	response.Response
	Id uuid.UUID `json:"id"`
}

type userSaver interface {
	SaveUser(ctx context.Context, name, lastname, email, city, password string) (uuid.UUID, error)
	//SendVerifyEmail(ctx context.Context, log *slog.Logger, name, emailPath string) error
}

func Save(ctx context.Context, log *slog.Logger, saver userSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		//_, m, err := jwtauth.FromContext(r.Context())
		//if err != nil {
		//	return
		//}

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestSaveUser

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

		hashPassword, err := password.HashPassword(req.Password)

		if err != nil {
			return
		}
		id, err := saver.SaveUser(ctx, req.Name, req.Lastname, req.Email, req.City, hashPassword)
		if err != nil {
			log.Error("failed to add user", err.Error())
			render.JSON(w, r, response.Error("failed to add url"))
			return
		}

		log.Info("user added", slog.String("id", id.String()))

		err = mail.SendVerifyEmail(ctx, log, req.Name, req.Email)
		if err != nil {
			log.Error("failed to send email to user", err.Error())
			render.JSON(w, r, response.Error("failed to send email to user"))
			return
		}
		log.Info("email verification has been sent")
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
