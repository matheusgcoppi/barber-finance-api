package tests

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/repository"
	service "github.com/matheusgcoppi/barber-finance-api/service/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetUser(t *testing.T) {
	// Create a new Echo instance and a request to test the route
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create an instance of your API server
	db, _ := database.NewPostgres()
	repositoryServer := repository.UserRepository{Store: db} // Create your repository server (if necessary)
	server := service.NewAPIServer(db, &repositoryServer)

	// Call the route handler function
	err := server.HandleGetUser(c)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert the HTTP status code (200 OK)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Add further assertions as needed to check the response body, headers, etc.
}

func TestHandleEnv(t *testing.T) {
	// Load the .env file
	if err := godotenv.Load("/Users/matheusgcoppi/Development/Golang/barber-finance/.env"); err != nil {
		t.Fatal("Error loading .env file")
	}

	// Your test logic here
}
