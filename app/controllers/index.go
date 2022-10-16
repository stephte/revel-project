package controllers

import (
	"github.com/revel/revel"
)

type BaseController struct {
	*revel.controller
}

// type ErrorResponse struct {
// 	Message string
// 	Error error
// }

// func (bc BaseController) ErrorResponse(error) {
// 	bc.RenderJSON()
// }
