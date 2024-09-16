package handler

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/winens/enterprise-backend-template/pkg/errdefs"
	handlers "github.com/winens/enterprise-backend-template/pkg/restapi/handler/interfaces"
	"github.com/winens/enterprise-backend-template/pkg/restapi/handler/request"
	usecases "github.com/winens/enterprise-backend-template/pkg/usecase/interfaces"
	"github.com/winens/enterprise-backend-template/pkg/utils"
)

type authHandler struct {
	authUseCase usecases.AuthUseCase
}

func NewAuthHandler(authUseCase usecases.AuthUseCase) handlers.AuthHandler {
	return &authHandler{authUseCase}
}

func (h *authHandler) SignUp(c *fiber.Ctx) error {
	var req request.SignUp

	if err := c.BodyParser(&req); err != nil {
		return errdefs.BadRequest
	}

	// TODO: add validator

	ctx := context.Background()
	if err := h.authUseCase.SignUp(ctx, req); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var req request.Login

	if err := c.BodyParser(&req); err != nil {
		return errdefs.BadRequest
	}

	ip := c.IP()
	userAgent := c.Get(fiber.HeaderUserAgent)

	sessionToken, err := h.authUseCase.Login(context.Background(), req, ip, userAgent)

	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return c.SendStatus(fiber.StatusOK)
}

func (h *authHandler) ConfirmEmail(c *fiber.Ctx) error {
	err := h.authUseCase.ConfirmEmail(context.Background(), c.Query("token"))
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *authHandler) Logout(c *fiber.Ctx) error {
	sess, err := utils.GetSessionFromCtx(c)
	if err != nil {
		return err
	}

	if err := h.authUseCase.Logout(context.Background(), sess); err != nil {
		return fmt.Errorf("failed to logout", err)
	}

	c.ClearCookie("session_token")
	return c.SendStatus(fiber.StatusOK)
}

func (h *authHandler) GetLoggedInUser(c *fiber.Ctx) error {
	sess, err := utils.GetSessionFromCtx(c)
	if err != nil {
		return err
	}

	user, err := h.authUseCase.GetLoggedInUser(context.Background(), sess)
	if err != nil {
		return fmt.Errorf("failed to get loggedin user: %w", err)
	}

	return c.JSON(user)
}
