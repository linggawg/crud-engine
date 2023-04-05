package handlers

import (
	"encoding/json"
	"engine/bin/middleware"
	rolesQueries "engine/bin/modules/roles/repositories/queries"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/modules/users/repositories/commands"
	"engine/bin/modules/users/repositories/queries"
	"engine/bin/modules/users/usecases"
	conn "engine/bin/pkg/databases"
	"engine/bin/pkg/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type HttpSqlx struct {
	CommandUsecase usecases.CommandUsecase
}

func New() *HttpSqlx {
	db := conn.InitSqlx()
	rolesPostgreQuery := rolesQueries.NewRolesQuery(db)
	postgreQuery := queries.NewUsersQuery(db)
	postgreCommand := commands.NewUsersCommand(db)
	commandUsecase := usecases.NewUsersCommandUsecase(postgreCommand, postgreQuery, rolesPostgreQuery)

	return &HttpSqlx{
		CommandUsecase: commandUsecase,
	}
}

func (h *HttpSqlx) Mount(echoGroup *echo.Group) {
	echoGroup.POST("/v1/login", h.Login, middleware.NoAuth())
	echoGroup.POST("/v1/register", h.RegisterUser, middleware.VerifyBearer())
}

func (h *HttpSqlx) Login(c echo.Context) error {
	var (
		params models.ReqLogin
	)
	err := json.NewDecoder(c.Request().Body).Decode(&params)
	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	err = validator.New().Struct(params)
	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	result := h.CommandUsecase.Login(c.Request().Context(), params)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Login User", http.StatusOK, c)
}

func (h *HttpSqlx) RegisterUser(c echo.Context) error {
	var (
		params models.ReqUser
	)
	err := json.NewDecoder(c.Request().Body).Decode(&params)
	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	err = validator.New().Struct(params)
	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	header, _ := json.Marshal(c.Get("opts"))
	json.Unmarshal(header, &params.Opts)

	result := h.CommandUsecase.RegisterUser(c.Request().Context(), params)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Register User", http.StatusCreated, c)
}
