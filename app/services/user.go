package services

import (
	"revel-project/app/models"
	"revel-project/app/services/dtos"
	"revel-project/app/services/mappers"
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

// takes in CreateUserDTO, returns UserDTO
func (us UserService) CreateUser(dto dtos.CreateUserDTO) (dtos.UserDTO, error) {
	user := mappers.MapCreateUserDTOToUser(dto)

	if createErr := us.DB.Create(&user).Error; createErr != nil {
		return dtos.UserDTO{}, createErr
	}

	rv := mappers.MapUserToUserDTO(user)

	return rv, nil
}

// TODO: update this to work with new DTOs
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