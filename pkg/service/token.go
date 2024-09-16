package service

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/winens/enterprise-backend-template/pkg/service/interfaces"
)

type service struct {
	// todo: we'll add KMS(key managament service) here in future.

	secretKey []byte
}

func NewTokenService() interfaces.TokenService {
	secretKey := viper.GetString("API_SECRET_KEY")
	if secretKey == "" {
		panic("API_SECRET_KEY is not set")
	}

	return &service{
		secretKey: []byte(secretKey),
	}
}

func (s *service) Sign(claims jwt.Claims, opts ...jwt.TokenOption) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims, opts...)
	return token.SignedString(s.secretKey)
}

func (s *service) Verify(tokenString string, claims jwt.Claims, opts ...jwt.ParserOption) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return s.secretKey, nil
	}, opts...)
}
