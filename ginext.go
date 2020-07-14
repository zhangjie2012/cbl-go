package cbl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

var (
	codeOK    = 0
	codeError = 1
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Resp{
		Code:  codeOK,
		Data:  data,
		Error: ErrNone,
	})
}

func ErrorResponse(c *gin.Context, err interface{}) {
	e := ""
	switch t := err.(type) {
	case string:
		e = t
	case error:
		e = t.Error()
	default:
		e = "server error"
	}

	c.JSON(http.StatusOK, &Resp{
		Code:  codeError,
		Data:  nil,
		Error: e,
	})
}
