package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/dugtriol/BarterApp/graph/model"
	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo"
	"github.com/dugtriol/BarterApp/internal/repo/repoerrs"
	"github.com/dugtriol/BarterApp/pkg/hasher"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	userRepo repo.User
	//passwordHasher hasher.PasswordHasher
	signKey  string
	tokenTTL time.Duration
}

func NewUserService(
	userRepo repo.User, signKey string, tokenTTL time.Duration,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		//passwordHasher: passwordHasher,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (s *UserService) Register(ctx context.Context, input AuthRegisterInput) (entity.User, error) {
	//log.Info(fmt.Sprintf("Service - UserService - Create"))
	password, err := hasher.HashPassword(input.Password)
	if err != nil {
		log.Errorf("UserService.Register - s.passwordHasher.HashPassword: %v", err)
	}
	user := entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: password,
		City:     input.City,
		Mode:     input.Mode,
	}

	output, err := s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return entity.User{}, ErrUserAlreadyExists
		}
		log.Errorf("UserService.Register - c.userRepo.Register: %v", err)
		return entity.User{}, ErrCannotCreateUser
	}
	return output, nil
}

func (s *UserService) GetById(ctx context.Context, log *slog.Logger, input UserGetByIdInput) (entity.User, error) {
	//log.Info(fmt.Sprintf("Service - UserService - GetById"))
	user, err := s.userRepo.GetById(ctx, input.Id)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return entity.User{}, ErrUserAlreadyExists
		}
		log.Error(fmt.Sprintf("Service - UserService - GetById: %v", err))
		return entity.User{}, ErrCannotCreateUser
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, log *slog.Logger, input UserGetByEmailInput) (entity.User, error) {
	//log.Info(fmt.Sprintf("Service - UserService - GetById"))
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return entity.User{}, ErrUserAlreadyExists
		}
		log.Error(fmt.Sprintf("Service - UserService - Create: %v", err))
		return entity.User{}, ErrCannotCreateUser
	}
	return user, nil
}

func (s *UserService) GenToken(id string) (*model.AuthToken, error) {
	expiredAt := time.Now().Add(time.Hour * 24 * 7) // a week

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiredAt.Unix(),
		Id:        id,
		IssuedAt:  time.Now().Unix(),
	})

	accessToken, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return nil, err
	}

	return &model.AuthToken{
		AccessToken: accessToken,
		ExpiredAt:   expiredAt,
	}, nil
}

var authHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromToken,
}

func stripBearerPrefixFromToken(token string) (string, error) {
	bearer := "BEARER"

	if len(token) > len(bearer) && strings.ToUpper(token[0:len(bearer)]) == bearer {
		return token[len(bearer)+1:], nil
	}

	return token, nil
}

var authExtractor = &request.MultiExtractor{
	authHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

func (s *UserService) ParseToken(r *http.Request) (*jwt.Token, error) {
	jwtToken, err := request.ParseFromRequest(
		r, authExtractor, func(token *jwt.Token) (interface{}, error) {
			t := []byte(s.signKey)
			//log.Info(fmt.Sprintf("ParseToken -  %v", t))
			return t, nil
		},
	)
	if err != nil {
		log.Errorf("UserService - parseToken: ", err)
		return nil, err
	}

	return jwtToken, nil
}

func (s *UserService) GetUsers(ctx context.Context, userIDs []string) ([]*entity.User, []error) {
	return s.userRepo.GetUsers(ctx, userIDs)
}
