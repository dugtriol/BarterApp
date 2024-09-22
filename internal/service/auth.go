package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dugtriol/BarterApp/internal/entity"
	"github.com/dugtriol/BarterApp/internal/repo"
	"github.com/dugtriol/BarterApp/internal/repo/repoerrs"
	"github.com/dugtriol/BarterApp/pkg/hasher"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type AuthService struct {
	userRepo repo.User
	//passwordHasher hasher.PasswordHasher
	signKey  string
	tokenTTL time.Duration
}

func NewAuthService(
	userRepo repo.User, signKey string, tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		//passwordHasher: passwordHasher,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, input AuthRegisterInput) (entity.User, error) {
	//log.Info(fmt.Sprintf("Service - AuthService - Create"))
	password, err := hasher.HashPassword(input.Password)
	if err != nil {
		log.Errorf("AuthService.Register - s.passwordHasher.HashPassword: %v", err)
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
		log.Errorf("AuthService.Register - c.userRepo.Register: %v", err)
		return entity.User{}, ErrCannotCreateUser
	}
	return output, nil
}

func (s *AuthService) GetById(ctx context.Context, log *slog.Logger, input UserGetByIdInput) (entity.User, error) {
	//log.Info(fmt.Sprintf("Service - AuthService - GetById"))
	user, err := s.userRepo.GetUserById(ctx, input.Id)
	if err != nil {
		if err != nil {
			if err == repoerrs.ErrAlreadyExists {
				return entity.User{}, ErrUserAlreadyExists
			}
			log.Error(fmt.Sprintf("Service - UserService - Create: %v", err))
			return entity.User{}, ErrCannotCreateUser
		}
	}
	return user, nil
}

//func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
//	// get user from DB
//	user, err := s.userRepo.GetUserByUsernameAndPassword(ctx, input.Username, s.passwordHasher.Hash(input.Password))
//	if err != nil {
//		if errors.Is(err, repoerrs.ErrNotFound) {
//			return "", ErrUserNotFound
//		}
//		log.Errorf("AuthService.GenerateToken: cannot get user: %v", err)
//		return "", ErrCannotGetUser
//	}
//
//	// generate token
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
//			IssuedAt:  time.Now().Unix(),
//		},
//		UserId: user.Id,
//	})
//
//	// sign token
//	tokenString, err := token.SignedString([]byte(s.signKey))
//	if err != nil {
//		log.Errorf("AuthService.GenerateToken: cannot sign token: %v", err)
//		return "", ErrCannotSignToken
//	}
//
//	return tokenString, nil
//}
//
//func (s *AuthService) ParseToken(accessToken string) (int, error) {
//	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//
//		return []byte(s.signKey), nil
//	})
//
//	if err != nil {
//		return 0, ErrCannotParseToken
//	}
//
//	claims, ok := token.Claims.(*TokenClaims)
//	if !ok {
//		return 0, ErrCannotParseToken
//	}
//
//	return claims.UserId, nil
//}
