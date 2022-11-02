package services

import (
	"revel-project/app/services/mappers"
	"revel-project/app/services/dtos"
	"revel-project/app/models"
	"github.com/google/uuid"
	"errors"
)

type UserService struct {
	BaseService
}

func (this UserService) GetUserByKey(userKey uuid.UUID) (dtos.UserDTO, dtos.ErrorDTO) {
	user, findErr := this.findUserByKey(userKey)

	if findErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(findErr, 404, false)
	}

	return mappers.MapUserToUserDTO(user), dtos.ErrorDTO{}
}

// TODO (low priority): implement query ordering and sorting
func (this UserService) GetUsers() ([]dtos.UserDTO, dtos.ErrorDTO) {
	if !this.validateUserHasAccess(1) {
		return nil, dtos.CreateErrorDTO(errors.New("Access Denied"), 401, false)
	}

	var users []models.User

	if err := this.DB.Order("created_at").Find(&users).Error; err != nil {
		return nil, dtos.CreateErrorDTO(err, 500, false)
	}

	return mappers.MapUsersToUserDTOs(users), dtos.ErrorDTO{}
}

// takes in CreateUserDTO, returns UserDTO
func (this UserService) CreateUser(dto dtos.CreateUserDTO) (dtos.UserDTO, dtos.ErrorDTO) {
	user := mappers.MapCreateUserDTOToUser(dto)

	if createErr := this.DB.Create(&user).Error; createErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(createErr, 0, false)
	}

	rv := mappers.MapUserToUserDTO(user)

	return rv, dtos.ErrorDTO{}
}

func (this UserService) UpdateUser(dto dtos.UserDTO) (dtos.UserDTO, dtos.ErrorDTO) {
	user, findErr := this.findUserByKey(dto.Key)

	if findErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(findErr, 404, false)
	}

	// handle validation (only super admins can update Role)
	if ((dto.Role != user.Role) && (this.currentUser.Role < 2)) || (!this.validateUserHasAccess(1) && !(this.currentUser.ID == user.ID)) {
		return dto, dtos.CreateErrorDTO(errors.New("Access Denied"), 401, false)
	}

	updatedUser := mappers.MapUserDTOToUser(dto)

	if updateErr := this.DB.Model(&user).Updates(updatedUser).Error; updateErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(updateErr, 0, false)
	}

	return mappers.MapUserToUserDTO(user), dtos.ErrorDTO{}
}
