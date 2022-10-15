package token

import (
	"engine/bin/config"
	"engine/bin/pkg/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func Generate(uid, roleName string, expire int) (token *ResToken, err error) {
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", expire))
	if err != nil {
		return nil, errors.New("failed parsing duration time")
	}

	claims := jwt.MapClaims{}
	claims["scope"] = roleName
	claims["authorized"] = true
	claims["userId"] = uid
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := tk.SignedString([]byte(config.GlobalEnv.APISecret))
	token = &ResToken{
		TokenType: "Bearer",
		Duration:  duration.Seconds(),
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
				UserID:   claims["userId"].(string),
				RoleName: claims["scope"].(string),
			}
		}
		output <- utils.Result{Data: tokenClaim}
	}()
	return output
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
