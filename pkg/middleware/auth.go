package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/service"
	"github.com/pkg/errors"
)

const CurrentUserKey = "currentUser"

func AuthMiddleware(ctx context.Context, log *slog.Logger, authService service.Auth) func(http.Handler) http.Handler {
	log.Info("start auth middleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				log.Info("parse token")
				token, err := authService.ParseToken(r)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				//claims, ok := token.Claims.(*service.TokenClaims)

				//log.Info(fmt.Sprintf("!ok =  %v || !token.Valid = %v", !ok, !token.Valid))
				//log.Info(fmt.Sprintf("%v", token))
				if !ok || !token.Valid {
					next.ServeHTTP(w, r)
					return
				}
				//log.Info(fmt.Sprintf("getbyid %s", claims["UserId"]))
				user, err := authService.GetById(ctx, log, service.UserGetByIdInput{Id: claims["UserId"].(string)})
				if err != nil {
					log.Error(err.Error())
					next.ServeHTTP(w, r)
					return
				}
				//log.Info("context.WithValue")
				ctx := context.WithValue(r.Context(), CurrentUserKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func GetCurrentUserFromCTX(ctx context.Context) (*entity.User, error) {
	errNoUserInContext := errors.New("no user in context")
	//log.Info(fmt.Sprintf("ctx.Value(CurrentUserKey): %s",ctx.Value(CurrentUserKey)))
	if ctx.Value(CurrentUserKey) == nil {
		return nil, errNoUserInContext
	}

	user, ok := ctx.Value(CurrentUserKey).(entity.User)
	//log.Info(fmt.Sprintf("ctx.Value(CurrentUserKey).(*entity.User) %v",user))
	if !ok || user.Id == "" {
		return nil, errNoUserInContext
	}

	return &user, nil
}
