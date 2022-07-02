package handlers

import (
	"encoding/json"
	"engine/bin/middleware"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/modules/users/repositories/commands"
	"engine/bin/modules/users/repositories/queries"
	"engine/bin/modules/users/usecases"
	conn "engine/bin/pkg/databases"
	"engine/bin/pkg/utils"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type HttpSqlx struct {
	db             *sqlx.DB
	commandUsecase usecases.CommandUsecase
}

func New() *HttpSqlx {
	db := conn.InitSqlx()
	postgreQuery := queries.NewUsersQuery(db)
	postgreCommand := commands.NewUsersCommand(db)
	commandUsecase := usecases.NewCommandUsecase(postgreCommand, postgreQuery)

	return &HttpSqlx{
		db:             db,
		commandUsecase: commandUsecase,
	}
}

// Mount function
func (h *HttpSqlx) Mount(echoGroup *echo.Group) {
	echoGroup.POST("/login", h.Login, middleware.NoAuth())
	echoGroup.POST("/register", h.RegisterUser, middleware.NoAuth())
}

// Post Login godoc
// @Summary      Login
// @Description  Login api
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param		 insertRequest body models.ReqLogin true "JSON request body based on column name"
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /login [post]
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

	result := h.commandUsecase.Login(c.Request().Context(), params)
	if result.Error != nil {
		return utils.Response(nil, "Login failed", http.StatusBadRequest, c)
	}

	return utils.Response(result.Data, "Login User", http.StatusOK, c)
}

// Post RegisterUser godoc
// @Summary      Register
// @Description  Register new user for login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param		 insertRequest body models.ReqUser true "JSON request body based on column name"
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /register [post]
func (h *HttpSqlx) RegisterUser(c echo.Context) error {
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

	result := h.commandUsecase.RegisterUser(c.Request().Context(), params)
	if result.Error != nil {
		return utils.Response(nil, "Register failed", http.StatusBadRequest, c)
	}

	return utils.Response(result.Data, "Register User", http.StatusCreated, c)
}
