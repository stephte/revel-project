package services

import (
	"revel-project/app/services/mappers"
	"revel-project/app/utilities/auth"
	"revel-project/app/services/dtos"
	"revel-project/app/models"
	"github.com/google/uuid"
	"errors"
)

type UserService struct {
	*BaseService
}


func(this UserService) GetUser(userKeyStr string) (dtos.UserDTO, dtos.ErrorDTO) {
	key, parseErr := uuid.Parse(userKeyStr)
	if parseErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(parseErr, 0, false)
	}

	user, findErr := this.findUserByKey(key)
	if findErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(findErr, 404, false)
	}

	if this.currentUser.ID != user.ID && !this.validateUserHasAccess(auth.AdminAccess()) {
		return dtos.UserDTO{}, dtos.AccessDeniedError()
	}

	return mappers.MapUserToUserDTO(user), dtos.ErrorDTO{}
}


// TODO (low priority): implement query ordering and sorting
func (this UserService) GetUsers() ([]dtos.UserDTO, dtos.ErrorDTO) {
	if !this.validateUserHasAccess(auth.AdminAccess()) {
		return nil, dtos.AccessDeniedError()
	}

	var users []models.User

	if err := this.db.Order("created_at").Find(&users).Error; err != nil {
		return nil, dtos.CreateErrorDTO(err, 500, false)
	}

	return mappers.MapUsersToUserDTOs(users), dtos.ErrorDTO{}
}


// takes in CreateUserDTO, returns UserDTO
func (this UserService) CreateUser(dto dtos.CreateUserDTO) (dtos.UserDTO, dtos.ErrorDTO) {
	user := mappers.MapCreateUserDTOToUser(dto)

	if createErr := this.db.Create(&user).Error; createErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(createErr, 0, false)
	}

	rv := mappers.MapUserToUserDTO(user)

	return rv, dtos.ErrorDTO{}
}


func (this UserService) UpdateUser(userKeyStr string, data map[string]interface{}) (dtos.UserDTO, dtos.ErrorDTO) {
	key, parseErr := uuid.Parse(userKeyStr)
	if parseErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(parseErr, 0, false)
	}

	user, findErr := this.findUserByKey(key)
	if findErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(findErr, 404, false)
	}

	// check if role exists in data; else resume as if its equal to users current role
	var role int
	_, exists := data["Role"]
	if exists {
		roleFloat, isFloat := data["Role"].(float64)
		if !isFloat {
			return dtos.UserDTO{}, dtos.CreateErrorDTO(errors.New("Role is not a float64"), 0, false)
		}

		role = int(roleFloat)
	} else {
		role = user.Role
	}

	if !this.validateUserHasAccess(auth.SuperAdminAccess()) && (role != user.Role || this.currentUser.ID != user.ID) {
		return dtos.UserDTO{}, dtos.AccessDeniedError()
	}

	if updateErr := this.db.Model(&user).Updates(data).Error; updateErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(updateErr, 0, false)
	}

	return mappers.MapUserToUserDTO(user), dtos.ErrorDTO{}
}
