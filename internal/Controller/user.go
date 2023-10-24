package Controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleGetUser(c echo.Context) error {
	return c.JSON(http.StatusOK, "hey")
}
