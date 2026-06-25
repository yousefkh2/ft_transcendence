package middleware

import (
	"net/http"
	"strings"

	"transcendence/backend/internal/auth"

	"github.com/labstack/echo/v4"
)

func EchoWithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		userID, err := auth.ParseJWT(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		c.Set("userID", userID)
		return next(c)
	}
}

func UserID(c echo.Context) string {
	return c.Get("userID").(string)
}
