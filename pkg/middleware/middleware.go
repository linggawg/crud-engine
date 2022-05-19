package middleware

import (
	"crud-engine/modules/models"
	"crud-engine/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
			errMsg := TokenValid(c.Request())
			if errMsg != "" {
				return utils.Response(nil, "Unauthorized", http.StatusUnauthorized, c)
			}
			return next(c)
		}
	}
}

// CreateToken . . .
func CreateToken(uid string) (token *models.ResToken, err error) {
	// duration set 3600 seconds
	duration := (time.Hour * 1).Seconds()

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["sub"] = uid
	claims["iat"] = time.Now().Unix()                    //Token create
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := tk.SignedString([]byte(os.Getenv("API_SECRET")))
	token = &models.ResToken{
		TokenType: "Bearer",
		Duration:  duration,
		Token:     accessToken,
	}
	return token, err

}

// TokenValid . . .
func TokenValid(r *http.Request) (errsMsg string) {
	tokenString := ExtractToken(r)
	if len(tokenString) == 0 {
		return "Missing authorization header"
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
		return "Token has been expired"
	}
	if err != nil {
		log.Println(err)
		return "Invalid Token"
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return ""
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func Pretty(data interface{}) {
	_, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}
}
