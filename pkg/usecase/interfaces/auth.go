package interfaces

import (
	"context"

	"github.com/winens/enterprise-backend-template/pkg/model"
	"github.com/winens/enterprise-backend-template/pkg/restapi/handler/request"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, signupDetails request.SignUp) error
	Login(ctx context.Context, loginDetails request.Login, ip, userAgent string) (sessionToken string, err error)
	Logout(ctx context.Context, sess *model.Session) error
	ConfirmEmail(ctx context.Context, token string) error

	FindSessionByToken(ctx context.Context, sessionToken string) (sess *model.Session, err error)
	GetLoggedInUser(ctx context.Context, sess *model.Session) (*model.User, error)
}
