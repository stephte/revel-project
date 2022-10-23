package controllers

import (
	"github.com/revel/revel"
)

type BaseController struct {
	*revel.Controller
}

type ErrorResponse struct {
	Err 		string		`json:"error"`
}

func (bc BaseController) RenderErrorJSON(err error, status int) revel.Result {
	if status == 0 {
		status = 400
	}

	bc.Response.ContentType = "application/json"
	bc.Response.Status = status

	return bc.RenderJSON(ErrorResponse{err.Error()})
}
