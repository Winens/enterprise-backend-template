package interfaces

import "github.com/golang-jwt/jwt/v5"

type TokenService interface {
	Sign(claims jwt.Claims, opts ...jwt.TokenOption) (string, error)
	Verify(tokenString string, claims jwt.Claims, opts ...jwt.ParserOption) (*jwt.Token, error)
}
