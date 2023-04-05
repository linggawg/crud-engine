package handlers

import (
	"encoding/json"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/pkg/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *EngineHTTPHandler) Get(c echo.Context) error {
	table := c.Param("table")
	query := make(map[string]interface{})
	header, _ := json.Marshal(c.Get("opts"))

	for key, value := range c.QueryParams() {
		query[key] = value[0]
		b, err := strconv.ParseBool(value[0])
		if err == nil {
			query[key] = b
		}
		v, err := strconv.Atoi(value[0])
		if err == nil {
			query[key] = v
		}
	}
	payload, _ := json.Marshal(query)
	var params models.GetList
	json.Unmarshal(payload, &params)
	json.Unmarshal(header, &params.Opts)

	result := h.DbsQueryUsecase.GetDbsConnection(c.Request().Context(), params.Opts.UserID, c.Request().Method, table, params.Key)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}
	var engineConfig models.EngineConfig
	byteSub, _ := json.Marshal(result.Data)
	json.Unmarshal(byteSub, &engineConfig)

	result = h.EngineQueryUsecase.Get(c.Request().Context(), engineConfig, table, &params)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.PaginationResponse(result.Data, result.MetaData, "List table ", http.StatusOK, c)
}
