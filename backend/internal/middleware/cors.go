package middleware

import(
	"net/http"

	echomw "github.com/labstack/echo/v4/middleware"
)

var WithCors = echomw.CORSWithConfig(echomw.CORSConfig{
	AllowOrigins: []string{"http://localhost:5173"},
	AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	AllowHeaders: []string{"Content-Type"},
})
