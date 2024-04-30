package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"jwt-mongo-auth/internal/entity/token"
	"jwt-mongo-auth/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	log      *logrus.Logger
	userRepo repository.User
}

func NewToken(log *logrus.Logger, userRepo repository.User) *Service {
	return &Service{
		log:      log,
		userRepo: userRepo,
	}
}

func(s *Service) hashRefreshToken(token string) ([]byte, error) {
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedRefreshToken, nil
}

func (s *Service) GenerateTokens(guide string) (token.Token, error) {
	accessTokenClaims := token.Claims{
		GUID: guide,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(token.JwtSecret)
	if err != nil {
		return token.Token{}, err
	}

	refreshTokenBytes := make([]byte, 64)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return token.Token{}, err
	}

	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)

	hashedRefreshToken,err := s.hashRefreshToken(refreshToken)
	if err != nil{
		return token.Token{},err
	}


	s.userRepo.SaveRefreshToken(context.Background(),hashedRefreshToken,guide)

	return token.Token{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
	}, nil

}
