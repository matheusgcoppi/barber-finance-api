package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/matheusgcoppi/barber-finance-api/internal/routes"
	"github.com/matheusgcoppi/barber-finance-api/service"
)

func main() {
	e := echo.New()
	err := service.DBconnection()
	if err != nil {
		log.Fatal(err)
	}

	routes.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
