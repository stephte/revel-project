package models

import (
	"revel-project/app/utilities/auth"
	"revel-project/app/utilities"
	"gorm.io/gorm"
	"strings"
	"errors"
	"time"
	"fmt"
)

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


func(this *User) BeforeCreate(tx *gorm.DB) error {
	this.Email = strings.ToLower(this.Email)

	if this.Role == 0 {
		this.Role = 1
	}

	if pwErr := this.handlePassword(); pwErr != nil {
		return pwErr
	}

	return nil
}


func(this *User) BeforeUpdate(tx *gorm.DB) (err error) {
	typ := utilities.GetType(tx.Statement.Dest)
	
	// normal User update is assumed to be with a map
	if typ == "map[string]interface {}" {
		mp, ok := tx.Statement.Dest.(map[string]interface{})
		if ok {
			return this.beforeSaveWithMap(mp, tx)
		}

		return errors.New("Invalid map")
	}

	// this works for password handling since we Save the user with the password already on the User
	if this.Password != "" {
		if pwErr := this.handlePassword(); pwErr != nil {
			return pwErr
		}
	}

	return nil
}


// pure validation checks should go in AfterSave
func(this *User) AfterSave(tx *gorm.DB) (err error) {
	if !utilities.IntArrContains(auth.GetUserRoles(), this.Role) {
		return errors.New("Not a valid User Role")
	}

	if !utilities.IsValidEmail(this.Email) {
		return errors.New("Invalid User Email")
	}

	return nil
}


func(this User) CheckPassword(givenPW string) bool {
	return auth.CompareStringWithHash(this.EncryptedPassword, givenPW)
}


func(this User) CheckPWResetToken(givenToken string) bool {
	return auth.CompareStringWithHash(this.PasswordResetToken, givenToken)
}


// ---------- Private ----------


func(this *User) beforeSaveWithMap(data map[string]interface{}, tx *gorm.DB) error {
	// first check if email key exists
	genericEmail, exists := data["Email"]
	if exists {
		// then get email string from map and update email with lowercase email
		email, isString := genericEmail.(string)
		if isString {
			tx.Statement.SetColumn("Email", strings.ToLower(email))
		} else {
			return errors.New("Email must be a string")
		}
	}

	return nil
}


func(u *User) handlePassword() error {
	if !auth.ValidatePassword(u.Password) {
		return errors.New("Password invalid")
	}

	hash, err := auth.CreateHash(u.Password)
	if err != nil {
		return err
	}

	u.EncryptedPassword = hash
	u.PasswordLastUpdated = time.Now().Unix()

	return nil
}
