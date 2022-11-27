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

	userDTOs, errDTO := service.GetUsers()
	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	// return response!
	return uc.RenderJSON(userDTOs)
}


func(uc UsersController) Find() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	userKeyStr := uc.Params.Route.Get("userKey")

	service := services.UserService{uc.baseService}

	userDTO, errDTO := service.GetUser(userKeyStr)
	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(userDTO)
}


func (uc UsersController) Create() revel.Result {
	var dto dtos.CreateUserDTO
	uc.Params.BindJSON(&dto)

	uc.setBaseService()
	service := services.UserService{uc.baseService}

	userDTO, errDTO := service.CreateUser(dto)

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(userDTO)
}


// this endpoint validates the request data against the UserDTO,
// but keeps it as a map so that only the included data is updated
// (GORM only updates non-zero fields when updating with struct)
func (uc UsersController) Update() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	userKeyStr := uc.Params.Route.Get("userKey")

	var data map[string]interface{}
	uc.Params.BindJSON(&data)

	// validate User update data
	validatedData, dataErr := dtos.ValidateUserMap(data)
	if dataErr != nil {
		return uc.renderErrorJSON(dtos.CreateErrorDTO(dataErr, 0, false))
	}

	service := services.UserService{uc.baseService}

	userDTO, errDTO := service.UpdateUser(userKeyStr, validatedData)

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(userDTO)
}
