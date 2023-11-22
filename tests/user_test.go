package tests

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/repository"
	service "github.com/matheusgcoppi/barber-finance-api/service/user"
	"github.com/matheusgcoppi/barber-finance-api/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *service.APIServer
	db     *database.CustomDB
)

func init() {
	var err error
	db, err = database.NewPostgresTest()
	if err != nil {
		panic("Failed to initialize the test database: " + err.Error())
	}

	repositoryServer := repository.UserRepository{Store: db}
	server = service.NewAPIServer(db, &repositoryServer)
}

func createServer(t *testing.T, method string, path string, body io.Reader, id string) (*service.APIServer, echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	method = strings.ToUpper(method)
	var req *http.Request
	if method == "GET" {
		req = httptest.NewRequest(http.MethodGet, path, body)
	} else if method == "POST" {
		req = httptest.NewRequest(http.MethodPost, path, body)
	} else if method == "PUT" {
		req = httptest.NewRequest(http.MethodPut, path, body)
	} else if method == "DELETE" {
		req = httptest.NewRequest(http.MethodDelete, path, body)
	} else {
		t.Error("Invalid method:", method)
	}

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if id != "" {
		c.SetParamNames("id")
		c.SetParamValues(id)
	}

	c.SetPath(path)

	return server, c, rec
}

func TestCreateUser(t *testing.T) {
	body := strings.NewReader(`{"type": 2, "username": "matheus","email": "johndoe@example.com", "password": "teste"}`)
	server, c, rec := createServer(t, "post", "/user", body, "")

	err := server.HandleCreateUser(c)

	jsonResponse := rec.Body.String()
	var words = []string{"matheus", "johndoe@example.com", "teste"}

	if !utils.ContainsUtil(words, jsonResponse) {
		t.Errorf("The JSON response does not contain expected words: %v", words)
	}

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert the HTTP status code (you can customize this based on your implementation)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUser(t *testing.T) {
	server, c, rec := createServer(t, "get", "/user", nil, "")

	// Call the route handler function
	err := server.HandleGetUser(c)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert the HTTP status code (200 OK)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserByID(t *testing.T) {
	server, c, rec := createServer(t, "get", "/user/:id", nil, "1")

	err := server.HandleGetUserByID(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateUser(t *testing.T) {
	body := strings.NewReader(`{"type": 3, "username": "update","email": "update@example.com", "password": "update"}`)
	server, c, rec := createServer(t, "put", "/user/:id", body, "1")

	err := server.HandleUpdateUser(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteUser(t *testing.T) {
	server, c, rec := createServer(t, "delete", "/user/:id", nil, "1")

	err := server.HandleDeleteUser(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleEnv(t *testing.T) {
	if err := godotenv.Load("/Users/matheusgcoppi/Development/Golang/barber-finance/.env"); err != nil {
		t.Fatal("Error loading .env file")
	}
}
