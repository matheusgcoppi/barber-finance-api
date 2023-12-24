package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/middleware"
	"github.com/matheusgcoppi/barber-finance-api/service"
)

func SetupRoutes(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {
	userRoutes(e, server, middleware)
	incomeRoutes(e, server, middleware)
}

func userRoutes(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {
	e.GET("/", server.IndexHandler)
	e.GET("/user", middleware.RequireAuth(server.HandleGetUser))
	e.GET("/user/:id", middleware.RequireAuth(server.HandleGetUserByID))
	e.POST("/user", server.HandleCreateUser)
	e.POST("/login", server.HandleLogin)
	e.DELETE("/user/:id", middleware.RequireAuth(server.HandleDeleteUser))
	e.PUT("/user/:id", middleware.RequireAuth(server.HandleUpdateUser))
	e.GET("validate", middleware.RequireAuth(server.Validate))
}

func incomeRoutes(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {
	e.GET("/incomes", middleware.RequireAuth(server.HandleGetIncome))
	e.GET("/income/:id", middleware.RequireAuth(server.HandleGetIncomeById))
	e.POST("/income", middleware.RequireAuth(server.HandleCreateIncome))
	e.DELETE("/income/:id", middleware.RequireAuth(server.HandleDeleteIncome))
	e.PUT("/income/:id", middleware.RequireAuth(server.HandleUpdateIncome))
}
