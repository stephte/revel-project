package controllers

import (
	"revel-app/app/services"
	"github.com/mitchellh/mapstructure"
)

type UsersController struct {
	BaseController
}

type UserParams struct {
	FirstName 	string	`json:firstName`
	LastName		string	`json:lastName`
	Email				string	`json:email`
}

type CreateUserParams struct {
	UserParams
	Password		string	`json:password`
}



func (uc UsersController) Create() revel.Result {
	// store params into struct, acts as param validation
	var params CreateUserParams
	uc.Params.BindJSON(&params)

	// convert struct to map to send to service
	var data map[string]interface{}
	mapstructure.Decode(params, &data)

	response, err := services.UserService.CreateUser(data)

	if err != nil {
		return uc.RenderError(err)
	}

	// return response!
	return uc.RenderJSON(response)
}
