package handler

import (
	"crud-engine/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"strings"
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
	var (
		jsonBody     map[string]interface{}
		err          error
		errorMessage = os.Getenv("PUT_ERROR_MESSAGE")
		db           = h.db
		table        = c.Param("table")
		value        = c.Param("value")
		field        = c.QueryParam("field_id")
	)

	err = json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	informationSchemas, err := sqlIsNullable(db, table, c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	primaryKey, err := getPrimaryKey(db, table, c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	sqlStatement, args, err := sqlStatementUpdate(table, field, value, primaryKey, jsonBody, informationSchemas)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage+err.Error(), http.StatusBadRequest, c)
	}

	_, err = db.ExecContext(c.Request().Context(), SetQuery(sqlStatement), args...)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	message := "successfully update " + table + " with " + field + " " + value
	return utils.Response(jsonBody, message, http.StatusOK, c)
}

func sqlStatementUpdate(table, fieldId, valueId string, primaryKey *PrimaryKey, jsonBody map[string]interface{}, informationSchemas []InformationSchema) (sql string, args []interface{}, err error) {
	var setData []string
	for key := range jsonBody {
		if key != primaryKey.column {
			if jsonBody[key] == nil {
				for _, i := range informationSchemas {
					if i.ColumName == key && i.IsNullable == "NO" {
						err = errors.New(fmt.Sprintf(": Error:validation for '%s' failed on the 'required' tag", i.ColumName))
						return "", args, err
					}
				}
			}
			setData = append(setData, fmt.Sprintf("%s=?", key))
			args = append(args, jsonBody[key])
		}
	}
	args = append(args, valueId)
	sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", table, strings.Join(setData, ","), fieldId)
	return SetQuery(sql), args, nil
}
