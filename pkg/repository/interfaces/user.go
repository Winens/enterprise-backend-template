package interfaces

import (
	"context"

	"github.com/winens/enterprise-backend-template/pkg/model"
	"github.com/winens/enterprise-backend-template/pkg/repository/response"
	"github.com/winens/enterprise-backend-template/pkg/restapi/handler/request"
)

type UserRepository interface {
	BeginTx(ctx context.Context, callback func(txRepo UserRepository) error) error

	FindUserById(ctx context.Context, id int64) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)

	CreateUser(ctx context.Context, details request.SignUp) (userId int64, err error)

	ConfirmEmailByUserId(ctx context.Context, userId int64) error

	FetchEmailPasswordLoginData(ctx context.Context, email string) (*response.UserFetchEmailPasswordLogin, error)
}
