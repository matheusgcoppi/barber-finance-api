package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"os"
	"strconv"
	"strings"
)

const (
	EmailPattern = `^([a-z\d\.-]+)@([a-z\d-]+)\.([a-z]{2,10})(\.[a-z]{2,8})?$`
)

func ContainsUtil(words []string, json string) bool {
	length := len(words)
	count := 0
	for i := 0; i < length; i++ {
		if strings.Contains(json, words[i]) {
			count++
		}
	}
	if length == count {
		return true
	}
	return false
}

func GetCurrentUserID(c echo.Context) (string, error) {
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	var userId string

	if id, ok := claims["sub"].(float64); ok {
		userId = strconv.Itoa(int(id))
	}
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	return userId, nil
}
