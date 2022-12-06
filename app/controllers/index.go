package controllers

import (
	"revel-project/app/services/dtos"
	"revel-project/app/services"
	"github.com/revel/revel"
	"strconv"
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


func(bc BaseController) getPaginationDTO() (dtos.PaginationDTO, dtos.ErrorDTO) {
	sort := bc.Params.Query.Get("order")
	pageStr := bc.Params.Query.Get("page")
	limitStr := bc.Params.Query.Get("limit")

	page, limit := 0, 0
	var err error
	
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return dtos.PaginationDTO{}, dtos.CreateErrorDTO(err, 400, false)
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return dtos.PaginationDTO{}, dtos.CreateErrorDTO(err, 400, false)
		}
	}

	rv := dtos.PaginationDTO{}
	rv.Init(sort, page, limit)

	return rv, dtos.ErrorDTO{}
}

func(bc BaseController) getRequestPath() string {
	url := bc.Request.URL.String()

	// now strip out everything after the '?' (if any)
	stringsArr := strings.Split(url, "?")

	return stringsArr[0]
}
