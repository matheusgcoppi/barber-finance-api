package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"net/http"
	"os"
	"time"
)

type DatabaseMiddleware struct {
	database *database.CustomDB
}

func NewDatabaseMiddleware(db *database.CustomDB) *DatabaseMiddleware {
	return &DatabaseMiddleware{
		database: db,
	}
}

func (s *DatabaseMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization cookie not found")
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if !token.Valid {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
			}
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_JWT")), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Error parsing token: %v", err))
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			var user model.User
			s.database.Db.First(&user, claims["sub"])
			if user.ID == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			c.Set("user", user)

			fmt.Println("In middleware")
			return next(c)
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
	}
}
