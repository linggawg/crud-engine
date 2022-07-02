package middleware

import (
	"engine/bin/pkg/token"
	"engine/bin/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NoAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(echo.HeaderContentType, "application/json; charset=utf-8")
			return next(c)
		}
	}
}

func VerifyBearer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("Content-Type", "application/json; charset=utf-8")

			parsedToken := <-token.Validate(c.Request())
			if parsedToken.Error != nil {
				return utils.Response(nil, "Unauthorized", http.StatusUnauthorized, c)
			}
			tokenData := parsedToken.Data.(token.Claim)
			tokenData.Authorization = c.Request().Header.Get(echo.HeaderAuthorization)
			c.Set("opts", tokenData)
			return next(c)
		}
	}
}
