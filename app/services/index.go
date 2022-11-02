package services

import (
	"github.com/revel/revel/logger"
	"revel-project/app/models"
	"github.com/google/uuid"
	"revel-project/app"
	"gorm.io/gorm"
)

type BaseService struct {
	DB *gorm.DB
	log					  logger.MultiLogger // what type does this need to be?
	currentUser		models.User
}

func (bs *BaseService) SetCurrentUser(userKey uuid.UUID) {
	user := models.User{}
	if findErr := bs.DB.Where("Key = $1", userKey.String()).First(&user).Error; findErr != nil {
		// log error
		bs.log.Errorf("User not found: %s", userKey.String())
	}

	bs.currentUser = user
}

func InitService(log logger.MultiLogger) BaseService {
	return BaseService{
		DB: app.DB,
		log: log,
	}
}

func InitServiceWithCurrentUser(log logger.MultiLogger, userKey uuid.UUID) BaseService {
	bs := InitService(log)
	bs.SetCurrentUser(userKey)

	return bs
}

func (bs BaseService) findUserByKey(userKey uuid.UUID) (models.User, error) {
	user := models.User{}
	if findErr := bs.DB.Where("Key = $1", userKey).First(&user).Error; findErr != nil {
		return user, findErr
	}

	return user, nil
}

func (bs BaseService) findUserByEmail(userEmail string) (models.User, error) {
	user := models.User{}
	if findErr := bs.DB.Where("Email = $1", userEmail).First(&user).Error; findErr != nil {
		return user, findErr
	}

	return user, nil
}


func (bs BaseService) validateUserHasAccess(accessNeeded int) bool {
	// bs.log.Infof("Needed: %d, Has: %d", accessNeeded, bs.currentUser.Role)
	return bs.currentUser.Role >= accessNeeded
}
