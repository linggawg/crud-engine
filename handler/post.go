package handler

import (
	"crud-engine/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Post UpdateData godoc
// @Summary      Insert Data
// @Description  Insert data by column name in format JSON
// @Tags         CrudEngine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param		 insertRequest body map[string]interface{} true "JSON request body based on column name"
// @Security Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /sql/{table} [post]
func (h *HttpSqlx) Post(c echo.Context) error {
	errorMessage := os.Getenv("POST_ERROR_MESSAGE")
	table := c.Param("table")
	db := h.db

	var jsonBody map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	primaryKey, err := getPrimaryKey(db, table, c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	columns, values := sqlStatement(primaryKey, jsonBody)
	sqlStatement := "INSERT INTO " + table + " (" + columns + ") VALUES (" + values + ");"

	_, err = db.ExecContext(c.Request().Context(), sqlStatement)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	return utils.Response(jsonBody, "successfully insert "+table, http.StatusCreated, c)
}

func sqlStatement(primaryKey *PrimaryKey, jsonBody map[string]interface{}) (columns string, values string) {
	for key := range jsonBody {
		if key != primaryKey.column {
			columns += key + ", "
			values += "'" + jsonBody[key].(string) + "', "
		}
	}
	if !strings.Contains(columns, primaryKey.column) && "int" != primaryKey.format {
		columns += primaryKey.column + ", "
		values += "'" + uuid.New().String() + "', "
	}
	return strings.TrimRight(columns, ", "), strings.TrimRight(values, ", ")
}
