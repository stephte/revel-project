package services

import (
	"gorm.io/gorm"
	"revel-project/app"
)

type BaseService struct {
	DB *gorm.DB
}

func InitService() BaseService {
	return BaseService{app.DB}
}
