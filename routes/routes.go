package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/service/user"
)

func SetupRoutes(e *echo.Echo, server *service.APIServer) {
	userRoutes(e, server)
}

func userRoutes(e *echo.Echo, server *service.APIServer) {
	e.GET("/", server.HandleIndex)
	e.GET("/user", server.HandleGetUser)
	e.GET("/user/:id", server.HandleGetUserByID)
	e.POST("/user", server.HandleCreateUser)
	e.DELETE("/user/:id", server.HandleDeleteUser)
	e.PUT("/user/:id", server.HandleUpdateUser)
}
