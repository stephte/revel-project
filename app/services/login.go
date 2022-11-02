package services

import (
	"revel-project/app/services/dtos"
	"revel-project/app/utilities"
	"revel-project/app/models"
	"errors"
	"time"
)

type LoginService struct {
	BaseService
}

func(this LoginService) LoginUser(credentials dtos.LoginDTO) (dtos.LoginTokenDTO, dtos.ErrorDTO) {
	// help protect against brute force attack
	killSomeTime(967, 2978)

	var user models.User
	var err error
	// make sure username + password are correct
	user, err = this.findUserByEmail(credentials.Email)

	if err != nil {
		return dtos.LoginTokenDTO{}, dtos.CreateErrorDTO(errors.New("Invalid Credentials"), 0, false)
	}

	if !utilities.ComparePasswords(user.EncryptedPassword, credentials.Password) {
		return dtos.LoginTokenDTO{}, dtos.CreateErrorDTO(errors.New("Invalid Credentials"), 0, false)
	}

	// then create JWT token and return it
	token, tokenErr := this.genToken(user)

	if tokenErr != nil {
		return dtos.LoginTokenDTO{}, dtos.CreateErrorDTO(errors.New("Error logging in"), 500, false)
	}

	return dtos.LoginTokenDTO{token}, dtos.ErrorDTO{}
}

func (this LoginService) genToken(user models.User) (string, error) {
	header := dtos.JWTHeaderDTO{
		Algorithm: "HS256",
		Type: "JWT",
	}

	payload := dtos.JWTPayloadDTO{
		Key: user.Key.String(),
		Expiration: time.Now().Add(time.Hour * 4).Unix(),
		Issuer: "revel-project",
	}

	return generateJWT(header, payload)
}
