package users

import (
	"context"
	"crud-engine/modules/models"
	"crud-engine/pkg/middleware"
	"crud-engine/pkg/utils"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type HttpSqlx struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *HttpSqlx {
	return &HttpSqlx{db}
}

// Mount function
func (h *HttpSqlx) Mount(echoGroup *echo.Group) {
	echoGroup.POST("/login", h.Login, middleware.NoAuth())
}

// Login function
func (h *HttpSqlx) Login(c echo.Context) error {
	var (
		params models.ReqLogin
	)
	err := json.NewDecoder(c.Request().Body).Decode(&params)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}
	err = validator.New().Struct(params)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}
	user, err := h.GetByEmail(c.Request().Context(), params.Email)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusInternalServerError, c)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println(err)
		return utils.Response(nil, "Invalid Username or Password", http.StatusBadRequest, c)
	}

	token, err := CreateToken(user.Username)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}
	return utils.Response(token, "Login User", http.StatusOK, c)
}
func (s *HttpSqlx) GetByEmail(ctx context.Context, email string) (user *models.User, err error) {
	var u models.User
	query := `
	SELECT
		id,
		username,
		email,
		password,
		created_at,
		created_by,
		updated_at,
		last_update_by,
		is_deleted
	FROM
		users
	WHERE is_deleted = false AND email = $1
		`
	err = s.db.GetContext(ctx, &u, query, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
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
