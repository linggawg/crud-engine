package handlers

import (
	"encoding/json"
	"engine/bin/config"
	models "engine/bin/modules/engine/models/domain"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Put UpdateData godoc
// @Summary      Put Update Data
// @Description  Update data by field_id and data by column name in format JSON, Requires sending complete data from the table, so that when there is data that is not sent it will be changed to the default value
// @Tags         Engine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param        value   path    string  true  "Value of id"
// @Param        field_id    query     string  true  "Update based on field_id "
// @Param		 updateRequest body map[string]interface{} true "JSON request body based on column name"
// @Security     Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /v1/{table}/{value} [put]
func (h *EngineHTTPHandler) Put(c echo.Context) error {
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
	payload := &models.EngineRequest{
		FieldId: c.QueryParam("field_id"),
		Table:   c.Param("table"),
		Value:   c.Param("value"),
		Data:    jsonBody,
	}
	json.Unmarshal(header, &payload.Opts)

	result := h.DbsQueryUsecase.GetDbsConnection(c.Request().Context(), payload.Opts.UserID, c.Request().Method, payload.Table, "")
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}
	var engineConfig models.EngineConfig
	byteSub, _ := json.Marshal(result.Data)
	json.Unmarshal(byteSub, &engineConfig)

	if config.GlobalEnv.StrictRestfulAPI {
		result = h.EngineCommandUsecase.Update(c.Request().Context(), engineConfig, payload)
	} else {
		result = h.EngineCommandUsecase.Patch(c.Request().Context(), engineConfig, payload)
	}
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	message := "successfully update " + payload.Table + " with " + payload.FieldId + ": " + payload.Value
	return utils.Response(result.Data, message, http.StatusOK, c)
}
