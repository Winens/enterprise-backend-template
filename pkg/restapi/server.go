package restapi

import (
	"context"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type HTTPServer struct {
	App *fiber.App
}

func NewHTTPServer(
// Here we pass all handlers as parameters so we can use them in routes.go
) *HTTPServer {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
	})

	// TODO: logger middleware w/ zerolog multiwriter (1: console, 2: es)
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	return &HTTPServer{
		App: app,
	}
}

func RunHTTPServer(lc fx.Lifecycle, s *HTTPServer) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := s.App.Listen(":8080"); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.App.ShutdownWithContext(ctx)
		},
	})
}
