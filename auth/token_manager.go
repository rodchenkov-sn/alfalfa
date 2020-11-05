package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/rodchenkov-sn/alfalfa/common"
	"github.com/rodchenkov-sn/alfalfa/service"
	"time"
)

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

type TokenManager struct {
	Repository *service.Repository
	PrivateKey []byte
}

func (tm *TokenManager) GenerateToken(info common.AuthInfo) (string, error) {
	exist, valid := tm.Repository.Authenticate(info)
	if !exist {
		return "", common.UserNotfoundError{Login: info.Login}
	}
	if !valid {
		return "", common.InvalidPasswordError{Login: info.Login}
	}
	claims := &Claims{
		Login: info.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.PrivateKey)
}

func (tm *TokenManager) ValidateToken(token string) (login string, err error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return tm.PrivateKey, nil
	})
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "", common.InvalidTokenError{}
	}
	if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return "", common.TokenExpiredError{}
	}
	return claims.Login, nil
}

func NewTokenManager(repository *service.Repository, privateKey string) *TokenManager {
	return &TokenManager{Repository: repository, PrivateKey: []byte(privateKey)}
}
