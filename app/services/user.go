package services

import (
	"gin-scaffold/app/common/request"
	"gin-scaffold/app/models"
	"gin-scaffold/global"
	"gin-scaffold/locales"
	"gin-scaffold/utils"
)

type userService struct{}

var UserService = new(userService)

func (userService *userService) Register(params request.Register) (user models.User, err error) {
	var result = global.App.DB.Where("mobile = ?", params.Mobile).Select("id").First(&models.User{})
	if result.RowsAffected != 0 {
		err = &locales.CErrors{Code: locales.CodeMobileExist}
		return
	}

	userID := utils.GenID()
	user = models.User{
		ID:       userID,
		Name:     params.Name,
		Mobile:   params.Mobile,
		Password: utils.BcryptMake([]byte(params.Password)),
	}
	err = global.App.DB.Create(&user).Error
	if err != nil {
		err = &locales.CErrors{Code: locales.CodeDBCreateError}
	}

	return
}

func (userService *userService) Login(params *request.Login) (err error, user *models.User) {
	err = global.App.DB.Where("mobile = ?", params.Mobile).First(&user).Error
	if err != nil || !utils.BcryptMakeCheck(params.Password, user.Password) {
		err = &locales.CErrors{Code: locales.CodeInvalidPassword}
	}
	return
}

func (userService *userService) GetUserInfo(id int64) (err error, user models.User) {
	err = global.App.DB.First(&user, id).Error
	if err != nil {
		err = &locales.CErrors{Code: locales.CodeUserNotExist}
	}
	return
}
