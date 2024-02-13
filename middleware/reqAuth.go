package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"github.com/matheusgcoppi/barber-finance-api/service"
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

func (s *DatabaseMiddleware) MiddlewareChain() echo.MiddlewareFunc {
	// Create a new CORS middleware instance
	corsMiddleware := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})
	return corsMiddleware
}

func (s *DatabaseMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization cookie not found")
		}

		var customClaims service.CustomClaims

		token, err := jwt.ParseWithClaims(cookie.Value, &customClaims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_JWT")), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Error parsing token: %v", err))
		}

		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Token Invalid"))
		}

		if claims, ok := token.Claims.(*service.CustomClaims); ok {
			if time.Now().After(claims.ExpiresAt.Time) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
			}

			var user model.User
			s.database.Db.First(&user, claims.Sub)
			if user.ID == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID in token")
			}

			c.Set("user", user)
			fmt.Println("In middleware")
			return next(c)
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
		}
	}
}
