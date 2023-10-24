package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/internal/Controller"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", Controller.HandleGetUser)
}
