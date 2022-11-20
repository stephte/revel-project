package controllers

import (
	"revel-project/app/services/dtos"
	"revel-project/app/services"
	"github.com/revel/revel"
)

type LoginController struct {
	BaseController
}


func(lc LoginController) Login() revel.Result {
	var dto dtos.LoginDTO
	lc.Params.BindJSON(&dto)
	
	lc.setBaseService()
	service := services.LoginService{lc.baseService}

	tokenDTO, errDTO := service.LoginUser(dto)

	if errDTO.Exists() {
		return lc.renderErrorJSON(errDTO)
	}

	return lc.RenderJSON(tokenDTO)
}


func(lc LoginController) StartPWReset() revel.Result {
	var dto dtos.EmailDTO
	lc.Params.BindJSON(&dto)

	lc.setBaseService()
	service := services.LoginService{lc.baseService}

	errDTO := service.StartPWReset(dto)

	if errDTO.Exists() {
		return lc.renderErrorJSON(errDTO)
	}

	return lc.RenderJSON(map[string]string{"msg": "Password reset email will be sent if a user with that email exists."})
}


func(lc LoginController) ConfirmPasswordResetToken() revel.Result {
	var dto dtos.ConfirmResetTokenDTO
	lc.Params.BindJSON(&dto)

	lc.setBaseService()
	service := services.LoginService{lc.baseService}

	res, errDTO := service.ConfirmResetToken(dto)

	if errDTO.Exists() {
		return lc.renderErrorJSON(errDTO)
	}

	return lc.RenderJSON(res)
}


func(lc LoginController) ResetPassword() revel.Result {
	errResponse := lc.validateJWT(true)
	if errResponse != nil {
		return errResponse
	}

	var dto dtos.ResetPWDTO
	lc.Params.BindJSON(&dto)

	service := services.LoginService{lc.baseService}
	
	tokenDTO, errDTO := service.UpdateUserPassword(dto)

	if errDTO.Exists() {
		return lc.renderErrorJSON(errDTO)
	}

	return lc.RenderJSON(tokenDTO)
}
