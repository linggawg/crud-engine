package handlers

import (
	"encoding/json"
	dbsModels "engine/bin/modules/dbs/models/domain"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Delete DeleteDate godoc
// @Summary      Delete Data
// @Description  Delete data by ID (primary key)
// @Tags         CrudEngine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param        id   path    string  true  "Primary Key"
// @Security Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /sql/{table}/{id} [delete]
func (h *EngineHTTPHandler) Delete(c echo.Context) error {
	header, _ := json.Marshal(c.Get("opts"))
	payload := &models.EngineRequest{
		FieldId: c.QueryParam("field_id"),
		Table:   c.Param("table"),
		Value:   c.Param("value"),
	}
	json.Unmarshal(header, &payload.Opts)

	result := h.dbsQueryUsecase.GetDbsConnection(c.Request().Context(), payload.Opts.UserID, c.Request().Method, payload.Table, false)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}
	var dbsConn dbsModels.Dbs
	byteSub, _ := json.Marshal(result.Data)
	json.Unmarshal(byteSub, &dbsConn)

	result = h.engineCommandUsecase.Delete(c.Request().Context(), dbsConn, payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	message := "successfully delete " + payload.Table + " with " + payload.FieldId + " " + payload.Value
	return utils.Response(result.Data, message, http.StatusOK, c)
}
