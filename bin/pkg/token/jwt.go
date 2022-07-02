package token

import (
	"engine/bin/config"
	"engine/bin/pkg/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func Generate(uid string) (token *ResToken, err error) {
	// duration set 3600 seconds
	duration := (time.Hour * 1).Seconds()

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = uid
	claims["iat"] = time.Now().Unix()                                          //Token create
	claims["exp"] = time.Now().Add(config.GlobalEnv.AccessTokenExpired).Unix() //Token expired
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := tk.SignedString([]byte(config.GlobalEnv.APISecret))
	token = &ResToken{
		TokenType: "Bearer",
		Duration:  duration,
		Token:     accessToken,
	}
	return token, err

}

func Validate(r *http.Request) <-chan utils.Result {
	output := make(chan utils.Result)
	go func() {
		defer close(output)
		tokenString := ExtractToken(r)
		if len(tokenString) == 0 {
			output <- utils.Result{Error: "Missing authorization header"}
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.GlobalEnv.APISecret), nil
		})
		if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
			output <- utils.Result{Error: "Token has been expired"}
			return
		}
		if err != nil {
			output <- utils.Result{Error: "Invalid Token"}
			return
		}
		var tokenClaim Claim
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			tokenClaim = Claim{
				UserID: claims["userId"].(string),
			}
		}
		output <- utils.Result{Data: tokenClaim}
	}()
	return output
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
