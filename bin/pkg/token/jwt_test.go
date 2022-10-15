package token

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		_, err := Generate("uud1", "Admin", 10000000000)
		assert.Error(t, err)
	})
}
func TestValidate(t *testing.T) {
	var accessToken string
	token, err := Generate("uud1", "Admin", 3600)
	assert.NoError(t, err)
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(10).Unix()
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/engine/v1/users", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.Token))
		res := <-Validate(req)
		assert.Nil(t, res.Error)
	})
	t.Run("missing-authorization-header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/engine/v1/users", strings.NewReader(""))
		res := <-Validate(req)
		assert.Equal(t, res.Error, "Missing authorization header")
	})
	t.Run("unexpected-signing-method", func(t *testing.T) {
		tk := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		accessToken, err = tk.SignedString(jwt.UnsafeAllowNoneSignatureType)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/engine/v1/users", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, accessToken))
		res := <-Validate(req)
		assert.Equal(t, res.Error, "Invalid Token")
	})
	t.Run("error-signed", func(t *testing.T) {
		tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		accessToken, err = tk.SignedString([]byte("api_secret"))
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/engine/v1/users", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, accessToken))
		res := <-Validate(req)
		assert.Equal(t, res.Error, "Invalid Token")
	})
	t.Run("token-expired", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/engine/v1/users", strings.NewReader(""))
		expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTk2MTQ0NjgsImlhdCI6MTY1OTYxNDQ2Nywic2NvcGUiOiJBZG1pbiIsInVzZXJJZCI6InV1ZDEifQ.ImSdaAXfnph4sqmtUbnlMml8iJJq7xtxYPhmy2oZGRo"
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, expiredToken))
		res := <-Validate(req)
		assert.Equal(t, res.Error, "Token has been expired")
	})
}
