package services

import (
	"revel-project/app/models"
	"revel-project/app/services/dtos"
	"revel-project/app/services/mappers"
	"github.com/google/uuid"
)

type UserService struct {
	BaseService
}

func (us UserService) findUserByKey(userKey uuid.UUID) (models.User, error) {
	user := models.User{}
	if findErr := us.DB.Where("Key = $1", userKey).First(&user).Error; findErr != nil {
		return user, findErr
	}

	return user, nil
}

// TODO (low priority): implement query ordering and sorting
func (us UserService) GetUsers() ([]dtos.UserDTO, error) {
	var users []models.User

	if err := us.DB.Order("created_at").Find(&users).Error; err != nil {
		return nil, err
	}

	return mappers.MapUsersToUserDTOs(users), nil
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

// TODO: handle Role validation once login validation implemented
func (us UserService) UpdateUser(dto dtos.UserDTO) (dtos.UserDTO, error) {
	user, findErr := us.findUserByKey(dto.Key)

	if findErr != nil {
		return dtos.UserDTO{}, findErr
	}

	updatedUser := mappers.MapUserDTOToUser(dto)

	if updateErr := us.DB.Model(&user).Updates(updatedUser).Error; updateErr != nil {
		return dtos.UserDTO{}, updateErr
	}

	return mappers.MapUserToUserDTO(user), nil
}
