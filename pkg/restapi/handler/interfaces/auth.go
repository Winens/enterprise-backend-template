package interfaces

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	SignUp(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	ConfirmEmail(c *fiber.Ctx) error

	GetLoggedInUser(c *fiber.Ctx) error
}
