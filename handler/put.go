package handler

import (
	"crud-engine/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// Put UpdateData godoc
// @Summary      Update Data
// @Description  Update data by ID (primary key) and data by column name in format JSON
// @Tags         CrudEngine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param        id   path    string  true  "Primary Key"
// @Param		 updateRequest body map[string]interface{} true "JSON request body based on column name"
// @Security Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /sql/{table}/{id} [put]
func (h *HttpSqlx) Put(c echo.Context) error {
	errorMessage := os.Getenv("PUT_ERROR_MESSAGE")
	table := c.Param("table")
	db := h.db

	var jsonBody map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	var setData string
	informationSchemas, err := sqlIsNullable(db, table, c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	for key := range jsonBody {
		if jsonBody[key] == nil {
			for _, i := range informationSchemas {
				if i.ColumName == key && i.IsNullable == "NO" {
					errM := fmt.Sprintf(", Error:validation for '%s' failed on the 'required' tag", i.ColumName)
					log.Println(errorMessage)
					return utils.Response(nil, errorMessage+errM, http.StatusBadRequest, c)
				}
			}
		}
		setData += key + fmt.Sprintf("='%s', ", jsonBody[key])
	}
	setData = strings.TrimRight(setData, ", ")

	value := c.Param("value")
	field := c.QueryParam("field_id")
	sqlStatement := "UPDATE " + table + " SET " + setData + " WHERE " + field + " ='" + value + "'"

	_, err = db.ExecContext(c.Request().Context(), sqlStatement)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	message := "successfully update " + table + " with " + field + " " + value
	return utils.Response(jsonBody, message, http.StatusOK, c)
}
