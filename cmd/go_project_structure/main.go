package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/middleware"
	"github.com/matheusgcoppi/barber-finance-api/repository"
	"github.com/matheusgcoppi/barber-finance-api/routes"
	"github.com/matheusgcoppi/barber-finance-api/service/user"
)

func main() {
	e := echo.New()
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	userRepository := repository.UserRepository{Store: db}

	server := service.NewAPIServer(db, &userRepository)

	middleware := middleware.NewDatabaseMiddleware(db)

	routes.SetupRoutes(e, server, middleware)

	e.Logger.Fatal(e.Start(":8080"))

}
