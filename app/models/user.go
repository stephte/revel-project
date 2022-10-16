package models

import (
	"strings"
	"errors"
	"revel-project/app/utilities"
)

type User struct {
	BaseModel
	// Key					uuid.UUID	`gorm:"type:uuid;uniqueIndex;not null;default:uuid_generate_v4()"`
	FirstName							string		`gorm:"not null"`
	LastName							string		`gorm:"not null"`
	Email									string		`gorm:"uniqueIndex;not null"`
	Password							string		`gorm:"-:all`
	EncryptedPassword			[]byte		`gorm:"not null"`
	Role									string		`gorm:"default:'regular'"`
	PasswordResetToken		string
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// first validate email
	newEmail, emailErr := utilities.HandleEmail(u.Email)

	if emailErr != nil {
		return emailErr
	}

	u.Email = newEmail

	if u.Password != nil && u.Password != "" {
		// err := u.handlePassword()
		if pwErr := u.handlePassword(); pwErr != nil {
			return pwErr
		}
	}

	return nil
}

func (u *User) handlePassword() error {
	if !utilities.ValidatePassword(u.Password) {
		return errors.New("Password invalid")
	}

	hash, err = utilities.CreateHash(u.Password)

	if err != nil {
		return err
	}

	u.EncryptedPassword = hash

	return nil
}
