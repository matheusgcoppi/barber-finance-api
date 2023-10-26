package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/repository"
	"github.com/matheusgcoppi/barber-finance-api/service/user"
)

func SetupRoutes(e *echo.Echo, server *service.APIServer, userRepository *repository.UserRepository) {
	e.GET("/", server.HandleIndex)
	e.GET("/user", server.HandleGetUser)
	e.POST("/user", server.HandleCreateUser)
}
