package service

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"github.com/matheusgcoppi/barber-finance-api/utils"
	"net/http"
	"reflect"
	"strconv"
)

func (a *APIServer) HandleCreateExpense(c echo.Context) error {
	expenseDTO := new(model.ExpenseDTO)

	err := c.Bind(&expenseDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if expenseDTO.When.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Time is required."})
	}

	if expenseDTO.Price == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Price is required."})
	}

	if expenseDTO.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Description is required."})
	}

	if expenseDTO.Payment == 0 {
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

	var result int

	a.store.Db.Raw("SELECT id FROM accounts WHERE user_id = ?;", parsedID).Scan(&result)

	if expenseDTO.AccountID == 0 {
		expenseDTO.AccountID = uint(result)
	} else if expenseDTO.AccountID != uint(result) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Mismatched user IDs."})
	}

	fmt.Println(expenseDTO.AccountID, reflect.TypeOf(expenseDTO.AccountID))

	newExpense := a.repositoryServer.CreateNewExpense(
		expenseDTO.AccountID,
		expenseDTO.Price,
		expenseDTO.Description,
		expenseDTO.When,
		expenseDTO.Payment,
	)

	expense, err := a.repositoryServer.CreateExpense(newExpense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	fmt.Println(expense)

	return c.JSON(http.StatusOK, expense)
}

func (a *APIServer) HandleGetExpenses(c echo.Context) error {
	expense, err := a.repositoryServer.GetExpense()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, expense)
}

func (a *APIServer) HandleGetExpenseById(c echo.Context) error {
	id := c.Param("id")
	expense, err := a.repositoryServer.GetExpenseById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, expense)
}

func (a *APIServer) HandleDeleteExpense(c echo.Context) error {
	id := c.Param("id")
	err := a.repositoryServer.DeleteExpense(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "Income id = " + id + " was deleted successfully"})
}

func (a *APIServer) HandleUpdateExpense(c echo.Context) error {
	id := c.Param("id")
	updatedExpense := new(model.ExpenseDTO)
	err := c.Bind(&updatedExpense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if updatedExpense.Payment != 0 {
		if (updatedExpense.Payment != model.Cash) &&
			(updatedExpense.Payment != model.CreditCard) &&
			(updatedExpense.Payment != model.DebitCard) &&
			(updatedExpense.Payment != model.Pix) &&
			(updatedExpense.Payment != model.Bill) &&
			(updatedExpense.Payment != model.Other) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "This Payment Type does not exist."})
		}
	}

	expense, err := a.repositoryServer.UpdateExpense(updatedExpense, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, expense)
}
