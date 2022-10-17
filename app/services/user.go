package services

import (
	_ "gorm.io/gorm"
	"revel-project/app/models"
	"github.com/mitchellh/mapstructure"
)

type UserService struct {
	BaseService
}

func (us UserService) FindUserByKey(userKey string) (models.User, error) {
	user := models.User{}
	if findErr := us.DB.Where("Key = $1", userKey).First(&user).Error; findErr != nil {
		return user, findErr
	}

	return user, nil
}

func (us UserService) CreateUser(data map[string]interface{}) (map[string]interface{}, error) {
	var user models.User
	decodeErr := mapstructure.Decode(data, &user)
	
	if decodeErr != nil {
		return nil, decodeErr
	}

	if createErr := us.DB.Create(&user).Error; createErr != nil {
		return nil, createErr
	}

	// convert struct back to map
	var rv map[string]interface{}
	mapErr := mapstructure.Decode(user, &rv)

	if mapErr != nil {
		return nil, mapErr
	}

	return rv, nil
}

func (us UserService) UpdateUser(userKey string, data map[string]interface{}) (error) {
	user, findErr := us.FindUserByKey(userKey)

	if findErr != nil {
		return findErr
	}

	if updateErr := us.DB.Model(&user).Updates(data).Error; updateErr != nil {
		return updateErr
	}

	return nil
}