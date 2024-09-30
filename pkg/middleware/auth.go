package middleware

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/dgrijalva/jwt-go"
	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/service"
	"github.com/pkg/errors"
)

const CurrentUserKey = "currentUser"

//func AuthMiddleware(ctx context.Context, log *slog.Logger, authService service.User) func(http.Handler) http.Handler {
//	//log.Info("start auth middleware")
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(
//			func(w http.ResponseWriter, r *http.Request) {
//				//log.Info("parse token")
//				operationContext := graphql.GetOperationContext(ctx)
//				name := operationContext.OperationName
//				log.Info(fmt.Sprintf("Auth middleware - name -  %v", name))
//				if name == "Login" || name == "Register" {
//					next.ServeHTTP(w, r)
//				}
//
//				token, err := authService.ParseTokenFromRequest(r)
//				if err != nil {
//					next.ServeHTTP(w, r)
//					return
//				}
//
//				claims, ok := token.Claims.(jwt.MapClaims)
//				//claims, ok := token.Claims.(*service.TokenClaims)
//
//				//log.Info(fmt.Sprintf("!ok =  %v || !token.Valid = %v", !ok, !token.Valid))
//				//log.Info(fmt.Sprintf("%v", token))
//				if !ok || !token.Valid {
//					next.ServeHTTP(w, r)
//					return
//				}
//				//log.Info(fmt.Sprintf("claims %v", claims))
//				//log.Info(fmt.Sprintf("getbyid %s", claims["jti"]))
//				user, err := authService.GetById(ctx, log, service.UserGetByIdInput{Id: claims["jti"].(string)})
//				if err != nil {
//					log.Error(err.Error())
//					next.ServeHTTP(w, r)
//					return
//				}
//				//log.Info("context.WithValue")
//				ctx := context.WithValue(r.Context(), CurrentUserKey, user)
//				next.ServeHTTP(w, r.WithContext(ctx))
//			},
//		)
//	}
//}

func AuthMiddleware(ctx context.Context, log *slog.Logger, authService service.User, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	operationContext := graphql.GetOperationContext(ctx)
	authHeader := operationContext.Headers.Get("Authorization")
	if authHeader == "" {
		log.Info("AuthMiddleware - operationContext.Headers.Get(Authorization)")
		return nil, errors.New("failed to get currrent user")
	}
	log.Info(fmt.Sprintf("Authheader: %s", authHeader))
	log.Info(fmt.Sprintf("obj: %v", obj))

	parseToken, err := authService.ParseToken(authHeader)
	if err != nil {
		log.Info("AuthMiddleware - services.User.ParseToken(token) - ", err)
		return nil, errors.New("failed to parse token services.User.ParseToken(token)")
	}
	claims, ok := parseToken.Claims.(jwt.MapClaims)

	if !ok || !parseToken.Valid {
		log.Info(fmt.Sprintf("AuthMiddleware - !ok = %v || !parseToken.Valid = %v", !ok, !parseToken.Valid))
		log.Info(fmt.Sprintf("parseToken: %v claims: %v", parseToken, parseToken.Claims.(jwt.MapClaims)))
		return nil, errors.New("failed to get claims from token")
	}

	user, err := authService.GetById(ctx, log, service.UserGetByIdInput{Id: claims["jti"].(string)})
	if err != nil {
		log.Info("AuthMiddleware - services.User.GetById", ok)
		return nil, errors.New("failed to get user")
	}

	// put it in context
	ctxNew := context.WithValue(ctx, CurrentUserKey, user)
	return next(ctxNew)
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
