package restapi

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	usecases "github.com/winens/enterprise-backend-template/pkg/usecase/interfaces"
	"github.com/winens/enterprise-backend-template/pkg/utils"
)

type middlewares struct {
	authUseCase usecases.AuthUseCase
}

func NewMiddlewares(authUseCase usecases.AuthUseCase) *middlewares {
	return &middlewares{authUseCase}
}

func (m *middlewares) FetchSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := m.authUseCase.FindSessionByToken(context.Background(), c.Cookies("session_token"))
		if err != nil {
			return fmt.Errorf("finding session failed: %w", err)
		}

		utils.SetSessionToCtx(c, sess)
		return c.Next()
	}
}
