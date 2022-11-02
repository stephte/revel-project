package models

import (
	"errors"
	"revel-project/app/utilities"
	"gorm.io/gorm"
	"strings"
	"fmt"
)

func getRoles() []int {
	return []int{
		0, // "regular",
		1, // "admin",
		2, // "super_admin"
	}
}

type User struct {
	BaseModel
	FirstName							string		`gorm:"not null"`
	LastName							string		`gorm:"not null"`
	Email									string		`gorm:"uniqueIndex;not null"`
	Role									int				`gorm:"default:0"`
	PasswordResetToken		string
	Password							string		`gorm:"-"`
	EncryptedPassword			[]byte		`gorm:"not null"`
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// first validate email
	emailErr := u.handleEmail()

	if emailErr != nil {
		return emailErr
	}

	if u.Password != "" {
		if pwErr := u.handlePassword(); pwErr != nil {
			return pwErr
		}
	}

	if !utilities.IntArrContains(getRoles(), u.Role) {
		return errors.New("Not a valid User Role")
	}

	return nil
}

func (u *User) handleEmail() error {
	if !utilities.IsValidEmail(u.Email) {
		return errors.New("Invalid email")
	}

	u.Email = strings.ToLower(u.Email)

	return nil
}

func (u *User) handlePassword() error {
	if !utilities.ValidatePassword(u.Password) {
		return errors.New("Password invalid")
	}

	hash, err := utilities.CreateHash(u.Password)

	if err != nil {
		return err
	}

	u.EncryptedPassword = hash

	return nil
}
