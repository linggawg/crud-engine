package users

import (
	"context"
	"crud-engine/modules/models"
	"crud-engine/pkg/middleware"
	"crud-engine/pkg/utils"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
	"log"
	"net/http"
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
	echoGroup.POST("/register", h.RegisterUser, middleware.NoAuth())
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

	token, err := middleware.CreateToken(user.Username)
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

func (s *HttpSqlx) RegisterUser(c echo.Context) error {
	var (
		params models.ReqUser
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

	_, err = s.GetByEmail(c.Request().Context(), params.Email)
	if err == nil {
		return utils.Response(nil, "Email has been used", http.StatusFound, c)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusInternalServerError, c)
	}
	password := string(hashedPassword)

	user := &models.User{
		Username:  params.Username,
		Email:     params.Email,
		Password:  password,
		CreatedBy: params.UserId,
		CreatedAt: null.TimeFrom(time.Now()),
		IsDeleted: false,
	}
	err = s.Insert(c.Request().Context(), user)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusInternalServerError, c)
	}

	return utils.Response(nil, "Register user", http.StatusOK, c)

}

func (s *HttpSqlx) Insert(ctx context.Context, user *models.User) (err error) {
	query := `
	INSERT INTO users
		(
			username,
			email,
			password,
			created_at,
			created_by,
			last_update_by,
			updated_at,
			is_deleted
		) 
		VALUES 
		(
			:username,
			:email,
			:password,
			:created_at,
			:created_by,
			:last_update_by,
			:updated_at,
			:is_deleted
		) RETURNING id ;
	`
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := s.db.NamedQueryContext(ctx, query, &models.User{
		Username:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		CreatedAt:    user.CreatedAt,
		CreatedBy:    user.CreatedBy,
		UpdatedAt:    user.CreatedAt,
		LastUpdateBy: &user.CreatedBy,
		IsDeleted:    user.IsDeleted,
	})
	if err != nil {
		return err
	}
	if res.Next() {
		res.Scan(&user.ID)
	}
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}
