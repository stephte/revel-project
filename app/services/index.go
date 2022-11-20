package services

import (
	"github.com/revel/revel/logger"
	"revel-project/app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseService struct {
	db 						*gorm.DB
	log					  logger.MultiLogger // what type does this need to be?
	currentUser		models.User
}


func (bs *BaseService) setCurrentUser(userKey uuid.UUID) error {
	user, findErr := bs.findUserByKey(userKey)
	if findErr != nil {
		// log error
		bs.log.Errorf("Can't find user with key: %s", userKey.String())
		return findErr
	}

	bs.currentUser = user

	return nil
}


func (bs *BaseService) setCurrentUserByEmail(email string) error {
	user, findErr := bs.findUserByEmail(email)
	if findErr != nil {
		// log error
		bs.log.Errorf("Can't find user with email: %s", email)
		return findErr
	}

	bs.currentUser = user

	return nil
}


func (bs BaseService) findUserByKey(userKey uuid.UUID) (models.User, error) {
	user := models.User{}
	if findErr := bs.db.Where("Key = $1", userKey).First(&user).Error; findErr != nil {
		return user, findErr
	}

	return user, nil
}


func (bs BaseService) findUserByEmail(userEmail string) (models.User, error) {
	user := models.User{}
	if findErr := bs.db.Where("Email = $1", userEmail).First(&user).Error; findErr != nil {
		return user, findErr
	}

	return user, nil
}


func (bs BaseService) validateUserHasAccess(accessNeeded int) bool {
	return bs.currentUser.Role >= accessNeeded
}
