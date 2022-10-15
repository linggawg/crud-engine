package middleware_test

import (
	"engine/bin/middleware"
	"engine/bin/pkg/token"
	"engine/bin/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNoAuth(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return utils.Response(nil, "success", http.StatusOK, c)
	}, middleware.NoAuth())
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "application/json; charset=UTF-8", res.Header().Get(echo.HeaderContentType))
}
func TestVerifyBearer(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return utils.Response(nil, "success", http.StatusOK, c)
	}, middleware.VerifyBearer())
	t.Run("success", func(t *testing.T) {
		generate, err := token.Generate("uud1", "Admin", 3600)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", generate.TokenType, generate.Token))
		e.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "application/json; charset=UTF-8", res.Header().Get(echo.HeaderContentType))
	})
	t.Run("missing-authorization-header", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		e.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Equal(t, "application/json; charset=UTF-8", res.Header().Get(echo.HeaderContentType))
	})
}
