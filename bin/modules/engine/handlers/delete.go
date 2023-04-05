package handlers

import (
	"encoding/json"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *EngineHTTPHandler) Delete(c echo.Context) error {
	header, _ := json.Marshal(c.Get("opts"))
	payload := &models.EngineRequest{
		FieldId: c.QueryParam("field_id"),
		Table:   c.Param("table"),
		Value:   c.Param("value"),
	}
	json.Unmarshal(header, &payload.Opts)

	result := h.DbsQueryUsecase.GetDbsConnection(c.Request().Context(), payload.Opts.UserID, c.Request().Method, payload.Table, "")
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}
	var engineConfig models.EngineConfig
	byteSub, _ := json.Marshal(result.Data)
	json.Unmarshal(byteSub, &engineConfig)

	result = h.EngineCommandUsecase.Delete(c.Request().Context(), engineConfig, payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	message := "successfully delete " + payload.Table + " with " + payload.FieldId + " " + payload.Value
	return utils.Response(result.Data, message, http.StatusOK, c)
}
