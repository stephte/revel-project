package models

import (
	"errors"
	"revel-project/app/utilities"
	"gorm.io/gorm"
	"fmt"
)

func getRoles() []int {
	return []int{
		1, // "regular",
		2, // "admin",
		3, // "super_admin"
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
	newEmail, emailErr := utilities.HandleEmail(u.Email)

	if emailErr != nil {
		return emailErr
	}

	u.Email = newEmail

	if u.Password != "" {
		// err := u.handlePassword()
		if pwErr := u.handlePassword(); pwErr != nil {
			return pwErr
		}
	}

	if u.Role == 0 {
		u.Role = 1
	} else if !utilities.IntArrContains(getRoles(), u.Role) {
		return errors.New("Not a valid User Role")
	}

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
