package errdefs

import "github.com/gofiber/fiber/v2"

var (
	NoRowsAffected      = fiber.NewError(fiber.StatusNotFound, "no rows affected")
	BadRequest          = fiber.NewError(fiber.StatusBadRequest, "bad request")
	NotFound            = fiber.NewError(fiber.StatusNotFound, "not found")
	InternalServerError = fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	Forbidden           = fiber.NewError(fiber.StatusForbidden, "forbidden")
	Unauthorized        = fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	InvalidToken        = fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	ExpiredToken        = fiber.NewError(fiber.StatusUnauthorized, "expired token")

	InvalidRequest = fiber.NewError(fiber.StatusBadRequest, "invalid request")
	TokenExpired   = fiber.NewError(fiber.StatusUnauthorized, "token expired")
)

// User route errors
var (
	UserNotFound          = fiber.NewError(fiber.StatusNotFound, "user not found")
	EmailAlreadyExists    = fiber.NewError(fiber.StatusConflict, "email already exists")
	EmailNotConfirmed     = fiber.NewError(fiber.StatusForbidden, "email not confirmed")
	EmailAlreadyConfirmed = fiber.NewError(fiber.StatusConflict, "email already confirmed")
	InvalidCredentials    = fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
)
