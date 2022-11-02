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
	
	service := services.LoginService{services.InitService(lc.Log)}

	tokenDTO, errDTO := service.LoginUser(dto)

	if errDTO.Exists() {
		return lc.renderErrorJSON(errDTO)
	}

	return lc.RenderJSON(tokenDTO)
}
