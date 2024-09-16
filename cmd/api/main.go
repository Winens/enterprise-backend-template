package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/winens/enterprise-backend-template/pkg/db"
	"github.com/winens/enterprise-backend-template/pkg/repository"
	"github.com/winens/enterprise-backend-template/pkg/restapi"
	"github.com/winens/enterprise-backend-template/pkg/restapi/handler"
	"github.com/winens/enterprise-backend-template/pkg/service"
	"github.com/winens/enterprise-backend-template/pkg/usecase"
	"go.uber.org/fx"
)

func main() {
	// Load configuration
	godotenv.Load()

	viper.AddConfigPath(".")
	viper.SetConfigFile("config.toml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Start the application
	fx.New(

		fx.Provide(db.NewPostgres, db.NewOneTimeTokenStore, service.NewTokenService,
			service.NewSMTPService),

		// Repositories
		fx.Provide(
			repository.NewAuthRepository,
			repository.NewUserRepository,
		),

		// UseCases
		fx.Provide(
			usecase.NewAuthUseCase,
		),

		// Handlers
		fx.Provide(
			handler.NewAuthHandler,
		),

		fx.Provide(restapi.NewMiddlewares, restapi.NewHTTPServer),

		fx.Invoke(restapi.SetupRoutes),
		fx.Invoke(restapi.RunHTTPServer),
	).Run()
}
