package user

import (
	"gorm.io/gorm"
	"encoding/json"
	"revel-project/app/models"
	"revel-project/app/utilities"
	"github.com/mitchellh/mapstructure"
)

type UserService struct {
	DB *gorm.DB
}

func (us UserService) FindUserByKey(userKey string) (models.User, error) {
	user := models.User{}
	us.DB.where("Key = $1", userKey).First(&user)

	return user, user.Error
}

func (us UserService) CreateUser(data map[string]interface{}) (map[string]interface{}, error) {
	var user models.User
	mapErr := mapstructure.Decode(data, &user)

	if mapErr != nil {
		return nil, mapErr
	}

	if createErr := us.DB.Create(&user).Error; createErr != nil {
		return nil, createErr
	}

	// convert struct back to map
	var rv map[string]interface{}
	mapErr := mapstructure.Decode(user, &rv)

	if mapErr != nil {
		return nil, e
	}

	return rv, nil
}

func (us UserService) UpdateUser(userKey string, data map[string]interface{}) (error) {
	user, findErr := us.FindUserByKey(userKey)

	if findErr != nil {
		return findErr
	}

	if updateErr := us.DB.Model(&user).Updates(data).Error; err != nil {
		return updateErr
	}

	return nil
}