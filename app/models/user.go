package models

import (
	"revel-project/app/utilities"
	"gorm.io/gorm"
	"strings"
	"errors"
	"time"
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
	FirstName									string		`gorm:"not null"`
	LastName									string		`gorm:"not null"`
	Email											string		`gorm:"uniqueIndex;not null"`
	Role											int				`gorm:"default:0"`
	PasswordResetToken				[]byte
	PasswordResetExpiration		int64
	Password									string		`gorm:"-"`
	EncryptedPassword					[]byte		`gorm:"not null"`
	PasswordLastUpdated				int64			// Probs dont need this
}


func(u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}


func(u *User) BeforeSave(tx *gorm.DB) (err error) {
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

	if u.Role == 0 {
		u.Role = 1
	} else if !utilities.IntArrContains(getRoles(), u.Role) {
		return errors.New("Not a valid User Role")
	}

	return nil
}


func(this User) CheckPassword(givenPW string) bool {
	return utilities.CompareStringWithHash(this.EncryptedPassword, givenPW)
}

func(this User) CheckPWResetToken(givenToken string) bool {
	return utilities.CompareStringWithHash(this.PasswordResetToken, givenToken)
}


func(u *User) handleEmail() error {
	if !utilities.IsValidEmail(u.Email) {
		return errors.New("Invalid email")
	}

	u.Email = strings.ToLower(u.Email)

	return nil
}


func(u *User) handlePassword() error {
	if !utilities.ValidatePassword(u.Password) {
		return errors.New("Password invalid")
	}

	hash, err := utilities.CreateHash(u.Password)

	if err != nil {
		return err
	}

	u.EncryptedPassword = hash
	u.PasswordLastUpdated = time.Now().Unix()

	return nil
}
