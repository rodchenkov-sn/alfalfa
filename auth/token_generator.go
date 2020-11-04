package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/rodchenkov-sn/alfalfa/service"
	"time"
)

type TokenGenerator struct {
	Repository *service.Repository
}

func (tg *TokenGenerator) GenerateToken(info service.AuthInfo) (string, error) {
	if err := tg.Repository.AddUser(info); err != nil {
		return "", err
	}
	secretKey := "my_secret_key"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": info.Login,
		"expired": time.Now().Add(48 * time.Hour).Unix(),
	})
	return token.SignedString(secretKey)
}

func NewTokenGenerator(repository *service.Repository) *TokenGenerator {
	return &TokenGenerator{Repository: repository}
}
