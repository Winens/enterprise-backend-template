package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/winens/enterprise-backend-template/pkg/model"
)

type HTTPCtxSessionTokenKeyType struct{}

var HTTPCtxSessionTokenKey HTTPCtxSessionTokenKeyType

func GetSessionFromCtx(c *fiber.Ctx) (*model.Session, error) {
	sess, ok := c.Locals(HTTPCtxSessionTokenKey).(*model.Session)
	if !ok {
		return nil, errors.New("session not found in context")
	}

	return sess, nil
}

func SetSessionToCtx(c *fiber.Ctx, sess *model.Session) {
	c.Locals(HTTPCtxSessionTokenKey, sess)
}
