package get

import (
	`context`
	`fmt`
	`log/slog`
	`net/http`
	`strconv`

	`github.com/dugtriol/BarterApp/internal/pkg/lib/api/response`
	`github.com/dugtriol/BarterApp/internal/pkg/storage`
	`github.com/go-chi/chi/v5`
	`github.com/go-chi/chi/v5/middleware`
	`github.com/go-chi/render`
)

type UserGet interface {
	GetByID(ctx context.Context, id int64) (*storage.User, error)
}

func New(ctx context.Context, log *slog.Logger, userGet UserGet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get.New"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Error("failed to get id from url param", err.Error())
			render.JSON(w, r, response.Error("failed to get id from url param"))
			return
		}
		//if id == "" {
		//	log.Error("id is empty")
		//	render.JSON(w, r, response.Error("id is empty"))
		//	return
		//}

		// find url via id in DB
		res, err := userGet.GetByID(ctx, id)
		//if errors.Is(err, storage.ErrURLNotFound) {
		//	log.Error("url not found", "id", id)
		//	render.JSON(w, r, response.Error("url not found"))
		//	return
		//}

		if err != nil {
			log.Error("failed to redirect url", "id", id)
			render.JSON(w, r, response.Error("failed to redirect url"))
			return
		}

		// redirect
		log.Info("got user", "url", res)
		w.Write([]byte(fmt.Sprintf("user: %s, %s, %s, %s", res.Username, res.Email, res.City, res.CreatedAt)))
	}
}
