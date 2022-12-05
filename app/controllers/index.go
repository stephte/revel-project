package controllers

import (
	"revel-project/app/services/dtos"
	"revel-project/app/services"
	"github.com/revel/revel"
	"strings"
	"errors"
)

type BaseController struct {
	*revel.Controller
	baseService					*services.BaseService
}


func (bc BaseController) renderErrorJSON(errorDTO dtos.ErrorDTO) revel.Result {
	if errorDTO.Status == 0 {
		errorDTO.Status = 400
	}

	bc.Response.Status = errorDTO.Status

	return bc.RenderJSON(errorDTO)
}


func (bc *BaseController) validateJWT(isPWReset bool) (revel.Result) {
	headers := bc.Request.Header

	jwt := headers.Get("Authorization")

	token := strings.Replace(jwt, "Bearer ", "", 1)

	bc.setBaseService()

	authService := services.AuthService{BaseService: bc.baseService}
	jwtValid, tokenErrDTO := authService.ValidateJWT(token, isPWReset)

	if tokenErrDTO.Exists() {
		return bc.renderErrorJSON(tokenErrDTO)
	} else if !jwtValid {
		errDTO := dtos.CreateErrorDTO(errors.New("Token Expired"), 401, true)
		return bc.renderErrorJSON(errDTO)
	}

	return nil
}


func(bc *BaseController) setBaseService() {
	service := services.InitService(bc.Log)
	bc.baseService = &service
}
