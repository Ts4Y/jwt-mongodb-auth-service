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

func NewService(log *logrus.Logger, userRepo repository.User) *Service {
	return &Service{
		log:      log,
		userRepo: userRepo,
	}
}

func (s *Service) HashRefreshToken(token string) ([]byte, error) {
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		s.log.Errorln("Не удалось захэшировать refresh token")
		return nil, err
	}

	return hashedRefreshToken, nil
}

func (s *Service) GenerateTokens(guid string) (token.Token, error) {
	accessTokenClaims := token.Claims{
		GUID: guid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(token.JwtSecret)
	if err != nil {
		s.log.Errorln("Не удалось создать jwt token", err)
		return token.Token{}, err
	}

	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		s.log.Errorln("Не удалось создать refresh token", err)
		return token.Token{}, err
	}

	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)

	hashedRefreshToken, err := s.HashRefreshToken(refreshToken)
	if err != nil {
		s.log.Errorln("Не удалось захэшировать refresh token", err)
		return token.Token{}, err
	}

	err = s.userRepo.SaveRefreshToken(context.Background(), hashedRefreshToken, guid)

	if err != nil {
		s.log.Errorln("Не удалось сохранить токены", err)
		return token.Token{}, err
	}

	return token.Token{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
	}, nil

}

func (s *Service) UpdateTokens(ctx context.Context, guid string, reftoken string) (token.Token, error) {

	hrefTok, err := s.userRepo.GetRefreshToken(ctx, guid)
	if err != nil {
		s.log.Errorln("не удалось получить refresh token", err)
		return token.Token{}, err
	}

	err = bcrypt.CompareHashAndPassword(hrefTok, []byte(reftoken))

	if err != nil {
		s.log.Errorln("Не удалось сравнить", err)
		return token.Token{}, err
	}

	accessTokenClaims := token.Claims{
		GUID: guid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(token.JwtSecret)
	if err != nil {
		s.log.Errorln("Не удалось создать jwt token", err)
		return token.Token{}, err
	}

	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		s.log.Errorln("Не удалось создать refresh token", err)
		return token.Token{}, err
	}

	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)

	hashedRefreshToken, err := s.HashRefreshToken(refreshToken)
	if err != nil {
		s.log.Errorln("Не удалось захэшировать refresh token", err)
		return token.Token{}, err
	}

	err = s.userRepo.UpdateRefreshToken(ctx, guid, hashedRefreshToken)

	if err != nil {
		s.log.Errorln("Не удалось обновить", err)
		return token.Token{}, err
	}

	return token.Token{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
	}, nil

}
