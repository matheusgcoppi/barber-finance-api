package service

import (
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"github.com/matheusgcoppi/barber-finance-api/utils"
	"net/http"
	"strconv"
)

func (a *APIServer) HandleGetIncome(c echo.Context) error {
	income, err := a.repositoryServer.GetIncome()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, income)

}

func (a *APIServer) HandleGetIncomeById(c echo.Context) error {
	id := c.Param("id")
	income, err := a.repositoryServer.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, income)
}

func (a *APIServer) HandleCreateIncome(c echo.Context) error {
	incomeDTO := new(model.IncomeDTO)
	err := c.Bind(&incomeDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if incomeDTO.When.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Time is required."})
	}

	if incomeDTO.Price == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Price is required."})
	}

	if incomeDTO.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Description is required."})
	}

	if incomeDTO.Payment == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payment Type is required."})
	}

	currentId, err := utils.GetCurrentUserID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	id, errP := strconv.ParseUint(currentId, 10, 64)
	if errP != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format."})
	}
	parsedID := uint(id)

	if incomeDTO.AccountID == 0 {
		incomeDTO.AccountID = parsedID
	} else if incomeDTO.AccountID != parsedID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Mismatched user IDs."})
	}

	newIncome := a.repositoryServer.CreateNewIncome(
		incomeDTO.Price,
		incomeDTO.Description,
		incomeDTO.When,
		incomeDTO.Payment,
	)

	Income, err := a.repositoryServer.CreateIncome(newIncome)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, Income)

}

func (a *APIServer) HandleUpdateIncome(c echo.Context) error {
	id := c.Param("id")
	err := a.repositoryServer.DeleteIncome(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "User with id = " + id + " was deleted successfully"})
}

func (a *APIServer) HandleDeleteIncome(c echo.Context) error {
	id := c.Param("id")
	updatedIncome := new(model.IncomeDTO)
	err := c.Bind(&updatedIncome)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if updatedIncome.Payment != 0 {
		if (updatedIncome.Payment != model.Cash) &&
			(updatedIncome.Payment != model.CreditCard) &&
			(updatedIncome.Payment != model.DebitCard) &&
			(updatedIncome.Payment != model.Pix) &&
			(updatedIncome.Payment != model.Bill) &&
			(updatedIncome.Payment != model.Other) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "This Payment Type does not exist."})
		}
	}

	//TODO: think about update.Income.AccountID

	income, err := a.repositoryServer.UpdateIncome(updatedIncome, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, income)
}
