package controllers

import (
	"revel-project/app/services/dtos"
	"revel-project/app/services"
	"github.com/google/uuid"
	"github.com/revel/revel"
	"strings"
	"errors"
)

type BaseController struct {
	*revel.Controller
	currentUserKey		uuid.UUID
}


func (bc BaseController) renderErrorJSON(errorDTO dtos.ErrorDTO) revel.Result {
	if errorDTO.Status == 0 {
		errorDTO.Status = 400
	}

	bc.Response.Status = errorDTO.Status

	return bc.RenderJSON(errorDTO)
}

func (bc *BaseController) validateJWT() (revel.Result) {
	headers := bc.Request.Header

	jwt := headers.Get("Authorization")

	token := strings.Replace(jwt, "Bearer ", "", 1)

	jwtValid, tokenErr := services.ValidateJWT(token)

	if tokenErr != nil {
		errDTO := dtos.CreateErrorDTO(tokenErr, 401, false)
		return bc.renderErrorJSON(errDTO)
	} else if !jwtValid {
		errDTO := dtos.CreateErrorDTO(errors.New("Token Expired"), 401, true)
		return bc.renderErrorJSON(errDTO)
	}

	bc.currentUserKey = services.GetJWTUserKey(token)

	return nil
}
