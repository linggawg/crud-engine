package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Result struct {
	Data  interface{}
	Error interface{}
}

type BaseWrapperModel struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Response(data interface{}, message string, code int, c echo.Context) error {
	success := false

	if code < http.StatusBadRequest {
		success = true
	}

	result := BaseWrapperModel{
		Success: success,
		Data:    data,
		Message: message,
		Code:    code,
	}

	return c.JSON(code, result)
}
