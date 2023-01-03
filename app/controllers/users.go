package controllers

import (
	"revel-project/app/services/dtos"
	"revel-project/app/services"
	"github.com/revel/revel"
)

type UsersController struct {
	BaseController
}


func(uc UsersController) Index() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	dto, paginationErrDTO := uc.getPaginationDTO()
	if paginationErrDTO.Exists() {
		return uc.renderErrorJSON(paginationErrDTO)
	}

	path := uc.getRequestPath()

	service := services.UserService{BaseService: uc.baseService}

	result, errDTO := service.GetUsers(dto, path)
	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(result)
}


func(uc UsersController) Find() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	userKeyStr := uc.Params.Route.Get("userKey")

	service := services.UserService{BaseService: uc.baseService}

	userDTO, errDTO := service.GetUser(userKeyStr)
	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(userDTO)
}


func(uc UsersController) Create() revel.Result {
	errResponse := uc.validateJWT(false)
	// Don't return response, don't always need to be authed to create an account
	if errResponse != nil {
		uc.Response.Status = 200
	}

	var dto dtos.CreateUserDTO
	uc.Params.BindJSON(&dto)

	service := services.UserService{BaseService: uc.baseService}

	userDTO, errDTO := service.CreateUser(dto)

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(userDTO)
}


// PATCH version of User update
// this endpoint validates the request data against the UserDTO,
// but keeps it as a map so that only the included data is updated
// (GORM only updates non-zero fields when updating with struct)
func(uc UsersController) Update() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	userKeyStr := uc.Params.Route.Get("userKey")

	var data map[string]interface{}
	uc.Params.BindJSON(&data)

	service := services.UserService{BaseService: uc.baseService}

	userDTO, errDTO := service.UpdateUser(userKeyStr, data)

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(userDTO)
}


// PUT version of User update (expects all user data) (prefer above PATCH version)
func (uc UsersController) UpdateOG() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	userKeyStr := uc.Params.Route.Get("userKey")

	var dto dtos.UserDTO
	uc.Params.BindJSON(&dto)

	service := services.UserService{BaseService: uc.baseService}

	response, errDTO := service.UpdateUserOG(userKeyStr, dto)

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.RenderJSON(response)
}


func(uc UsersController) Delete() revel.Result {
	errResponse := uc.validateJWT(false)
	if errResponse != nil {
		return errResponse
	}

	userKeyStr := uc.Params.Route.Get("userKey")

	service := services.UserService{BaseService: uc.baseService}

	errDTO := service.DeleteUser(userKeyStr)

	if errDTO.Exists() {
		return uc.renderErrorJSON(errDTO)
	}

	return uc.blankSuccessResponse()
}
