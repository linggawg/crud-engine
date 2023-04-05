package handlers

import (
	"encoding/json"
	models "engine/bin/modules/engine/models/domain"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *EngineHTTPHandler) Post(c echo.Context) error {
	var (
		jsonBody map[string]interface{}
	)
	header, _ := json.Marshal(c.Get("opts"))

	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = err.Error()
		return utils.ResponseError(errObj, c)
	}
	payload := &models.EngineRequest{Table: c.Param("table"), Data: jsonBody}
	json.Unmarshal(header, &payload.Opts)

	result := h.DbsQueryUsecase.GetDbsConnection(c.Request().Context(), payload.Opts.UserID, c.Request().Method, payload.Table, "")
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}
	var engineConfig models.EngineConfig
	byteSub, _ := json.Marshal(result.Data)
	json.Unmarshal(byteSub, &engineConfig)

	result = h.EngineCommandUsecase.Insert(c.Request().Context(), engineConfig, payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}
	return utils.Response(result.Data, "Successfully insert", http.StatusOK, c)
}
