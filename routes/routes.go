package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/middleware"
	"github.com/matheusgcoppi/barber-finance-api/service"
)

func SetupRoutes(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {
	user(e, server, middleware)
	income(e, server, middleware)
	expense(e, server, middleware)
}

func user(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {
	e.GET("/", server.IndexHandler)
	e.GET("/user", middleware.RequireAuth(server.HandleGetUser))
	e.GET("/user/:id", middleware.RequireAuth(server.HandleGetUserByID))
	e.POST("/user", server.HandleCreateUser)
	e.POST("/login", server.HandleLogin)
	e.DELETE("/user/:id", middleware.RequireAuth(server.HandleDeleteUser))
	e.PUT("/user/:id", middleware.RequireAuth(server.HandleUpdateUser))
	e.GET("/validate", middleware.RequireAuth(server.Validate))
}

func income(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {
	e.GET("/incomes", middleware.RequireAuth(server.HandleGetIncome))
	e.GET("/income/:id", middleware.RequireAuth(server.HandleGetIncomeById))
	e.POST("/income", middleware.RequireAuth(server.HandleCreateIncome))
	e.DELETE("/income/:id", middleware.RequireAuth(server.HandleDeleteIncome))
	e.PUT("/income/:id", middleware.RequireAuth(server.HandleUpdateIncome))
}

func expense(e *echo.Echo, server *service.APIServer, middleware *middleware.DatabaseMiddleware) {
	e.GET("/expenses", middleware.RequireAuth(server.HandleGetExpenses))
	e.GET("/expense/:id", middleware.RequireAuth(server.HandleGetExpenseById))
	e.POST("/expense", middleware.RequireAuth(server.HandleCreateExpense))
	e.DELETE("/expense/:id", middleware.RequireAuth(server.HandleDeleteExpense))
	e.PUT("/expense/:id", middleware.RequireAuth(server.HandleUpdateExpense))
}
