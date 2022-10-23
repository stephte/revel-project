package controllers

import (
	"github.com/revel/revel"
	"revel-project/app/services/dtos"
)

type BaseController struct {
	*revel.Controller
}


func (bc BaseController) RenderErrorJSON(err error, status int) revel.Result {
	if status == 0 {
		status = 400
	}

	bc.Response.ContentType = "application/json"
	bc.Response.Status = status

	return bc.RenderJSON(dtos.ErrorResponse{err.Error()})
}
