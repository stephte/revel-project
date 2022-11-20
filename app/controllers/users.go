package controllers

import (
	"revel-project/app/services/dtos"
	"revel-project/app/services"
	"github.com/revel/revel"
)

type UsersController struct {
	BaseController
}

// TODO (low priority): implement query ordering and sorting
func (uc UsersController) Index() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	service := services.UserService{uc.baseService}

	response, errDTO := service.GetUsers()

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	// return response!
	return uc.RenderJSON(response)
}

func (uc UsersController) Create() revel.Result {
	var dto dtos.CreateUserDTO
	uc.Params.BindJSON(&dto)

	uc.setBaseService()
	service := services.UserService{uc.baseService}

	response, errDTO := service.CreateUser(dto)

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(response)
}

func (uc UsersController) Update() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	var dto dtos.UserDTO
	uc.Params.BindJSON(&dto)

	service := services.UserService{uc.baseService}

	response, errDTO := service.UpdateUser(dto)

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(response)
}
