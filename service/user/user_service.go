package service

import (
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"github.com/matheusgcoppi/barber-finance-api/repository"
	"net/http"
)

type APIServer struct {
	store            *database.CustomDB
	repositoryServer *repository.UserRepository
}

func NewAPIServer(store *database.CustomDB, repository *repository.UserRepository) *APIServer {
	return &APIServer{
		store:            store,
		repositoryServer: repository,
	}
}

func (a *APIServer) HandleIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, "hey")
}

func (a *APIServer) HandleCreateUser(c echo.Context) error {
	userDTO := new(model.UserDTO)
	if err := c.Bind(userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if userDTO.Type == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Type is required."})
	}

	if (userDTO.Type < 1) || (userDTO.Type > 3) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "This Role does not exist."})
	}

	if userDTO.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username is required."})
	}

	if userDTO.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required."})
	}

	if userDTO.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password is required."})
	}

	newUser := repository.NewUser(
		true,
		int(userDTO.Type),
		userDTO.Username,
		userDTO.Email,
		userDTO.Password,
	)

	if err := a.repositoryServer.CreateUser(newUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, newUser)
}

func (a *APIServer) HandleGetUser(c echo.Context) error {
	err, result := a.repositoryServer.GetUser()
	if err != nil {
		err := c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		if err != nil {
			return err
		}
	}
	return c.JSON(http.StatusOK, result)
}
