package controllers

import (
	"revel-project/app/services"
	"revel-project/app/services/dtos"
	"github.com/revel/revel"
)

type UsersController struct {
	BaseController
}

// TODO (low priority): implement query ordering and sorting
func (uc UsersController) Index() revel.Result {
	service := services.UserService{services.InitService()}

	response, err := service.GetUsers()

	if err != nil {
		return uc.RenderErrorJSON(err, 401)
	}

	// return response!
	return uc.RenderJSON(response)
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

func (uc UsersController) Update() revel.Result {
	var dto dtos.UserDTO
	uc.Params.BindJSON(&dto)

	service := services.UserService{services.InitService()}

	response, err := service.UpdateUser(dto)

	if err != nil {
		return uc.RenderErrorJSON(err, 0)
	}

	// return response!
	return uc.RenderJSON(response)
}
