package response

import (
	"gin-scaffold/locales"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    locales.ResCode `json:"code"`
	Data    interface{}     `json:"data"`
	Message string          `json:"message"`
}

type ValidateResponse struct {
	Code    locales.ResCode `json:"code"`
	Errors  interface{}     `json:"errors"`
	Message string          `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		locales.CodeSuccess,
		data,
		"success",
	})
}

func FailC(c *gin.Context, code locales.ResCode) {
	c.JSON(http.StatusOK, Response{
		code,
		nil,
		code.Msg(c.GetString("locale")),
	})
}

func Fail(c *gin.Context, err error) {
	if cErr, ok := err.(*locales.CErrors); ok {
		FailC(c, cErr.Code)
	} else {
		c.JSON(http.StatusOK, Response{
			locales.CodeServerBusy,
			nil,
			err.Error(),
		})
	}
}

func ValidateFail(c *gin.Context, errors interface{}) {

	c.JSON(http.StatusOK, ValidateResponse{
		locales.CodeInvalidParam,
		errors,
		locales.CodeInvalidParam.Msg(c.GetString("locale")),
	})
}
