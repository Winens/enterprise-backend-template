package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	crand "crypto/rand"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/winens/enterprise-backend-template/pkg/errdefs"
	"github.com/winens/enterprise-backend-template/pkg/model"
	repos "github.com/winens/enterprise-backend-template/pkg/repository/interfaces"
	"github.com/winens/enterprise-backend-template/pkg/restapi/handler/request"
	services "github.com/winens/enterprise-backend-template/pkg/service/interfaces"
	"github.com/winens/enterprise-backend-template/pkg/usecase/interfaces"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	passwordHashCost = 10
)

type authUseCase struct {
	authRepo     repos.AuthRepository
	userRepo     repos.UserRepository
	tokenService services.TokenService
	smtpService  services.SMTPService
}

func NewAuthUseCase(
	authRepo repos.AuthRepository, userRepo repos.UserRepository, tokenService services.TokenService,
	smtpService services.SMTPService,
) interfaces.AuthUseCase {
	return &authUseCase{authRepo, userRepo, tokenService, smtpService}
}

func (u *authUseCase) SignUp(ctx context.Context, details request.SignUp) error {
	// TODO: repo.CheckEmai	lExists(...)

	// check if email already exists
	user, err := u.userRepo.FindUserByEmail(ctx, details.Email)
	if err != nil && !errors.Is(err, errdefs.UserNotFound) {
		return err
	}

	if user != nil {
		return errdefs.EmailAlreadyExists
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(details.Password), passwordHashCost)
	if err != nil {
		return fmt.Errorf("failed to generate bcrypt hash", err)
	}

	// pass the hash into object
	details.Password = string(passwordHash)

	return u.userRepo.BeginTx(ctx, func(txRepo repos.UserRepository) error {
		userId, err := txRepo.CreateUser(ctx, details)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// generate verification token
		claims := jwt.RegisteredClaims{
			Issuer:    "enterprise-backend-template",
			Subject:   strconv.FormatInt(userId, 10),
			Audience:  []string{"email_confirmation"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		}

		verificationToken, err := u.tokenService.Sign(claims)
		if err != nil {
			return fmt.Errorf("failed to generate verification token", err)
		}

		// send verification email
		return u.smtpService.SendUserVerificationEmail(details.Email, details.FirstName, verificationToken)

	})

}

func (u *authUseCase) Login(ctx context.Context, details request.Login, ip, userAgent string) (sessionToken string, err error) {
	userLoginData, err := u.userRepo.FetchEmailPasswordLoginData(ctx, details.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user login info: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(userLoginData.PasswordHash), []byte(details.Password)) != nil {
		return "", errdefs.InvalidCredentials
	}

	// generate session token

	sessionId, err := u.authRepo.NewSession(ctx, userLoginData.UserId, ip, userAgent)
	if err != nil {
		return "", err
	}

	key, err := base64.StdEncoding.DecodeString(os.Getenv("SESSION_SECRET"))
	if err != nil {
		return "", fmt.Errorf("failed to decode SESSION_SECRET: %w", err)
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AEAD: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	n, err := crand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to create nonce: %w", err)
	}

	if n != aead.NonceSize() {
		return "", fmt.Errorf("failed to create nonce: not enough random bytes")
	}

	// session_id'yi şifreleyip string olarak döndürüyoruz.
	tkn := aead.Seal(nonce, nonce, sessionId[:], nil)

	return base64.URLEncoding.EncodeToString(tkn), nil

}

func (u *authUseCase) ConfirmEmail(ctx context.Context, token string) error {
	var claims jwt.RegisteredClaims

	_, err := u.tokenService.Verify(token, &claims, jwt.WithAudience("email_confirmation"),
		jwt.WithIssuer("enterprise-backend-template"))

	if err != nil {
		return fmt.Errorf("failed to verify email verification token: %w", errdefs.InvalidToken, err)
	}

	userId, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return errdefs.InvalidToken
	}

	return u.userRepo.ConfirmEmailByUserId(ctx, userId)
}

func (u *authUseCase) FindSessionByToken(ctx context.Context, sessionToken string) (sess *model.Session, err error) {
	// todo: dexcrypt session token into -> id

	key, err := base64.StdEncoding.DecodeString(os.Getenv("SESSION_SECRET"))
	if err != nil {
		return nil, fmt.Errorf("failed key", err)
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD: %w", err)
	}

	decodedToken, err := base64.URLEncoding.DecodeString(sessionToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(decodedToken) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := decodedToken[:nonceSize], decodedToken[nonceSize:]

	sessionIdBytes, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	sessionId, err := uuid.FromBytes(sessionIdBytes[:])

	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid bytes: %w", err)
	}

	return u.authRepo.FindSessionById(ctx, sessionId)
}

func (u *authUseCase) GetLoggedInUser(ctx context.Context, sess *model.Session) (*model.User, error) {
	return u.userRepo.FindUserById(ctx, sess.UserId)
}

func (u *authUseCase) Logout(ctx context.Context, sess *model.Session) error {
	return u.authRepo.DeleteSession(ctx, sess.Id)
}
