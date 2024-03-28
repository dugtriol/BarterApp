package save

import (
	`context`
	`errors`
	`io`
	`log/slog`
	`net/http`

	`github.com/dugtriol/BarterApp/internal/pkg/lib/api/response`
	`github.com/dugtriol/BarterApp/internal/pkg/lib/logger/sl`
	`github.com/go-chi/chi/v5/middleware`
	`github.com/go-chi/render`
	`github.com/go-playground/validator/v10`
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	City     string `json:"city" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	response.Response
	Id int64 `json:"id"`
}

type UserSaver interface {
	SaveUser(ctx context.Context, name, email, city, password string) (int64, error)
}

func New(ctx context.Context, log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

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

		//// alias
		//alias := req.Alias
		//if alias == "" {
		//	alias = random.NewRandomString(aliasLength)
		//}

		id, err := userSaver.SaveUser(ctx, req.Username, req.Email, req.City, req.Password)
		//if errors.Is(err, storage.ErrURLExists) {
		//	log.Info("url already exists", slog.String("url", req.URL))
		//	render.JSON(w, r, response.Error("url already exists"))
		//	return
		//}
		if err != nil {
			log.Error("failed to add user", err.Error())
			render.JSON(w, r, response.Error("failed to add url"))
			return
		}

		log.Info("user added", slog.Int64("id", id))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int64) {
	render.JSON(
		w, r, Response{
			Response: response.OK(),
			Id:       id,
		},
	)
}
