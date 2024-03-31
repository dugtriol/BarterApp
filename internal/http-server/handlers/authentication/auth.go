package authentication

import (
	"context"
	"errors"
	`fmt`
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/dugtriol/BarterApp/internal/pkg/lib/api/response"
	"github.com/dugtriol/BarterApp/internal/pkg/lib/logger/sl"
	`github.com/dugtriol/BarterApp/internal/pkg/storage`
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	`golang.org/x/crypto/bcrypt`
)

type signInInput struct {
	Email    string `json:"email" binding:"required,mail"`
	Password string `json:"password" binding:"required"`
}

func MakeToken(email string, tokenAuth *jwtauth.JWTAuth) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"email": email})
	return tokenString
}

type userGetByEmail interface {
	GetUserByEmail(ctx context.Context, email string) (*storage.User, error)
}

func LogIn(ctx context.Context, log *slog.Logger, tokenAuth *jwtauth.JWTAuth, get userGetByEmail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signInInput
		const op = "handlers.authentication.auth.LogIn"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

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

		// TODO: проверка email и пароля
		u, err := get.GetUserByEmail(ctx, req.Email)
		if err != nil {
			log.Error("Not found email in database", err.Error())
			render.JSON(w, r, response.Error("Not found email inn database"))
			return
		}
		log.Info("%s %s", u.Password, req.Password)
		if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
			log.Error("Invalid password", err.Error())
			render.JSON(w, r, response.Error("invalid password"))
			return
		}

		token := MakeToken(req.Email, tokenAuth)
		log.Info("make token")
		http.SetCookie(
			w, &http.Cookie{
				HttpOnly: true,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				SameSite: http.SameSiteLaxMode,
				Name:     "jwt",
				Value:    token,
			},
		)
		log.Info("End og LogIn method")
		w.Write([]byte(fmt.Sprintf("Log In! %s", tokenAuth)))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func LogOut(ctx context.Context, log *slog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			HttpOnly: true,
			MaxAge: -1, // Delete the cookie.
			SameSite: http.SameSiteLaxMode,
			// Uncomment below for HTTPS:
			// Secure: true,
			Name:  "jwt",
			Value: "",
		})
		log.Info("Log out")
		//http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
