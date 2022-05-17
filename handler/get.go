package handler

import (
	"crud-engine/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PageFetchInput struct {
	Page *int
	Size *int
	Sort string
}

func Get(c echo.Context) error {
	db := config.CreateCon()
	var sqlStatement string

	pagination, err := getPagination(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	isQuery := c.QueryParam("isQuery")
	if isQuery == "true" {
		sqlStatement = c.Param("table")
	} else {
		table := c.Param("table")
		primaryKey, err := getPrimaryKey(db, table, c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		isDistinct := ""
		if c.QueryParam("isDistinct") == "true" {
			isDistinct = "DISTINCT "
		}

		query := ""
		if c.QueryParam("query") != "" {
			query = " WHERE " + c.QueryParam("query")
		}

		colls := c.QueryParam("colls")
		if colls == "" {
			colls = "*"
		}

		sqlStatement = "SELECT " + isDistinct + colls + " FROM " + table + query
		sqlStatement = setQueryPagination(sqlStatement, primaryKey, pagination)
	}
	fmt.Println(sqlStatement)

	rows, err := db.QueryContext(c.Request().Context(), sqlStatement)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
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

	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, string(jsonData))
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
