package handlers

import (
	"encoding/json"
	"engine/bin/middleware"
	models "engine/bin/modules/services/models/domain"
	"engine/bin/modules/services/repositories/commands"
	"engine/bin/modules/services/repositories/queries"
	"engine/bin/modules/services/usecases"
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
	postgreQuery := queries.NewServicesQuery(db)
	postgreCommand := commands.NewServicesCommand(db)
	commandUsecase := usecases.NewServicesCommandUsecase(postgreCommand, postgreQuery)

	return &HttpSqlx{
		CommandUsecase: commandUsecase,
	}
}

func (h *HttpSqlx) Mount(echoGroup *echo.Group) {
	echoGroup.DELETE("/v1/services/default", h.DeleteDefaultServicesMethod, middleware.VerifyBearer())
}

func (h *HttpSqlx) DeleteDefaultServicesMethod(c echo.Context) error {
	var (
		payload models.ServicesRequest
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

	result := h.CommandUsecase.DeleteByServiceUrl(c.Request().Context(), payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Delete Default Services Method", http.StatusOK, c)
}
