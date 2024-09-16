package restapi

import (
	handlers "github.com/winens/enterprise-backend-template/pkg/restapi/handler/interfaces"
)

func SetupRoutes(srv *HTTPServer, middlewares *middlewares, authHandler handlers.AuthHandler) {
	r := srv.App.Group("/")
	{
		rAuth := r.Group("/auth")
		{
			rAuth.Post("/sign-up", authHandler.SignUp)
			// TODOs: rate limiting and max session number must be added.
			rAuth.Post("/login", authHandler.Login)
			rAuth.Post("/logout", middlewares.FetchSession(), authHandler.Logout)
			rAuth.Get("/confirm-email", authHandler.ConfirmEmail)

			// TODOS:
			// [] Add forgot password
		}

		r.Get("/users/current", middlewares.FetchSession(), authHandler.GetLoggedInUser)

		// TODOS:
		// [] Add user update
		// [] Password update
		// [] Avatar update
	}

}
