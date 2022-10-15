package handlers

import (
	"encoding/json"
	"engine/bin/middleware"
	commandsServices "engine/bin/modules/services/repositories/commands"
	queriesServices "engine/bin/modules/services/repositories/queries"
	models "engine/bin/modules/users-services/models/domain"
	"engine/bin/modules/users-services/repositories/commands"
	"engine/bin/modules/users-services/repositories/queries"
	"engine/bin/modules/users-services/usecases"
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
	servicesPostgreQuery := queriesServices.NewServicesQuery(db)
	servicesPostgreCommand := commandsServices.NewServicesCommand(db)

	postgreQuery := queries.NewUsersServicesQuery(db)
	postgreCommand := commands.NewUsersServicesCommand(db)
	commandUsecase := usecases.NewUsersServicesCommandUsecase(postgreCommand, postgreQuery, servicesPostgreCommand, servicesPostgreQuery)

	return &HttpSqlx{
		CommandUsecase: commandUsecase,
	}
}

// Mount function
func (h *HttpSqlx) Mount(echoGroup *echo.Group) {
	echoGroup.POST("/v1/users-services/default", h.GenerateDefaultUsersServices, middleware.VerifyBearer())
	echoGroup.DELETE("/v1/users-services/default", h.DeleteByUsersIdAndServiceUrl, middleware.VerifyBearer())
}

// Post GenerateDefaultUsersServices godoc
// @Summary      Generate Default services and users-services
// @Description  Generate Default services GET, POST, PUT, PATCH, DELETE and create users services
// @Tags         Services
// @Accept       json
// @Produce      json
// @Param		 UsersServicesRequest body models.UsersServicesRequest true "JSON request body based on column name, required service_url, user_id, db_id"
// @Security     Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /v1/users-services/default [post]
func (h *HttpSqlx) GenerateDefaultUsersServices(c echo.Context) error {
	var (
		payload models.UsersServicesRequest
	)
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	err = validator.New().Struct(payload)
	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	header, _ := json.Marshal(c.Get("opts"))
	json.Unmarshal(header, &payload.Opts)

	result := h.CommandUsecase.InsertAllByServices(c.Request().Context(), payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Generate Default Users Services Method", http.StatusCreated, c)
}

// Delete DeleteByUsersIdAndServiceUrl godoc
// @Summary      Delete Default users-services
// @Description  Delete all users-services by user_id and services_url
// @Tags         Services
// @Accept       json
// @Produce      json
// @Param		 UsersServicesRequest body models.UsersServicesRequest true "JSON request body based on column name, required service_url, user_id"
// @Security     Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /v1/users-services/default [delete]
func (h *HttpSqlx) DeleteByUsersIdAndServiceUrl(c echo.Context) error {
	var (
		payload models.UsersServicesRequest
	)
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	err = validator.New().Struct(payload)
	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	header, _ := json.Marshal(c.Get("opts"))
	json.Unmarshal(header, &payload.Opts)

	result := h.CommandUsecase.DeleteByUsersIdAndServiceUrl(c.Request().Context(), payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Delete Default Users Services Method", http.StatusOK, c)
}
