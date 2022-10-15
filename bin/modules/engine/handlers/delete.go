package handlers

import (
	"encoding/json"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Delete DeleteDate godoc
// @Summary      Delete Data
// @Description  Delete data by field_id
// @Tags         Engine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param        value   path    string  true  "Value Of id"
// @Param        field_id    query     string  true  "Delete based on field_id "
// @Security     Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /v1/{table}/{value} [delete]
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
