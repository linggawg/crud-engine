package handler

import (
	conn "crud-engine/pkg/database"
	"crud-engine/pkg/utils"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xwb1989/sqlparser"
)

type PageFetchInput struct {
	Page *int
	Size *int
	Sort string
}

// Get ShowData godoc
// @Summary      Find all Data
// @Description  Find all data by statement parameter
// @Tags         CrudEngine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param        isQuery    query     boolean  false  "if isQuery is true, the sql query statement is fetched directly from the path table"
// @Param        isDistinct    query     boolean  false  " DISTINCT statement is used to return only distinct (different) values. "
// @Param        colls    query     string  false  "column name (ex : username, email)"
// @Param        query    query     string  false  "where condition query sql"
// @Param        pageSize    query     int  false  "limit per page"
// @Param        pageNo    query     int  false  "page number list data "
// @Param        sortBy    query     string  false  "sorting data by column name (ex : name ASC / name DESC)"
// @Security Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /sql/{table} [get]
func (h *HttpSqlx) Get(c echo.Context) error {
	errorMessage := os.Getenv("GET_ERROR_MESSAGE")
	var (
		sqlStatement string
		sqlTotal     string
	)
	table := c.Param("table")

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

	pagination, err := getPagination(c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	if c.QueryParam("isQuery") == "true" {
		table, _ = url.QueryUnescape(table)
		sqlStatement = table
		sqlTotal = table
	} else {
		primaryKey, err := getPrimaryKey(database, table, dbs.Dialect, c)
		if err != nil {
			log.Println(err)
			return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
		}

		isDistinct := c.QueryParam("isDistinct")
		if isDistinct == "true" {
			isDistinct = "DISTINCT "
		}

		query := c.QueryParam("query")
		if query != "" {
			query = " WHERE " + query
		}

		colls := c.QueryParam("colls")
		if colls == "" {
			colls = "*"
		}

		sqlStatement = "SELECT " + isDistinct + colls + " FROM " + table + query
		sqlTotal = sqlStatement
		sqlStatement = setQueryPagination(sqlStatement, primaryKey.column, pagination)
	}

	_, err = sqlparser.Parse(sqlStatement)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	var totalItems int64
	err = database.QueryRow("SELECT COUNT('total') FROM (" + sqlTotal + ") as total").Scan(&totalItems)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	rows, err := database.QueryContext(c.Request().Context(), sqlStatement)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	_, err = json.Marshal(tableData)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	result := map[string]interface{}{
		"content":       tableData,
		"totalElements": totalItems,
		"maxPage": func() *float64 {
			if pagination.Size != nil {
				maxPage := math.Ceil(float64(totalItems)/float64(*pagination.Size)) - 1
				return &maxPage
			} else {
				return nil
			}
		}(),
		"page": pagination.Page,
		"size": pagination.Size,
	}

	return utils.Response(result, "List table "+table, http.StatusOK, c)
}

func setQueryPagination(query string, primaryKey string, p *PageFetchInput) (newQuery string) {
	if p != nil {
		if p.Sort != "" {
			query += " ORDER BY " + p.Sort
		} else {
			query += " ORDER BY " + primaryKey + " ASC"
		}
		if p.Size != nil {
			sz := *p.Size
			size := strconv.Itoa(sz)
			query += " LIMIT " + size
			if p.Page != nil {
				pg := *p.Page
				page := strconv.Itoa(pg * sz)
				query += " OFFSET " + page
			}
		}
	}
	return query
}

func getPagination(c echo.Context) (p *PageFetchInput, e error) {
	var (
		getParam = c.QueryParams()
		page     *int
		size     *int
	)
	allParams := getParam.Get("pageSize")
	if allParams != "" {
		sz, err := strconv.Atoi(allParams)
		if err != nil {
			return p, err
		}
		size = &sz
		allParams := getParam.Get("pageNo")
		if allParams != "" {
			pg, err := strconv.Atoi(allParams)
			if err != nil {
				return p, err
			}
			page = &pg
		}
	}

	p = &PageFetchInput{
		Page: page,
		Size: size,
		Sort: getParam.Get("sortBy"),
	}
	return p, nil
}
