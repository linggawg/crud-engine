package handler

import (
	conn "crud-engine/pkg/database"
	"crud-engine/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	var (
		err          error
		errorMessage = os.Getenv("POST_ERROR_MESSAGE")
		table        = c.Param("table")
		jsonBody     map[string]interface{}
	)

	dbs, err := h.GetDbsConn(c, h.db)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	database, err := conn.InitDbs(conn.SQLXConfig{
		Host:     dbs.Host,
		Port:     strconv.Itoa(dbs.Port),
		Name:     dbs.Name,
		Username: dbs.Username,
		Password: func() string {
			if dbs.Password != nil {
				return *dbs.Password
			} else {
				return ""
			}
		}(),
		Dialect: dbs.Dialect,
	})
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}
	defer database.Close()

	err = json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	primaryKey, err := getPrimaryKey(database, table, dbs.Dialect, c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	informationSchemas, err := sqlIsNullable(database, table, dbs.Dialect, c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	sqlStatement, args, err := sqlStatementInsert(table, primaryKey, jsonBody, informationSchemas)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage+err.Error(), http.StatusBadRequest, c)
	}
	_, err = database.ExecContext(c.Request().Context(), sqlStatement, args...)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	return utils.Response(jsonBody, "successfully insert "+table, http.StatusCreated, c)
}

func sqlStatementInsert(table string, primaryKey *PrimaryKey, jsonBody map[string]interface{}, informationSchemas []InformationSchema) (sql string, args []interface{}, err error) {
	var columns, values []string
	for key := range jsonBody {
		if key != primaryKey.column {
			columns = append(columns, key)
			values = append(values, "?")
			if jsonBody[key] == nil {
				for _, i := range informationSchemas {
					if i.IsNullable == "NO" && i.ColumName == key {
						errorMessage := fmt.Sprintf(": Error:validation for '%s' failed on the 'required' tag", i.ColumName)
						return "", args, errors.New(errorMessage)
					}
				}
			}
			args = append(args, jsonBody[key])
		}
	}
	if !strings.EqualFold(primaryKey.format, "int") {
		columns = append(columns, primaryKey.column)
		values = append(values, "?")
		args = append(args, uuid.New().String())
	}

	for _, i := range informationSchemas {
		if i.IsNullable == "NO" && i.ColumName != primaryKey.column {
			if !strings.Contains(strings.Join(columns, ","), i.ColumName) {
				errorMessage := fmt.Sprintf(": Error:validation for '%s' failed on the 'required' tag", i.ColumName)
				return "", args, errors.New(errorMessage)
			}
		}
	}
	sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", table, strings.Join(columns, ","), strings.Join(values, ","))
	return SetQuery(sql), args, nil
}
