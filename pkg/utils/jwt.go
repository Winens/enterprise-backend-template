package utils

import (
	"crypto/ed25519"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JWTKeyEd25519 struct {
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}

// Sign signs the claims and returns JWT token.
func (k JWTKeyEd25519) Sign(claims jwt.Claims, opts ...jwt.TokenOption) (string, error) {
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims, opts...)
	return token.SignedString(k.Private)
}

// Verify verifies the JWT token and returns the token.
func (k JWTKeyEd25519) Verify(tokenString string, claims jwt.Claims, opts ...jwt.ParserOption) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return k.Public, nil
	}, opts...)
}

func LoadJWTKeyEd25519(privateKeyPath, publicKeyPath string) (*JWTKeyEd25519, error) {
	// Load PEM encoded private key and public key.
	privatePem, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	publicPem, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	// Decode the PEM encoded data.
	privateDecoded, err := jwt.ParseEdPrivateKeyFromPEM(privatePem)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	publicDecoded, err := jwt.ParseEdPublicKeyFromPEM(publicPem)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	// Type conversion.
	privateKey, ok := privateDecoded.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("failed to convert private key to ed25519.PrivateKey")
	}

	publicKey, ok := publicDecoded.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to convert public key to ed25519.PublicKey")
	}

	return &JWTKeyEd25519{
		Private: privateKey,
		Public:  publicKey,
	}, nil
}
