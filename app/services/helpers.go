package services

import (
	"github.com/revel/revel/logger"
	"revel-project/app"
)

func InitService(log logger.MultiLogger) BaseService {
	return BaseService{
		db: app.DBConnection.GetDB(),
		log: log,
	}
}
