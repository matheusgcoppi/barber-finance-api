package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"github.com/matheusgcoppi/barber-finance-api/repository"
	"github.com/matheusgcoppi/barber-finance-api/utils"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"regexp"
	"time"
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (a *APIServer) HandleLogin(c echo.Context) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := a.repositoryServer.LoginUser(body.Email, body.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_JWT")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid to create token",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30)
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
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

	if (userDTO.Type != model.System) && (userDTO.Type != model.Support) && (userDTO.Type != model.Admin) {
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

	password, err := HashPassword(userDTO.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	newUser := repository.NewUser(
		true,
		int(userDTO.Type),
		userDTO.Username,
		userDTO.Email,
		password,
	)

	err, createdUser := a.repositoryServer.CreateUser(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, createdUser)
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

func (a *APIServer) HandleGetUserByID(c echo.Context) error {
	id := c.Param("id")
	result, err := a.repositoryServer.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *APIServer) HandleDeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := a.repositoryServer.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "User with id = " + id + " was deleted successfully"})
}

func (a *APIServer) HandleUpdateUser(c echo.Context) error {
	id := c.Param("id")
	updatedUser := new(model.UserDTO)
	if err := c.Bind(updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if updatedUser.Type != 0 {
		if (updatedUser.Type != model.System) && (updatedUser.Type != model.Support) && (updatedUser.Type != model.Admin) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "This Role does not exist."})
		}
	}

	if updatedUser.Email != "" {
		match, _ := regexp.MatchString(utils.EmailPattern, updatedUser.Email)
		if match == false {
			println(match)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Email Address"})
		}
	}

	if updatedUser.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password cannot be null"})
	}

	user, err := a.repositoryServer.UpdateUser(updatedUser, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (a *APIServer) Validate(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "I'm logged in"})
}
