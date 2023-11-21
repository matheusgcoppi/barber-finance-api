package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/middleware"
	"github.com/matheusgcoppi/barber-finance-api/service/user"
)

func SetupRoutes(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {

	userRoutes(e, server, middleware)
}

func userRoutes(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {
	e.GET("/", server.HandleIndex)
	e.GET("/user", server.HandleGetUser)
	e.GET("/user/:id", server.HandleGetUserByID)
	e.POST("/user", server.HandleCreateUser)
	e.POST("/login", server.HandleLogin)
	e.DELETE("/user/:id", server.HandleDeleteUser)
	e.PUT("/user/:id", server.HandleUpdateUser)
	e.GET("validate", middleware.RequireAuth(server.Validate))
}
