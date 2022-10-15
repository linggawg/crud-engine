package handlers

import (
	"encoding/json"
	models "engine/bin/modules/engine/models/domain"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Post UpdateData godoc
// @Summary      Insert Data
// @Description  Insert data by column name in format JSON
// @Tags         Engine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param		 insertRequest body map[string]interface{} true "JSON request body based on column name"
// @Security     Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /v1/{table} [post]
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
