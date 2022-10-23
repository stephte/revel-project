package controllers

import (
	"revel-project/app/services"
	"revel-project/app/services/dtos"
	"github.com/revel/revel"
)

type UsersController struct {
	BaseController
}


func (uc UsersController) Create() revel.Result {
	var dto dtos.CreateUserDTO
	uc.Params.BindJSON(&dto)

	service := services.UserService{services.InitService()}

	response, err := service.CreateUser(dto)

	if err != nil {
		return uc.RenderErrorJSON(err, 0)
	}

	// return response!
	return uc.RenderJSON(response)
}
