package v1

import (
	"gin-scaffold/app/common/keys"
	"gin-scaffold/app/common/request"
	"gin-scaffold/app/common/response"
	"gin-scaffold/app/services"
	"gin-scaffold/locales"
	"gin-scaffold/utils/validate"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {
	form := new(request.Login)
	if err := c.ShouldBindJSON(form); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.FailC(c, locales.CodeInvalidParam)
			return
		}
		response.ValidateFail(c, validate.GetValidationErrors(
			c.GetString("lang"), form, errs,
		))
		return
	}

	if err, user := services.UserService.Login(form); err != nil {
		response.Fail(c, err)
		return
	} else {
		tokenData, err := services.JwtService.GetToken(services.AppGuardName, user)
		if err != nil {
			response.FailC(c, locales.CodeInvalidGenToken)
		}
		response.Success(c, tokenData)
	}
}

func Register(c *gin.Context) {
	form := new(request.Register)
	if err := c.ShouldBindJSON(form); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.FailC(c, locales.CodeInvalidParam)
			return
		}
		response.ValidateFail(c, validate.GetValidationErrors(
			c.GetString("lang"), form, errs,
		))
		return
	}

	if user, err := services.UserService.Register(*form); err != nil {
		response.Fail(c, err)
	} else {
		response.Success(c, user)
	}
}

func Info(c *gin.Context) {
	err, user := services.UserService.GetUserInfo(c.GetInt64(keys.CtxUserIDKey))
	if err != nil {
		response.Fail(c, err)
	}
	response.Success(c, user)
}

func Logout(c *gin.Context) {
	token, ok := c.Keys[keys.CtxUserToken]
	if ok {
		err := services.JwtService.JoinBlackList(token.(*jwt.Token))
		if err != nil {
			response.FailC(c, locales.CodeInvalidLogout)
			return
		}
	}

	response.Success(c, nil)
}
